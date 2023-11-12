package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/adityachandla/s3Bench/pkg/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	runs := flag.Int("runs", 10, "Number of objects to fetch")
	size := flag.Int("size", 128, "Size of fetched info in bytes")
	flag.Parse()
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(cfg)
	timesMillis := make([]int64, *runs)
	for i := 0; i < *runs; i++ {
		timesMillis[i] = timeToFetch(s3Client, *size)
	}
	for i := 0; i < *runs; i++ {
		fmt.Printf("%d,", timesMillis[i])
	}
	fmt.Printf("%d\n", timesMillis[*runs-1])
}

func timeToFetch(client *s3.Client, size int) int64 {
	startRange := rand.Intn(common.FILE_SIZE_BYTES - size)
	endRange := startRange + size
	req := s3.GetObjectInput{
		Bucket: aws.String(common.BUCKET_NAME),
		Key:    aws.String(common.KEY),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", startRange, endRange)),
	}
	log.Printf("Requesting byte range %d-%d", startRange, endRange)
	start := time.Now().UnixMilli()
	_, err := client.GetObject(context.TODO(), &req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	end := time.Now().UnixMilli()
	return end - start
}
