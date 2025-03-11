package loyalty

import "time"

type BalanceScheme struct {
	ID         int32     `json:"-"`
	OwnerID    int32     `json:"-"`
	Balance    float64   `json:"current"`
	Withdraw   float64   `json:"withdraw"`
	LastUpdate time.Time `json:"-"`
}
