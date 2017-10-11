package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

const (
	envKey     = "HOOKSHOT_URLS"
	failMsg    = "Failed to hit URL %d"
	successMsg = "Success! Hit %d URLs"
)

var urls []string

// Handle responds to lambda invocations by hitting all known URLs
func Handle(_ interface{}, _ *runtime.Context) (string, error) {
	for index, url := range urls {
		err := hitURL(url)
		if err != nil {
			return fmt.Sprintf(failMsg, index), fmt.Errorf(failMsg, index)
		}
	}
	return fmt.Sprintf(successMsg, len(urls)), nil
}

func hitURL(url string) error {
	_, err := http.Get(url)
	return err
}

func init() {
	urlString := os.Getenv(envKey)
	if urlString == "" {
		panic("No URLs provided")
	}
	urls = strings.Split(urlString, "|")
}

func main() {}
