package sqlstore

import (
	"github.com/xformation/sdp/pkg/bus"
	m "github.com/xformation/sdp/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetDBHealthQuery)
}

func GetDBHealthQuery(query *m.GetDBHealthQuery) error {
	return x.Ping()
}
