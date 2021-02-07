package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/assembla/cony"
	"github.com/snlansky/coral/pkg/gateway"
)

type Consume struct {
	cli     *cony.Client
	handler gateway.Consumer
	cns     *cony.Consumer
}

func NewConsume(url string, handler gateway.Consumer) *Consume {
	// Construct new client with the flag url
	// and default backoff policy
	cli := cony.NewClient(
		cony.URL("amqp://admin:admin@localhost/"),
		cony.Backoff(cony.DefaultBackoff),
	)

	// Declarations
	// The queue name will be supplied by the AMQP server
	que := &cony.Queue{
		Name:       "",
		Durable:    false,
		AutoDelete: true,
	}

	bnd1 := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      GetawayRoute,
	}

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd1),
	})

	// Declare and register a consumer
	cns := cony.NewConsumer(
		que,
		cony.AutoAck(), // Auto sign the deliveries
	)
	cli.Consume(cns)

	c := &Consume{
		cli:     cli,
		handler: handler,
		cns:     cns,
	}
	return c
}

func (c *Consume) Start() {
	for c.cli.Loop() {
		select {
		case msg := <-c.cns.Deliveries():
			switch msg.RoutingKey {
			case GatewayApiRoute:
				var api gateway.Api
				err := json.Unmarshal(msg.Body, &api)
				if err != nil {
					logger.Errorf("json unmarshal api error: %v, body: %q ", err, msg.Body)
					continue
				}
				err = c.handler.ApiHandler(&api)
				if err != nil {
					logger.Errorf("api handler error: %v, api: %v ", err, api)
				}
			case GatewayApiCallRoute:
				var entity gateway.ApiCallEntity
				err := json.Unmarshal(msg.Body, &entity)
				if err != nil {
					logger.Errorf("json unmarshal ApiCallEntity error: %v, body: %q ", err, msg.Body)
					continue
				}
				err = c.handler.ApiCallHandler(&entity)
				if err != nil {
					logger.Errorf("api handler error: %v, entity: %v ", err, entity)
				}
			case GetawayContractCallRoute:
				var entity gateway.ContractCallEntity
				err := json.Unmarshal(msg.Body, &entity)
				if err != nil {
					logger.Errorf("json unmarshal ContractCallEntity error: %v, body: %q ", err, msg.Body)
					continue
				}
				err = c.handler.ContractCallHandler(&entity)
				if err != nil {
					logger.Errorf("api handler error: %v, entity: %v ", err, entity)
				}
			}
			// msg.Ack(false)
			// msg.Nack(false)
			// msg.Reject(false)
		case err := <-c.cns.Errors():
			fmt.Printf("Consumer error: %v\n", err)
		case err := <-c.cli.Errors():
			fmt.Printf("Client error: %v\n", err)
		}
	}
}
