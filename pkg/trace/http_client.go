package trace

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

func DoRequest(ctx context.Context, tr *http.Transport, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req, nethttp.OperationName(req.URL.EscapedPath()), nethttp.ComponentName("http client"))
	defer ht.Finish()

	var client *http.Client
	if tr != nil {
		client = &http.Client{Transport: &nethttp.Transport{RoundTripper: tr}}
	} else {
		client = &http.Client{Transport: &nethttp.Transport{}}
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(v)
}

func DoPost(ctx context.Context, url string, req interface{}, v interface{}) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	err = DoRequest(ctx, nil, request, v)
	if err != nil {
		return err
	}
	return nil
}

func DefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
