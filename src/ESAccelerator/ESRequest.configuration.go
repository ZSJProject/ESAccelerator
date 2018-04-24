package ESAccelerator

type ESRequestImplGenerator func(*GlobalHTTPHandler) ESRequestImpl
type ESRequestImplHelper map[string]ESRequestImplGenerator

var (
	RecognizableRequests = ESRequestImplHelper {
		"_search":	func(Delegate *GlobalHTTPHandler) ESRequestImpl {
			return ConvertToESRequestImpl(&ESSearchRequest{Delegate })
		}}
)

func ConvertToESRequestImpl(Implement ESRequestImpl) ESRequestImpl {
	return Implement
}

func GetRecognizableRequests() ESRequestImplHelper {
	return RecognizableRequests
}