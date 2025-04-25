package producer

import (
	"assm/service-todo/ctx"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type TodoProducer struct {
	ServiceContext ctx.ServiceCtx
}

func NewTodoProducer(serviceContext ctx.ServiceCtx) *TodoProducer {
	return &TodoProducer{ServiceContext: serviceContext}
}

func (tp *TodoProducer) Run() {
	for e := range tp.ServiceContext.KafkaProducer().Events() {
		switch ev := (e).(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

func (tp *TodoProducer) Shutdown() func() {
	return func() {
		log.Println("Shutting down Producer server")
		tp.ServiceContext.KafkaProducer().Close()
	}
}
