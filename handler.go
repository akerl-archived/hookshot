package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

const (
	bucketKey  = "S3_BUCKET"
	pathKey    = "S3_KEY"
	failMsg    = "Failed to hit URL %d"
	successMsg = "Success! Hit %d URLs"
)

var urls []string

// Handle responds to lambda invocations by hitting all known URLs
func Handle(_ interface{}, _ *runtime.Context) (string, error) {
	for index, url := range urls {
		err := hitURL(url)
		if err != nil {
			return logAndExit(failMsg, index, err)
		}
	}
	return logAndExit(successMsg, len(urls), nil)
}

func logAndExit(template string, data interface{}, err error) (string, error) {
	msg := fmt.Sprintf(template, data)
	fmt.Println(msg)
	if err != nil {
		err = fmt.Errorf(msg)
	}
	return msg, err
}

func hitURL(url string) error {
	_, err := http.Get(url)
	return err
}

func init() {
	bucket, path := parseEnv()
	client := s3Client()
	obj, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		panic(err)
	}
	urlString, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		panic(err)
	}
	urls = strings.Split(string(urlString), "|")
}

func parseEnv() (string, string) {
	bucket := os.Getenv(bucketKey)
	if bucket == "" {
		panic("No bucket provided")
	}
	path := os.Getenv(pathKey)
	if path == "" {
		panic("No path given")
	}
	return bucket, path
}

func s3Client() *s3.S3 {
	awsConfig := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	s := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))
	return s3.New(s)
}

func main() {}
