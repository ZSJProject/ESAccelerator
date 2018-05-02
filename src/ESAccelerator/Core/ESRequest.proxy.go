package Core

import (
	"net/url"
	"net/http/httputil"
	"log"
	"strings"
)

const ESProxyDestination = "http://es.vm.zsj.co.kr:9200"

type ESProxyRequest struct { /* no special implementation */ }

func (this *ESProxyRequest) Name() string {
	return "proxy"
}

func (this *ESProxyRequest) Endpoint() string {
	return "-proxy"
}

func (this *ESProxyRequest) Acceptable(Object *HTTPConnection) bool {
	return true
}

func (this *ESProxyRequest) Compatible(Object ESRequestImpl) bool {
	return Object.Name() == this.Name()
}

func (__DO_NOT_USE___ *ESProxyRequest) DoRequest(Self *Circulator, Bodies ...ESRequestBody) {
	for _, V := range Bodies {
		Connection		:= V.Origin.Connection
		OriginRequest	:= Connection.MyBody
		Flag			:= &Connection.MyFlag

		Flag.Interrupted = true

		//Reverse proxy
		V.Body.(*httputil.ReverseProxy).ServeHTTP(Connection.MyWriter, Connection.MyBody)

		Self.SendResponse(V.Origin, false, "__INTERRUPTED__", 999)

		log.Printf("PROXY (%s): %s -> %s",
			OriginRequest.RemoteAddr,
			OriginRequest.URL.Path,
			strings.Join([]string{ ESProxyDestination, OriginRequest.URL.Path }, ""))
	}
}

func (this *ESProxyRequest) GetRequestBody(Request *ESRequest) (*ESRequestBody, error) {
	DestURL, _	:= url.Parse(ESProxyDestination)

	return &ESRequestBody{
		Request,
		httputil.NewSingleHostReverseProxy(DestURL)}, nil
}