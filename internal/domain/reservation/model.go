package reservation

import "time"

type Reservation struct {
	Id        int64
	Sku       uint64
	Count     uint32
	CreatedAt time.Time
}
