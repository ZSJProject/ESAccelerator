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
}

type JSONResponse struct {
	Error   bool
	Message interface{}
}

func (this *GlobalHTTPHandler) SendError(Connection HTTPConnection, Context []byte, StatusCode int) {
	Delegate := Connection.MyWriter
	Body := Connection.MyBody

	Delegate.WriteHeader(StatusCode)
	Delegate.Write(Context)

	log.Printf("ERR (%s): %s -> %d", Body.RemoteAddr, Body.URL.Path, StatusCode)
}

func (this *GlobalHTTPHandler) SendJSON(Connection HTTPConnection, ToJSON JSONResponse, StatusCode int) {
	Delegate := Connection.MyWriter
	Body := Connection.MyBody

	Stringfied, Exception := json.Marshal(ToJSON)

	if Exception != nil {
		Connection.SendError(
			Connection,
			[]byte(fmt.Sprintf("다음과 같은 타입을 JSON 형식으로 변환할 수 없었습니다: %T", ToJSON.Message)),
			http.StatusInternalServerError)

		return
	}

	Delegate.WriteHeader(StatusCode)
	Delegate.Write(Stringfied)

	log.Printf("OK (%s): %s -> %d", Body.RemoteAddr, Body.URL.Path, StatusCode)
}

func (this *GlobalHTTPHandler) ServeHTTP(Writer http.ResponseWriter, Body *http.Request) {
	if Body.Method != "POST" {
		this.SendJSON(JSONResponse{
			true,
			fmt.Sprintf("다음 메서드는 지원하지 않습니다: %s", Body.Method)},

			403)

		return
	}

	log.Printf("통과")

	ESResponse := <-__S_Circulator.AddESRequestToCirculator(CreateESRequest(this))

	this.SendJSON(JSONResponse{
		ESResponse.Error,
		ESResponse.Response},

		ESResponse.StatusCode)
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
