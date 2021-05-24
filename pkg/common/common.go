package common

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"path"
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

type traceTransport struct {
}

func (t traceTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return res, err
	}

	log.Printf("[DEBUG] Request to %+v\n", request.URL)

	stream := res.Body
	buf, err := ioutil.ReadAll(stream)
	if err != nil {
		return res, err
	}

	log.Printf("[DEBUG] Response Struct %+v\n", res.Header)
	log.Printf("[DEBUG] Response Body %+v\n", string(buf))

	newStream := bytes.NewReader(buf)
	newCloseableStream := ioutil.NopCloser(newStream)
	res.Body = newCloseableStream

	return res, err
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

	client := &http.Client{Timeout: args.TimeoutSec * time.Second, Transport: &traceTransport{}}
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

func GetApiUrl(baseUrl string, subPath string) (string, error) {
	u, err := url2.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, "/api/v1", subPath)
	return u.String(), nil
}
