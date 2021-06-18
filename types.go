package shresiesbot

type Credentials struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type BuyOrder struct {
	Reference string  `validate:"required"`
	Id        string  `validate:"required"`
	Amount    float64 `validate:"required"`
}

type BuyConfiguration struct {
	Scheduler string     `validate:"required"`
	Orders    []BuyOrder `validate:"required"`
}

type AutoInvest struct {
	Sharesies *Credentials
	Buy       *BuyConfiguration
}
