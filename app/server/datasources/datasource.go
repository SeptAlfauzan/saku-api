package datasources

import (
	"github.com/septalfauzan/saku-api/app/server/datasources/remote"
)

type Datasources struct {
	Remote *remote.RemoteDatasource
}

func NewDatasources(apiKey, apiURL, ocrPrompt, ocrModel string) *Datasources {
	return &Datasources{
		Remote: remote.NewRemoteDatasource(apiKey, apiURL, ocrPrompt, ocrModel),
	}
}
