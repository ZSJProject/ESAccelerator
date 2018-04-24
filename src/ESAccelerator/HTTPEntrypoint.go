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

	Notifier interface{}
}

type JSONResponse struct {
	Error   bool
	Message interface{}
}

func (this *HTTPConnection) SendError(Context []byte, StatusCode int) {
	Delegate := this.MyWriter
	Body := this.MyBody

	Delegate.WriteHeader(StatusCode)
	Delegate.Write(Context)

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

	Delegate.WriteHeader(StatusCode)
	Delegate.Write(Stringfied)

	log.Printf("OK (%s): %s -> %d", Body.RemoteAddr, Body.URL.Path, StatusCode)
}

func (this *GlobalHTTPHandler) ServeHTTP(Writer http.ResponseWriter, Body *http.Request) {
	MyConnection := HTTPConnection{ Writer, Body, nil }

	if Body.Method != "POST" {
		MyConnection.SendJSON(
			JSONResponse{
				true,
				fmt.Sprintf("다음 메서드는 지원하지 않습니다: %s", Body.Method)},

			403)

		return
	}

	//log.Printf("통과")

	ESResponse := <-__S_Circulator.AddESRequestToCirculator(CreateESRequest(&MyConnection))

	MyConnection.SendJSON(ESResponse.Response, ESResponse.StatusCode)
}

func OpenHTTPServer(Address string) *http.Server {
	ServerInstance := &http.Server{Addr: Address}

	ServerInstance.SetKeepAlivesEnabled(false)

	http.Handle("/", &GlobalHTTPHandler{})

	go func() {
		if Exception := ServerInstance.ListenAndServe(); Exception != nil {
			log.Printf("HTTP 서버를 활성화 하려던 중 예외가 발생했습니다: %s", Exception)
		}
	}()

	return ServerInstance
}
