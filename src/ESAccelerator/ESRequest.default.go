package ESAccelerator

type ESDefaultRequest struct { /* no special implementation */ }

func (this *ESDefaultRequest) Name() string {
	return "default"
}

func (this *ESDefaultRequest) Endpoint() string {
	return "/"
}

func (this *ESDefaultRequest) Acceptable(Object *HTTPConnection) bool {
	switch Method := Object.MyBody.Method; Method {
	case "GET":
		fallthrough
	case "PATCH":
		Object.MyFlag.HasBody = true
		fallthrough
	case "OPTIONS":
		fallthrough
	case "HEAD":
		return true

	default:
		return false
	}
}

func (this *ESDefaultRequest) Compatible(Object ESRequestImpl) bool {
	return Object.Endpoint() == this.Endpoint()
}

func (__DO_NOT_USE___ *ESDefaultRequest) DoRequest(Self *Circulator, Bodies ...ESRequestBody) {
	ES 			:= GetGlobalESConnector()
	ESMeta		:= ES.Metadata.Origin

	for _, V := range Bodies {
		Self.SendResponse(V.Origin, false, ESMeta, 200)
	}
}

func (this *ESDefaultRequest) GetRequestBody(Request *ESRequest) (*ESRequestBody, error) {
	Connection 	:= Request.Connection
	Body		:= Connection.MyBody.Method

	return &ESRequestBody{
		Request,
		Body}, nil
}