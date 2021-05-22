package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpRequestArgs struct {
	Method     string
	Url        string
	TimeoutSec time.Duration
	Body       io.Reader
	Headers    map[string]string
}

type Config struct {
	Address string
}

func HttpRequest(args *HttpRequestArgs) (int, []byte, error) {
	var err error

	req, err := http.NewRequest(args.Method, args.Url, args.Body)
	if err != nil {
		return 0, nil, err
	}
	if len(args.Headers) > 0 {
		for key, value := range args.Headers {
			req.Header.Add(key, value)
		}
	}
	client := &http.Client{Timeout: args.TimeoutSec * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, err
}

func GetApiUrl(baseUrl string, subPath string) string {
	if subPath == "" {
		return baseUrl
	}
	url := baseUrl + "/api/v1" + subPath
	return url
}
