package units

import (
	"fmt"
)

type Unit struct {
	Ip           string
	Path         string
	Method       string
	Request_uuid string
	Time         float64
}

func (u *Unit) Print() {
	fmt.Printf("ip: %s path: %s method: %s request_uuid: %s time: %d \r\n", u.Ip, u.Path, u.Method, u.Request_uuid, u.Time)
}
