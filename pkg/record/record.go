package record

import (
	"context"

	"github.com/snlansky/coral/pkg/application"

	"github.com/snlansky/coral/pkg/logging"
	"github.com/snlansky/coral/pkg/trace"

	"github.com/snlansky/coral/pkg/protos"
	"github.com/snlansky/coral/pkg/xgrpc"
)

var _record *recorder

var logger = logging.MustGetLogger("record")

type recorder struct {
	url    string
	client protos.BehaviorLogClient
}

func NewRecorder(url string) (*recorder, error) {
	client, err := xgrpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	return &recorder{
		url:    url,
		client: protos.NewBehaviorLogClient(client),
	}, nil
}

func (r *recorder) Record(ctx context.Context, record *protos.Record) error {
	_, err := r.client.Recode(ctx, record)
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

		err := _record.Record(ctx, record)
		if err != nil {
			logger.Errorf("record option error: %v", err)
		}
	}()
}
