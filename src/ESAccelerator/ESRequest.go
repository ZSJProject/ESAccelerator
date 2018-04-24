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
	Body *http.Request

	Notifier interface{}
}

type ESRequestBody struct {
	Origin *ESRequest
	Body   interface{}
}

func (this *ESRequest) GetLinearly() (ESRequestImpl, *http.Request) {
	return this.Type, this.Body
}

func CreateESRequest(Delegate *GlobalHTTPHandler) *ESRequest {
	Body := Delegate.MyBody
	Path := Body.URL.Path

	OPIdx := strings.LastIndex(Path, "_")

	if OPIdx == -1 {
		return nil
	}

	OP := Path[OPIdx:]
	Impl := func() ESRequestImpl {
		for k, v := range GetRecognizableRequests() {
			if k == OP {
				return v(Delegate)
			}
		}

		return nil
	}()

	if Impl == nil {
		return nil
	}

	return &ESRequest{
		Impl,
		Delegate.MyBody,
		nil}
}
