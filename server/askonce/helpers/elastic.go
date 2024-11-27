package helpers

import (
	"askonce/conf"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"net/url"
	"os"
)

var (
	EsClient *elasticsearch.Client
)

func InitElastic() {
	config := conf.WebConf.ElasticSearch
	endpointUrl, err := url.Parse(config.Addr)
	if err != nil {
		panic("init elastic failed, err:" + err.Error())
	}
	cfg := elasticsearch.Config{
		Addresses: []string{endpointUrl.String()},
		Username:  config.Username,
		Password:  config.Password,
		Logger: &elastictransport.JSONLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	}
	EsClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic("init elastic failed, err:" + err.Error())
	}
}
