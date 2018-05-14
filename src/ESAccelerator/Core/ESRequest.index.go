package Core

import (
	"strings"
	"gopkg.in/olivere/elastic.v5"
	"bytes"
	"encoding/json"
	"errors"
)

type ESIndexRequest struct { /* no special implementation */ }

func (this *ESIndexRequest) Name() string {
	return "index"
}

func (this *ESIndexRequest) Endpoint() string {
	return "-"
}

func (this *ESIndexRequest) Acceptable(Object *HTTPConnection) bool {
	switch Method := Object.MyBody.Method; Method {
	case "POST":
		fallthrough
	case "PUT":
		return true

	default:
		return false
	}
}

func (this *ESIndexRequest) Compatible(Object ESRequestImpl) bool {
	return Object.Endpoint() == this.Endpoint()
}

func (__DO_NOT_USE___ *ESIndexRequest) DoRequest(Self *Circulator, Bodies ...ESRequestBody) {
	ES 			:= GetGlobalESConnector()
	ESClient	:= ES.Client
	ESContext	:= ES.Context

	Requests 	:= make([]*elastic.BulkIndexRequest, 0, len(Bodies))
	Receivers	:= make([]*ESRequest, 0, len(Bodies))
	Send		:= func(IsException bool, StatusCode int, Results interface{}) {
		for Acc, V := range Receivers {
			if IsException {
				Self.SendResponse(V, IsException, Results, StatusCode)

				return
			}

			Response := Results.([]map[string]*elastic.BulkResponseItem)[Acc]["index"]

			Self.SendResponse(V, IsException, Response, Response.Status)
		}
	}

	for _, V := range Bodies {
		Requests 	= append(Requests, V.Body.(*elastic.BulkIndexRequest))
		Receivers	= append(Receivers, V.Origin)
	}

	Bulk := ESClient.Bulk()

	for _, V := range Requests {
		Bulk.Add(V)
	}

	Transition, Exception := Bulk.Do(ESContext)

	if Exception != nil {
		Send(
			true,
			500,
			"ESIndexRequest 객체에서 병합된 인덱스 요청을 원본 서버에 보내던 도중 문제가 발생했습니다.")

		return
	}

	Send(
		false,
		200,
		Transition.Items)
}

func (this *ESIndexRequest) GetRequestBody(Request *ESRequest) (*ESRequestBody, error) {
	Connection 		:= Request.Connection

	GivenPath		:= Connection.MyBody.URL.Path
	Separated		:= strings.Split(GivenPath, "/")

	Body			:= elastic.NewBulkIndexRequest()
	Raw, Exception 	:= func() (string, error) {
		_1, _2 := new(bytes.Buffer), new(bytes.Buffer)

		_1.ReadFrom(Connection.MyBody.Body)

		Exception := json.Compact(_2, _1.Bytes())

		return _2.String(), Exception
	}()

	if Exception != nil {
		return nil, errors.New("이스케이프 문자 처리 도중 문제가 발생했습니다.")
	}

	switch len(Separated) {
	case 2:
		{
			//exists index, no type, no id
			Body = Body.Index(Separated[1]).Doc(Raw)
			break
		}

	case 3:
		{
			//exists index, type, no id
			Body = Body.Index(Separated[1]).Type(Separated[2]).Doc(Raw)
			break
		}

	case 4:
		{
			Body = Body.Index(Separated[1]).Type(Separated[2]).Doc(Raw).Id(Separated[3])
			break
		}

	default:
		return nil, errors.New("인식할 수 없는 경로입니다.")
	}

	return &ESRequestBody{
		Request,
		Body}, nil
}