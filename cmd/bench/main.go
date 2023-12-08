package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/adityachandla/s3Bench/pkg/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	smithyrand "github.com/aws/smithy-go/rand"
)

var s3Client *s3.Client

var runs, parallelism, size int

func init() {
	flag.IntVar(&runs, "runs", 10, "Number of objects to fetch")
	flag.IntVar(&parallelism, "parallelism", 2, "Number of concurrent fetchers")
	flag.IntVar(&size, "size", 1<<19, "Size of fetched info in bytes")
}

func WithOperationTiming(msgHandler func(string)) func(*s3.Options) {
	return func(o *s3.Options) {
		o.APIOptions = append(o.APIOptions, addTimingMiddlewares(msgHandler))
		o.HTTPClient = &timedHTTPClient{
			client:     awshttp.NewBuildableClient(),
			msgHandler: msgHandler,
		}
	}
}

// PrintfMSGHandler writes messages to stdout.
func PrintfMSGHandler(msg string) {
	fmt.Printf("%s\n", msg)
}

type invokeIDKey struct{}

func setInvokeID(ctx context.Context, id string) context.Context {
	return middleware.WithStackValue(ctx, invokeIDKey{}, id)
}

func getInvokeID(ctx context.Context) string {
	id, _ := middleware.GetStackValue(ctx, invokeIDKey{}).(string)
	return id
}

func timeSpan(ctx context.Context, name string, consumer func(string)) func() {
	start := time.Now()
	return func() {
		elapsed := time.Now().Sub(start)
		consumer(fmt.Sprintf("[%s] %s: %s", getInvokeID(ctx), name, elapsed))
	}
}

type timedHTTPClient struct {
	client     *awshttp.BuildableClient
	msgHandler func(string)
}

func (c *timedHTTPClient) Do(r *http.Request) (*http.Response, error) {
	defer timeSpan(r.Context(), "http", c.msgHandler)()

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("inner client do: %v", err)
	}

	return resp, nil
}

type addInvokeIDMiddleware struct {
	msgHandler func(string)
}

func (*addInvokeIDMiddleware) ID() string { return "addInvokeID" }

func (*addInvokeIDMiddleware) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, md middleware.Metadata, err error,
) {
	id, err := smithyrand.NewUUID(smithyrand.Reader).GetUUID()
	if err != nil {
		return out, md, fmt.Errorf("new uuid: %v", err)
	}

	return next.HandleInitialize(setInvokeID(ctx, id), in)
}

type timeOperationMiddleware struct {
	msgHandler func(string)
}

func (*timeOperationMiddleware) ID() string { return "timeOperation" }

func (m *timeOperationMiddleware) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	middleware.InitializeOutput, middleware.Metadata, error,
) {
	defer timeSpan(ctx, "operation", m.msgHandler)()
	return next.HandleInitialize(ctx, in)
}

type emitMetadataMiddleware struct {
	msgHandler func(string)
}

func (*emitMetadataMiddleware) ID() string { return "emitMetadata" }

func (m *emitMetadataMiddleware) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	middleware.InitializeOutput, middleware.Metadata, error,
) {
	out, md, err := next.HandleInitialize(ctx, in)

	invokeID := getInvokeID(ctx)
	requestID, _ := awsmiddleware.GetRequestIDMetadata(md)
	m.msgHandler(fmt.Sprintf(`[%s] requestID = "%s"`, invokeID, requestID))

	return out, md, err
}

func addTimingMiddlewares(mh func(string)) func(*middleware.Stack) error {
	return func(s *middleware.Stack) error {
		if err := s.Initialize.Add(&timeOperationMiddleware{msgHandler: mh}, middleware.Before); err != nil {
			return fmt.Errorf("add time operation middleware: %v", err)
		}
		if err := s.Initialize.Add(&addInvokeIDMiddleware{msgHandler: mh}, middleware.Before); err != nil {
			return fmt.Errorf("add invoke id middleware: %v", err)
		}
		if err := s.Initialize.Insert(&emitMetadataMiddleware{msgHandler: mh}, "RegisterServiceMetadata", middleware.After); err != nil {
			return fmt.Errorf("add emit metadata middleware: %v", err)
		}
		return nil
	}
}

func main() {
	flag.Parse()
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatal(err)
	}
	s3Client = s3.NewFromConfig(cfg, WithOperationTiming(PrintfMSGHandler))
	fetchTimes := runBenchmark()
	printCsv(fetchTimes)
}

func runBenchmark() []int64 {
	timesChannel := make(chan int64)
	for p := 0; p < parallelism; p++ {
		go fetchBytes(timesChannel)
	}

	fetchTimes := make([]int64, parallelism*runs)
	for i := 0; i < runs*parallelism; i++ {
		fetchTimes[i] = <-timesChannel
	}
	return fetchTimes
}

func fetchBytes(timesChannel chan<- int64) {
	for i := 0; i < runs; i++ {
		timesChannel <- timeToFetch()
	}
}

func timeToFetch() int64 {
	startRange := rand.Intn(common.FILE_SIZE_BYTES - size)
	endRange := startRange + size
	req := s3.GetObjectInput{
		Bucket: aws.String(common.BUCKET_NAME),
		Key:    aws.String(common.KEY),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", startRange, endRange)),
	}
	//log.Printf("Requesting byte range %d-%d", startRange, endRange)
	start := time.Now().UnixMilli()
	_, err := s3Client.GetObject(context.TODO(), &req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	end := time.Now().UnixMilli()
	return end - start
}

func printCsv(arr []int64) {
	s := len(arr)
	for i := 0; i < s-1; i++ {
		fmt.Printf("%d,", arr[i])
	}
	fmt.Printf("%d\n", arr[s-1])
}
