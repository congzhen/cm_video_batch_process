package process

import "github.com/rs/xid"

func GetXid() string {
	return xid.New().String()
}
