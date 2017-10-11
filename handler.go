package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

const (
	KEY     = "HOOKSHOT_URLS"
	FAIL    = "Failed to hit URL %d"
	SUCCESS = "Success! Hit %d URLs"
)

var URLS []string

// Handle responds to lambda invocations by hitting all known URLs
func Handle(_ interface{}, _ *runtime.Context) (string, error) {
	for index, url := range URLS {
		err := hitURL(url)
		if err != nil {
			return fmt.Sprintf(FAIL, index), fmt.Errorf(FAIL, index)
		}
	}
	return fmt.Sprintf(SUCCESS, index), nil
}

func hitURL(url string) error {
	_, err := http.Get(url)
	return err
}

func init() {
	urlString := os.GetEnv(ENV_KEY)
	if urls == "" {
		panic("No URLs provided")
	}
	URLS = strings.Split(urlString, "|")
}

func main() {}
