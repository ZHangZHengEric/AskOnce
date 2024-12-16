package helpers

import (
	"askonce/conf"
	"github.com/xiangtao94/golib/pkg/elastic8"
)

var (
	EsClient *elastic8.ElasticsearchClient
)

func InitElastic() {
	config := conf.WebConf.ElasticSearch
	var err error
	EsClient, err = elastic8.InitESClient(config)
	if err != nil {
		panic("elastic connect error: %v" + err.Error())
	}
}
