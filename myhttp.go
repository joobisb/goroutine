package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"crypto/md5"
	"net/url"
	"sync"
	"time"
)

const (
	httpsScheme     = "https://"
	defaultWorkers  = 10
	apiTimeOutInSec = 30
)

type RequestClient interface {
	GetResponseHash(uri string, timeout time.Duration) (string, error)
}

type HTTPClient struct{}

func main() {

	var uris []string
	var workers int

	flag.IntVar(&workers, "parallel", 0, "parallel request limit")
	flag.Parse()

	if workers == 0 {
		workers = defaultWorkers
		uris = os.Args[1:]
	} else {
		uris = os.Args[3:]
	}

	if len(uris) == 0 {
		fmt.Fprintf(os.Stderr, "error: no uris found")
		return
	}

	apiClient := HTTPClient{}
	urlMD5HashMap, errs := execute(apiClient, uris, workers)

	//Incase of error from any of the APIs log the errors and
	//print the MD5 hash of rest of the successful APIs
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "error %s", err.Error())
		}
	}

	for k, v := range urlMD5HashMap {
		fmt.Printf("%s %s\n", k, v)
	}
}

func execute(apiClient RequestClient, uris []string, workers int) (map[string]string, []error) {
	urlMD5HashMap := make(map[string]string)
	var errs []error
	workQueue := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		worker := i

		go func(worker int, workQueue chan string) {
			for uri := range workQueue {
				responseHash, err := apiClient.GetResponseHash(uri, time.Second*apiTimeOutInSec)
				if err != nil {
					errs = append(errs, err)
				}

				urlMD5HashMap[uri] = responseHash
			}

			wg.Done()
		}(worker, workQueue)
	}

	go func() {
		for _, u := range uris {
			workQueue <- u
		}
		close(workQueue)
	}()
	wg.Wait()

	return urlMD5HashMap, errs
}

func (HTTPClient) GetResponseHash(uri string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if len(parsedURL.Scheme) == 0 {
		uri = httpsScheme + uri
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.Body != nil {
		defer res.Body.Close()
		res, _ := ioutil.ReadAll(res.Body)
		return fmt.Sprintf("%x", md5.Sum(res)), nil
	}
	return "", fmt.Errorf("no response found")
}
