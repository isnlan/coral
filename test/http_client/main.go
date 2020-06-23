package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/snlansky/coral/pkg/errors"

	"github.com/snlansky/coral/pkg/trace"
)

func main() {
	t, closer, err := trace.NewTracer("scpkg-gin-client", "127.0.0.1:6831")
	errors.Check(err)
	defer closer.Close()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8090/ping", nil)
	errors.Check(err)

	var i map[string]interface{}
	err = trace.DoRequest(context.Background(), req, &i)
	errors.Check(err)

	{
		f := func() error {
			time.Sleep(time.Second)
			return nil
		}

		ctx, err := trace.StartSpan(context.Background(), "test callback", f)
		errors.Check(err)
		fmt.Println(ctx)

		trace.StartSpan(ctx, "test callback2", f)
	}

	{
		//req, err := http.NewRequest("GET", "http://127.0.0.1:8090/ping", nil)
		//errors.Check(err)
		//
		//var i map[string]interface{}
		//err = trace.DoRequest(t, "http-request", trace.NewClient(), req, &i)
		//errors.Check(err)
		spanA := t.StartSpan("A7")
		time.Sleep(time.Second)
		spanA.Finish()
		spanB := t.StartSpan("B7", opentracing.ChildOf(spanA.Context()))
		time.Sleep(time.Second)
		spanB.Finish()
		spanC := t.StartSpan("C7", opentracing.ChildOf(spanB.Context()))
		time.Sleep(time.Second)
		spanC.Finish()

	}

	fmt.Println("->", i)
}
