package rabbitmq

import (
	"encoding/json"
	"fmt"

	gateway2 "github.com/isnlan/coral/pkg/blink/gateway"

	"github.com/isnlan/coral/pkg/logging"

	"github.com/streadway/amqp"

	"github.com/assembla/cony"
)

const GatewayApiRoute = "gateway.api"
const GatewayApiCallRoute = "gateway.api.call"
const GetawayContractCallRoute = "gateway.contract.call"
const GetawayRoute = "gateway.#"

var exc = cony.Exchange{
	Name:       "bink.gateway",
	Kind:       "topic",
	Durable:    true,
	AutoDelete: false,
}

var logger = logging.MustGetLogger("rabbitmq")

type produceImpl struct {
	cli                  *cony.Client
	apiProducer          *cony.Publisher
	apiCallProducer      *cony.Publisher
	contractCallProducer *cony.Publisher
}

func NewProduce(url string) gateway2.Producer {
	// Construct new client with the flag url
	// and default backoff policy
	cli := cony.NewClient(
		cony.URL(url),
		cony.Backoff(cony.DefaultBackoff),
	)

	// Declare the exchange we'll be using
	cli.Declare([]cony.Declaration{
		cony.DeclareExchange(exc),
	})

	p := &produceImpl{cli: cli}
	p.setup()

	return p
}

func (p *produceImpl) setup() {
	// Declare and register a publisher
	// with the cony client

	apiProducer := cony.NewPublisher(exc.Name, GatewayApiRoute)
	p.cli.Publish(apiProducer)
	p.apiProducer = apiProducer

	apiCallProducer := cony.NewPublisher(exc.Name, GatewayApiCallRoute)
	p.cli.Publish(apiCallProducer)
	p.apiCallProducer = apiCallProducer

	contractCallProducer := cony.NewPublisher(exc.Name, GetawayContractCallRoute)
	p.cli.Publish(contractCallProducer)
	p.contractCallProducer = contractCallProducer

	// Client loop sends out declarations(exchanges, queues, bindings
	// etc) to the AMQP server. It also handles reconnecting.
	go func() {
		for p.cli.Loop() {
			select {
			case err := <-p.cli.Errors():
				fmt.Printf("Client error: %v\n", err)
			case blocked := <-p.cli.Blocking():
				fmt.Printf("Client is blocked %v\n", blocked)
			}
		}
	}()
}

func (p *produceImpl) ApiUpload(api *gateway2.Api) error {
	bytes, err := json.Marshal(api)
	if err != nil {
		return err
	}
	return p.apiProducer.Publish(amqp.Publishing{
		Body: bytes,
	})
}

func (p *produceImpl) ApiCallRecord(entity *gateway2.ApiCallEntity) error {
	bytes, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return p.apiCallProducer.Publish(amqp.Publishing{
		Body: bytes,
	})
}

func (p *produceImpl) ContractCallRecord(entity *gateway2.ContractCallEntity) error {
	bytes, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return p.contractCallProducer.Publish(amqp.Publishing{
		Body: bytes,
	})
}
