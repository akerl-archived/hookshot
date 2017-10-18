package main

import (
	"fmt"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

const (
	failMsg    = "Failed to hit URL %s"
	successMsg = "Success! Hit %d URLs"
)

var tl targetList

// Handle responds to lambda invocations by hitting all known URLs
func Handle(_ interface{}, _ *runtime.Context) (string, error) {
	for name, target := range tl {
		err := target.hit()
		if err != nil {
			return logAndExit(failMsg, name, err)
		}
	}
	return logAndExit(successMsg, len(tl), nil)
}

func logAndExit(template string, data interface{}, err error) (string, error) {
	msg := fmt.Sprintf(template, data)
	fmt.Println(msg)
	if err != nil {
		err = fmt.Errorf(msg)
	}
	return msg, err
}

func init() {
	var err error
	tl, err = loadConfig()
	if err != nil {
		panic(err)
	}
}

func main() {}
