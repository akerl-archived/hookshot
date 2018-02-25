package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akerl/go-lambda/s3"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/yaml.v2"
)

type target struct {
	URL    string `json:"url"`
	Method string `json:"method,omitempty"`
}

type config struct {
	Targets map[string]target `json:"targets"`
}

type targetList map[string]*http.Request

var targets targetList

func handler() error {
	if len(targets) == 0 {
		return fmt.Errorf("no targets found in config")
	}

	client := &http.Client{}

	for name, req := range targets {
		_, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to hit url: %s", name)
		}
	}
	log.Printf("success: hit %d urls", len(targets))
	return nil
}

func loadConfig() {
	bucket := os.Getenv("S3_BUCKET")
	path := os.Getenv("S3_KEY")
	fmt.Printf("%+v\n", os.Environ())
	if bucket == "" || path == "" {
		log.Print("variables not provided")
		return
	}

	obj, err := s3.GetObject(bucket, path)
	if err != nil {
		log.Print(err)
		return
	}

	c := config{}
	err = yaml.Unmarshal(obj, &c)
	if err != nil {
		log.Print(err)
		return
	}

	tl := make(targetList)
	for name, t := range c.Targets {
		req, err := http.NewRequest(t.Method, t.URL, nil)
		if err != nil {
			log.Print(err)
			return
		}
		tl[name] = req
	}
	targets = tl
}

func main() {
	loadConfig()
	lambda.Start(handler)
}
