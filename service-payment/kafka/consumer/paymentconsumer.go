package consumer

import (
	"assm/service-payment/ctx"
	"assm/service-payment/db"
	"assm/service-payment/model"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

type PaymentConsumer struct {
	ServiceContext ctx.ServiceCtx
}

func NewPaymentConsumer(serviceContext ctx.ServiceCtx) *PaymentConsumer {
	return &PaymentConsumer{ServiceContext: serviceContext}
}

func (pc *PaymentConsumer) Run() {
	c := pc.ServiceContext.KafkaConsumer()
	for {
		msg, err := c.ReadMessage(time.Second)

		mongoClient := pc.ServiceContext.MongoClient()
		mongoCollection := pc.ServiceContext.MongoCollection()

		if err == nil {
			var event model.KafkaEvent
			err := bson.Unmarshal(msg.Value, &event)
			if err != nil {
				fmt.Println("Error unmarshalling message: ", err)
			}

			fmt.Println("got this message from kafka: ", event)

			switch event.Event {
			case "Withdrawal":
				if event.Debtor == "" {
					log.Println("Invalid data")
					return
				}

				if !db.EditBalance(mongoClient, mongoCollection, event.Debtor, -event.Amount) {
					log.Println("failed to update balance")
					return
				}
			case "Deposit":
				if event.Amount < 0 || event.Creditor == "" {
					log.Println("Invalid data")
					return
				}
				if !db.EditBalance(mongoClient, mongoCollection, event.Creditor, event.Amount) {
					log.Println("failed to update balance")
					return
				}

			case "P2PTransfer":
				if event.Amount < 0 || event.Creditor == "" || event.Debtor == "" {
					log.Println("Invalid data")
					return
				}

				transferStatus, err := db.Transfer(mongoClient, mongoCollection, event.Debtor, event.Creditor, event.Amount)

				if err != nil {
					log.Println(err)
					return
				}

				if !transferStatus {
					log.Println("failed to update balance")
					return
				}

			default:
				fmt.Println("Got unrecognized event with these values: ", string(msg.Value))
			}

		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

func (pc *PaymentConsumer) Shutdown() {
	log.Println("Shutting down Consumer service")
	err := pc.ServiceContext.KafkaConsumer().Close()
	if err != nil {
		fmt.Println("Error shutting down Consumer service")
		return
	}
}
