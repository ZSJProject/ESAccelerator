package ESAccelerator

type ESRequestImplGenerator func(connection *HTTPConnection) ESRequestImpl
type ESRequestImplHelper map[string]ESRequestImplGenerator

var (
	RecognizableRequests = ESRequestImplHelper{
		"_search": func(Delegate *HTTPConnection) ESRequestImpl {
			return ConvertToESRequestImpl(&ESSearchRequest{})
		}}
)

func ConvertToESRequestImpl(Implement ESRequestImpl) ESRequestImpl {
	return Implement
}

func GetRecognizableRequests() ESRequestImplHelper {
	return RecognizableRequests
}
