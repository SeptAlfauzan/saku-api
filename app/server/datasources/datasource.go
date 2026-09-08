package datasources

import (
	"github.com/septalfauzan/saku-api/app/server/config"
	"github.com/septalfauzan/saku-api/app/server/datasources/remote"
)

type Datasources struct {
	Remote *remote.RemoteDatasource
}

func NewDatasources(cfg *config.Config) *Datasources {
	return &Datasources{
		Remote: remote.NewRemoteDatasource(cfg),
	}
}
