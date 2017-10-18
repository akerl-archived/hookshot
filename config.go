package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

const (
	bucketKey = "S3_BUCKET"
	pathKey   = "S3_KEY"
)

type target struct {
	URL    string
	Method string
}

type targetList map[string]target

type config struct {
	Targets targetList
}

func (t target) hit() error {
	m := t.Method
	if m == "" {
		m = "GET"
	}
	client := &http.Client{}
	req, err := http.NewRequest(m, t.URL, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}

func loadConfig() (config, error) {
	var c config

	bucket, path, err := parseEnv()
	if err != nil {
		return c, err
	}

	client := s3Client()

	obj, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return c, err
	}
	body, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(body, &c)
	return c, err
}

func parseEnv() (string, string, error) {
	bucket := os.Getenv(bucketKey)
	if bucket == "" {
		return "", "", fmt.Errorf("No bucket provided")
	}
	path := os.Getenv(pathKey)
	if path == "" {
		return "", "", fmt.Errorf("No path given")
	}
	return bucket, path, nil
}

func s3Client() *s3.S3 {
	awsConfig := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	s := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))
	return s3.New(s)
}
