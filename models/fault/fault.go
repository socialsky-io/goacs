package fault

import "time"

type Fault struct {
	UUID      string    `json:"uuid" db:"uuid"`
	CPEUUID   string    `json:"cpe_uuid" db:"cpe_uuid"`
	Code      string    `json:"code" db:"code"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
