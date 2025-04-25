package model

type HttpResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

type Transaction struct {
	Amount float32 `json:"amount" binding:"required,gt=0"`
	To     string  `json:"to"`
}

type TransactionResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
