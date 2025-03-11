package loyalty

import (
	"time"

	"github.com/niksmo/gophermart/pkg/di"
)

type BalanceScheme struct {
	ID         int32     `json:"-"`
	OwnerID    int32     `json:"-"`
	Balance    float64   `json:"current"`
	Withdraw   float64   `json:"withdraw"`
	LastUpdate time.Time `json:"-"`
}

func (bs *BalanceScheme) ScanRow(row di.Row) error {
	return row.Scan(
		&bs.ID, &bs.OwnerID, &bs.Balance, &bs.Withdraw, &bs.LastUpdate,
	)
}
