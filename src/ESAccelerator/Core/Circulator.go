package Core

import (
	"log"
	"time"
)

type CirculatorResponse struct {
	Error      bool
	StatusCode int

	Response interface{}
}

type Circulator struct {
	MyQueue Queue
}

var __S_Circulator = MakeCirculator(80)

func (this *Circulator) AddESRequestToCirculator(MyESRequest *ESRequest) <-chan CirculatorResponse {
	MyChannel := make(chan CirculatorResponse)

	if MyESRequest == nil {
		go func() {
			MyChannel <- CirculatorResponse{
				true,
				403,
				JSONResponse{
					true,
					"ESRequest 객체를 생성하지 못했습니다."}}
		}()

	} else {
		MyESRequest.Connection.Notifier = func(Response interface{}, Error bool, StatusCode int) {
			MyChannel <- CirculatorResponse{
				Error,
				StatusCode,
				Response}

			close(MyChannel)
		}

		this.MyQueue.Push(MyESRequest, nil)
	}

	return MyChannel
}

func (this *Circulator) DoCirculate(Ticker *time.Ticker) {
	Q := &this.MyQueue

	for {
		<-Ticker.C

		Jobs := Q.MPop(ESTimestamp(1000 * time.Millisecond))

		if Jobs != nil {
			PendedRequests := map[string][]ESRequestBody{}
			Specimens := make([]ESRequestImpl, 0, len(GetRecognizableRequests()))

			for _, V := range *Jobs {
				SpecimenIdx := 0
				Impl, Request := V.GetLinearly()

				Body, Exception := Impl.GetRequestBody(V)

				if Exception != nil {
					log.Printf("ERROR (%s): %s -> %s", Request.RemoteAddr, Request.URL.Path, Exception.Error())

					this.SendResponse(
						V,
						true,
						Exception.Error(),
						400)

					continue
				}

				if len(Specimens) == 0 ||
					!func() bool {
						CompatibleSpecimenWasFound := false

						for Acc, V_ := range Specimens {
							SpecimenIdx = Acc

							if CompatibleSpecimenWasFound = V_.Compatible(Impl); CompatibleSpecimenWasFound {
								break
							}
						}

						if !CompatibleSpecimenWasFound {
							SpecimenIdx++
						}

						return CompatibleSpecimenWasFound
					}() {
					Specimens = append(Specimens, Impl)
				}

				Endpoint := Specimens[SpecimenIdx].Endpoint()
				Box, Exists := PendedRequests[Endpoint]

				if !Exists {
					Box = []ESRequestBody{}
					PendedRequests[Endpoint] = Box
				}

				Box = append(Box, *Body)
				PendedRequests[Endpoint] = Box
			}

			for Acc, _ := range PendedRequests {
				Requests := PendedRequests[Acc]

				if len(Requests) > 0 {
					//Specimen
					Impl, _ := Requests[0].Origin.GetLinearly()

					Impl.DoRequest(this, Requests...)
				}
			}
		}
	}
}

func (this *Circulator) SendResponse(Request *ESRequest, Error bool, Response interface{}, StatusCode int) {
	if Error {
		switch Response.(type) {
		case string:
			{
				Request.Connection.Notifier.(func(interface{}, bool, int))(Response.(string), Error, StatusCode)

				return
			}

		default:
			{
				log.Fatalf("Error 인자가 참인 SendResponse 메소드 호출에 대하여 Response 인자는 항상 문자열 형식이여야만 합니다.")

				return
			}
		}
	}

	Request.Connection.Notifier.(func(interface{}, bool, int))(Response, Error, StatusCode)
}

func MakeCirculator(Interval time.Duration) *Circulator {
	MyCirculator := Circulator{CreateNewQueue()}
	MyTicker := time.NewTicker(Interval * time.Millisecond)

	go MyCirculator.DoCirculate(MyTicker)

	return &MyCirculator
}
