package metrics

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/prometheus/pkg/labels"

	"github.com/alecthomas/units"

	"github.com/isnlan/coral/pkg/errors"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/prometheus/prometheus/pkg/textparse"
)

const acceptHeader = `application/openmetrics-text; version=0.0.1,text/plain;version=0.0.4;q=0.5,*/*;q=0.1`

var (
	errBodySizeLimit      = errors.New("body size limit exceeded")
	errNameLabelMandatory = fmt.Errorf("missing metric name (%s label)", labels.MetricName)
)

type Metrics struct {
	Metrics string
	Name    string
	Hash    uint64
	Labels  labels.Labels
	Value   float64
	Time    int64
}

// targetScraper implements the scraper interface for a target.
type targetScraper struct {
	url     string
	client  *http.Client
	req     *http.Request
	timeout time.Duration

	gzipr *gzip.Reader
	buf   *bufio.Reader

	bodySizeLimit int64
}

func (s *targetScraper) scrape(ctx context.Context, w io.Writer) (string, error) {
	if s.req == nil {
		req, err := http.NewRequest("GET", s.url, nil)
		if err != nil {
			return "", err
		}
		req.Header.Add("Accept", acceptHeader)
		req.Header.Add("Accept-Encoding", "gzip")
		req.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", strconv.FormatFloat(s.timeout.Seconds(), 'f', -1, 64))

		s.req = req
	}

	resp, err := s.client.Do(s.req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("server returned HTTP status %s", resp.Status)
	}

	if s.bodySizeLimit <= 0 {
		s.bodySizeLimit = math.MaxInt64
	}
	if resp.Header.Get("Content-Encoding") != "gzip" {
		n, err := io.Copy(w, io.LimitReader(resp.Body, s.bodySizeLimit))
		if err != nil {
			return "", err
		}
		if n >= s.bodySizeLimit {
			return "", errBodySizeLimit
		}
		return resp.Header.Get("Content-Type"), nil
	}

	if s.gzipr == nil {
		s.buf = bufio.NewReader(resp.Body)
		s.gzipr, err = gzip.NewReader(s.buf)
		if err != nil {
			return "", err
		}
	} else {
		s.buf.Reset(resp.Body)
		if err = s.gzipr.Reset(s.buf); err != nil {
			return "", err
		}
	}

	n, err := io.Copy(w, io.LimitReader(s.gzipr, s.bodySizeLimit))

	s.gzipr.Close()

	if err != nil {
		return "", err
	}

	if n >= s.bodySizeLimit {
		return "", errBodySizeLimit
	}

	return resp.Header.Get("Content-Type"), nil
}

func ScrapeMetrics(ctx context.Context, url string) ([]*Metrics, error) {
	buf := bytes.NewBuffer(nil)

	s := &targetScraper{
		url:           url,
		client:        &http.Client{Transport: &nethttp.Transport{}},
		timeout:       time.Second * 30,
		bodySizeLimit: int64(units.MiB * 10),
	}

	contentType, err := s.scrape(ctx, buf)
	if err != nil {
		return nil, err
	}

	p := textparse.New(buf.Bytes(), contentType)

	var list []*Metrics

	for {
		et, err := p.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
			}

			break
		}

		switch et {
		case textparse.EntryType:
			continue
		case textparse.EntryHelp:
			continue
		case textparse.EntryUnit:
			continue
		case textparse.EntryComment:
			continue
		default:
		}

		t := time.Now().Unix()
		_, tp, v := p.Series()

		if tp != nil {
			t = *tp
		}

		var (
			lset labels.Labels
			mets string
			hash uint64
		)

		mets = p.Metric(&lset)
		hash = lset.Hash()

		if !lset.Has(labels.MetricName) {
			return nil, errNameLabelMandatory
		}

		m := &Metrics{
			Metrics: mets,
			Name:    lset.Get(labels.MetricName),
			Hash:    hash,
			Labels:  lset,
			Value:   v,
			Time:    t,
		}

		list = append(list, m)
	}

	return list, nil
}
