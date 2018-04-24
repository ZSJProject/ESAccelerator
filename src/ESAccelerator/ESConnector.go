package ESAccelerator

import (
	"context"
	"gopkg.in/olivere/elastic.v5"
)

type ESMetadata struct {
	ESEndpoint      string
	ESVersion       string
	ESBootstrapCode int
}

type ESConnection struct {
	Context context.Context
	Client  *elastic.Client

	Metadata *ESMetadata
}

var __S_ESMetadata = &ESMetadata{ESEndpoint: "http://es.vm.zsj.co.kr:9200"}
var __S_ESConnection = MakeESConnection()

func MakeESConnection() *ESConnection {
	ESContext := context.Background()
	ESClient, ESException := elastic.NewClient(elastic.SetURL(__S_ESMetadata.ESEndpoint))

	if ESException != nil {
		panic(ESException)
	}

	{
		Info, Code, Exception := ESClient.Ping(__S_ESMetadata.ESEndpoint).Do(ESContext)

		if Exception != nil {
			panic(Exception)
		}

		__S_ESMetadata.ESVersion = Info.Version.Number
		__S_ESMetadata.ESBootstrapCode = Code
	}

	return &ESConnection{
		ESContext,
		ESClient,
		__S_ESMetadata}
}

func GetGlobalESConnector() *ESConnection {
	return __S_ESConnection
}
