package main

import (
	"fmt"
	"github.com/avast/retry-go/v4"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	url := "https://www.baidu.com"
	var body []byte
	count := 0

	err := retry.Do(
		func() error {
			count++
			fmt.Sprintln(count, time.Now())
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		fmt.Println("error: %w", err)
	}

	fmt.Println(string(body))
}
