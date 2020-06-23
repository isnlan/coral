package trace

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

var client = &http.Client{Transport: &nethttp.Transport{}}

func DoRequest(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req, nethttp.OperationName(req.URL.EscapedPath()), nethttp.ComponentName("http client"))
	defer ht.Finish()

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
	err = DoRequest(ctx, request, v)
	if err != nil {
		return err
	}
	return nil
}
