package ESAccelerator

import (
	"errors"
	"strings"
	"gopkg.in/olivere/elastic.v5"
	"bytes"
)

type ESSearchRequest struct {
	MyDelegate *GlobalHTTPHandler
}

func (this *ESSearchRequest) Name() string {
	return "_search"
}

func (this *ESSearchRequest) Endpoint() string {
	return "_msearch"
}

func (this *ESSearchRequest) Compatible(Object ESRequestImpl) bool {
	return Object.Endpoint() == this.Endpoint()
}

func (__DO_NOT_USE___ *ESSearchRequest) DoRequest(Self *Circulator, Bodies ...ESRequestBody) {
	for _, V := range Bodies {
		//MyBody := V.Body.(*elastic.SearchRequest)

		//log.Print(MyBody)

		Self.SendResponse(V.Origin, false, "히힛 완료!", 200)
	}
}

func (this *ESSearchRequest) GetRequestBody(Request *ESRequest) (*ESRequestBody, error) {
	//ES 			:= GetGlobalESConnector()

	GivenPath 	:= Request.Body.URL.Path
	Separated 	:= strings.Split(GivenPath, "/")

	Body		:= elastic.NewSearchRequest()
	Buffer		:= new(bytes.Buffer)
	Raw			:= func() string {
		Buffer.ReadFrom(Request.Body.Body)

		return Buffer.String()
	}()

	switch len(Separated) {
	case 3:
		{
			//exists index, no type
			Body = Body.Index(Separated[0]).Source(Raw)
			break
		}

	case 4:
		{
			//exists index, type
			Body = Body.Index(Separated[0]).Type(Separated[1]).Source(Raw)
			break
		}

	default:
		return nil, errors.New("인식할 수 없는 경로입니다.")
	}

	return &ESRequestBody{
		Request,
		Body}, nil
}