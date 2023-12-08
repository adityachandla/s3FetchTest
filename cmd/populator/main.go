package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"

	"github.com/adityachandla/s3Bench/pkg/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type byteArray struct {
	barr []byte
	read int
}

func (ba *byteArray) Read(val []byte) (int, error) {
	if ba.read >= len(ba.barr) {
		return 0, io.EOF
	}
	r := copy(val, ba.barr[ba.read:])
	ba.read += r
	return r, nil
}

func RandomArray(size int) byteArray {
	arr := make([]byte, size)
	n, err := rand.Read(arr)
	if n != size || err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return byteArray{
		barr: arr,
		read: 0,
	}
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)
	log.Printf("Created client and manager")

	randomArray := RandomArray(common.FILE_SIZE_BYTES)
	log.Printf("Generated random array")

	input := &s3.PutObjectInput{
		Bucket: aws.String(common.BUCKET_NAME),
		Key:    aws.String(common.KEY),
		Body:   &randomArray,
	}
	_, err = uploader.Upload(context.TODO(), input)
	fmt.Println(err)
}
