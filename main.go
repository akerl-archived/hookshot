package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akerl/go-lambda/s3"
	"github.com/aws/aws-lambda-go/lambda"
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
	c := config{}
	cf, err := s3.GetConfigFromEnv(&c)
	if err != nil {
		log.Print(err)
		return
	}
	cf.OnError = func(_ *ConfigFile, err error) {
		log.Print(err)
	}
	cf.Autoreload(60)

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
