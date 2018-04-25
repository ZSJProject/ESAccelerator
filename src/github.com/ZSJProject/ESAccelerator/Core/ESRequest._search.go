package Core

import (
	"bytes"
	"errors"
	"gopkg.in/olivere/elastic.v5"
	"strings"
	"encoding/json"
)

type ESSearchRequest struct { /* no special implementation */ }

func (this *ESSearchRequest) Name() string {
	return "_search"
}

func (this *ESSearchRequest) Endpoint() string {
	return "_msearch"
}

func (this *ESSearchRequest) Acceptable(Object *HTTPConnection) bool {
	switch Method := Object.MyBody.Method; Method {
	case "GET":
		fallthrough
	case "POST":
		return true

	default:
		return false
	}
}

func (this *ESSearchRequest) Compatible(Object ESRequestImpl) bool {
	return Object.Endpoint() == this.Endpoint()
}

func (__DO_NOT_USE___ *ESSearchRequest) DoRequest(Self *Circulator, Bodies ...ESRequestBody) {
	ES 			:= GetGlobalESConnector()
	ESClient	:= ES.Client
	ESContext	:= ES.Context

	Requests 	:= make([]*elastic.SearchRequest, 0, len(Bodies))
	Receivers	:= make([]*ESRequest, 0, len(Bodies))
	Send		:= func(IsException bool, StatusCode int, Results interface{}) {
		for Acc, V := range Receivers {
			if IsException {
				Self.SendResponse(V, IsException, Results, StatusCode)

				return
			}

			Self.SendResponse(V, IsException, Results.([]*elastic.SearchResult)[Acc], StatusCode)
		}
	}

	for _, V := range Bodies {
		Requests 	= append(Requests, V.Body.(*elastic.SearchRequest))
		Receivers	= append(Receivers, V.Origin)
	}

	Transition, Exception := ESClient.MultiSearch().
		Add(Requests...).
		Do(ESContext)

	if Exception != nil {
		Send(
			true,
			500,
			"ESSearchRequest 객체에서 병합된 검색 요청을 원본 서버에 보내던 도중 문제가 발생했습니다.")

		return
	}

	Send(
		false,
		200,
		Transition.Responses)
}

func (this *ESSearchRequest) GetRequestBody(Request *ESRequest) (*ESRequestBody, error) {
	//ES 			:= GetGlobalESConnector()
	Connection 	:= Request.Connection

	GivenPath 	:= Connection.MyBody.URL.Path
	Separated 	:= strings.Split(GivenPath, "/")

	Body := elastic.NewSearchRequest()
	Raw, Exception := func() (string, error) {
		_1, _2 := new(bytes.Buffer), new(bytes.Buffer)

		_1.ReadFrom(Connection.MyBody.Body)

		Exception := json.Compact(_2, _1.Bytes())

		return _2.String(), Exception
	}()

	if Exception != nil {
		return nil, errors.New("이스케이프 문자 처리 도중 문제가 발생했습니다.")
	}

	switch len(Separated) {
	case 3:
		{
			//exists index, no type
			Body = Body.Index(Separated[1]).Source(Raw)
			break
		}

	case 4:
		{
			//exists index, type
			Body = Body.Index(Separated[1]).Type(Separated[2]).Source(Raw)
			break
		}

	default:
		return nil, errors.New("인식할 수 없는 경로입니다.")
	}

	return &ESRequestBody{
		Request,
		Body}, nil
}
