package ESAccelerator

import (
	"net/http"
	"strings"
)

type ESRequestImpl interface {
	Name() string
	Endpoint() string

	Compatible(ESRequestImpl) bool
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
	Body := Connection.MyBody
	Path := Body.URL.Path

	OPIdx := strings.LastIndex(Path, "_")

	if OPIdx == -1 {
		return nil
	}

	OP := Path[OPIdx:]
	Impl := func() ESRequestImpl {
		for k, v := range GetRecognizableRequests() {
			if k == OP {
				return v(Connection)
			}
		}

		return nil
	}()

	if Impl == nil {
		return nil
	}

	return &ESRequest{
		Impl,
		Connection}
}
