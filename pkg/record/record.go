package record

import (
	"context"
	"time"

	"github.com/snlansky/coral/pkg/application"

	"github.com/snlansky/coral/pkg/logging"
	"github.com/snlansky/coral/pkg/trace"

	"github.com/snlansky/coral/pkg/protos"
	"github.com/snlansky/coral/pkg/xgrpc"
)

var _record *recorder

var logger = logging.MustGetLogger("record")

type recorder struct {
	url string
	cli *xgrpc.Client
}

func NewRecorder(url string) (*recorder, error) {
	client, err := xgrpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	return &recorder{
		url: url,
		cli: client,
	}, nil
}

func (r *recorder) Record(record *protos.Record) error {
	conn, err := r.cli.Get()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	client := protos.NewBehaviorLogClient(conn.ClientConn)
	_, err = client.Recode(ctx, record)
	return err
}

func InitRecorder(url string) (err error) {
	if url == "" {
		logger.Warn("trace service uninstall...")
		return nil
	}
	_record, err = NewRecorder(url)
	return
}

func AsyncRecode(ctx context.Context, record *protos.Record) {
	go func() {
		if _record == nil {
			return
		}

		record.Url = trace.GetUrlFromContext(ctx)
		record.TraceId = trace.GetTraceIDFromContext(ctx)
		record.Service = application.Name

		err := _record.Record(record)
		if err != nil {
			logger.Errorf("record option error: %v", err)
		}
	}()
}
