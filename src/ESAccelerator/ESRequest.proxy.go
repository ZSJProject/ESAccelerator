package ESAccelerator

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
	return true
}

func (__DO_NOT_USE___ *ESProxyRequest) DoRequest(Self *Circulator, Bodies ...ESRequestBody) {
	ES 			:= GetGlobalESConnector()
	ESMeta		:= ES.Metadata.Origin

	for _, V := range Bodies {
		Self.SendResponse(V.Origin, false, ESMeta, 200)
	}
}

func (this *ESProxyRequest) GetRequestBody(Request *ESRequest) (*ESRequestBody, error) {
	Connection 	:= Request.Connection

	return &ESRequestBody{
		Request,
		Connection.MyBody}, nil
}