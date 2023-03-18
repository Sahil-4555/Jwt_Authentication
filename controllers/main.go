package controllers

type Controllers struct {
}

func NewController() *Controllers {
	return &Controllers{}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}
