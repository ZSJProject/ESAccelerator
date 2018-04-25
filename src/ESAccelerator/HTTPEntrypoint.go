package ESAccelerator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GlobalHTTPHandler struct {
	Handler http.Handler
}

type HTTPConnection struct {
	MyWriter http.ResponseWriter
	MyBody   *http.Request

	MyFlag	HTTPResponseFlag

	Notifier interface{}
}

type HTTPResponseFlag struct {
	HasBody			bool
	Interrupted		bool
}

type JSONResponse struct {
	Error   bool
	Message interface{}
}

func (this *HTTPConnection) SendError(Context []byte, StatusCode int) {
	Delegate := this.MyWriter
	Body := this.MyBody

	Delegate.WriteHeader(StatusCode)

	if this.MyFlag.HasBody {
		Delegate.Write(Context)
	}

	log.Printf("ERR (%s): %s -> %d", Body.RemoteAddr, Body.URL.Path, StatusCode)
}

func (this *HTTPConnection) SendJSON(ToJSON interface{}, StatusCode int) {
	Delegate := this.MyWriter
	Body := this.MyBody

	Stringfied, Exception := json.Marshal(ToJSON)

	if Exception != nil {
		this.SendError(
			[]byte(fmt.Sprintf("다음과 같은 타입을 JSON 형식으로 변환할 수 없었습니다: %T", ToJSON)),
			http.StatusInternalServerError)

		return
	}

	Delegate.Header().Set("Content-Type", "application/json")
	Delegate.WriteHeader(StatusCode)

	if this.MyFlag.HasBody {
		Delegate.Write(Stringfied)
	}

	log.Printf("OK (%s): %s -> %d", Body.RemoteAddr, Body.URL.Path, StatusCode)
}

func (this *GlobalHTTPHandler) ServeHTTP(Writer http.ResponseWriter, Body *http.Request) {
	MyConnection 	:= HTTPConnection{
		Writer,
		Body,
		HTTPResponseFlag{ false, false },
		nil}

	ESResponse 		:= <-__S_Circulator.AddESRequestToCirculator(CreateESRequest(&MyConnection))
	Flag			:= &MyConnection.MyFlag

	switch Method := Body.Method; Method {
	case "GET":
		fallthrough
	case "POST":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		{
			Flag.HasBody = true

			break
		}

	case "TRACE":
		fallthrough
	case "CONNECT":
		{
			MyConnection.SendError([]byte{}, 405)

			return
		}

	default:
		break
	}

	if !Flag.Interrupted {
		MyConnection.SendJSON(ESResponse.Response, ESResponse.StatusCode)
	}
}

func OpenHTTPServer(Address string) *http.Server {
	ServerInstance := &http.Server{Addr: Address}

	//ServerInstance.SetKeepAlivesEnabled(false)

	http.Handle("/", &GlobalHTTPHandler{})

	go func() {
		if Exception := ServerInstance.ListenAndServe(); Exception != nil {
			log.Printf("HTTP 서버를 활성화 하려던 중 예외가 발생했습니다: %s", Exception)
		}
	}()

	return ServerInstance
}
