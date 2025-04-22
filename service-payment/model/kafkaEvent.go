package model

type KafkaEvent struct {
	Debtor   string
	Amount   float32
	Creditor string
	Event    string
}
