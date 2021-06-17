package main

type Credentials struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type BuyOrder struct {
	Reference string
	Id        string
	Amount    float64
}

type AutoInvest struct {
	Scheduler string     `validate:"required"`
	Buy       []BuyOrder `validate:"required"`
}
