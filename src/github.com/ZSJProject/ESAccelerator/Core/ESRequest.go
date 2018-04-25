package Core

import (
	"net/http"
	"strings"
)

type ESRequestImpl interface {
	Name() string
	Endpoint() string

	Compatible(ESRequestImpl) bool
	Acceptable(*HTTPConnection) bool

	GetRequestBody(*ESRequest) (*ESRequestBody, error)

	DoRequest(*Circulator, ...ESRequestBody)
}

type ESRequest struct {
	Type ESRequestImpl
	Connection *HTTPConnection
}

type ESRequestBody struct {
	Origin *ESRequest
	Body   interface{}
}

func (this *ESRequest) GetLinearly() (ESRequestImpl, *http.Request) {
	return this.Type, this.Connection.MyBody
}

func CreateESRequest(Connection *HTTPConnection) *ESRequest {
	Body 	:= Connection.MyBody
	Path 	:= Body.URL.Path

	OPIdx 	:= strings.LastIndex(Path, "_")
	OP		:= func() string {
		if OPIdx == -1 {
			return Path
		}

		return Path[OPIdx:]
	}()

	Impl := func() ESRequestImpl {
		for _, V := range GetRecognizableRequests() {
			switch V.Hint.(type) {
			case string:
				if V.Hint.(string) == OP {
					return V.Generator(Connection)
				}

				break

			case func(*HTTPConnection) bool:
				if V.Hint.(func(*HTTPConnection) bool)(Connection) {
					return V.Generator(Connection)
				}
			}
		}

		return nil
	}()

	if Impl == nil || !Impl.Acceptable(Connection) {
		Impl =  ConvertToESRequestImpl(&ESProxyRequest{})
	}

	return &ESRequest{
		Impl,
		Connection}
}
