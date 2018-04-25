package ESAccelerator

import "strings"


type ESRequestImplGenerator func(connection *HTTPConnection) ESRequestImpl
type ESRequestImplGeneratorWrapper struct {
	Hint		interface{}
	Generator	ESRequestImplGenerator
}

type ESRequestImplHelper []ESRequestImplGeneratorWrapper

var (
	RecognizableRequests = ESRequestImplHelper{
		ESRequestImplGeneratorWrapper{
			"/",
			func(Delegate *HTTPConnection) ESRequestImpl {
				return ConvertToESRequestImpl(&ESDefaultRequest{})
			}},

		ESRequestImplGeneratorWrapper{
			"_search",
			func(Delegate *HTTPConnection) ESRequestImpl {
				return ConvertToESRequestImpl(&ESSearchRequest{})
			}},

		ESRequestImplGeneratorWrapper{
			func(Connection *HTTPConnection) bool {
				Path	:= Connection.MyBody.URL.Path

				IsNotOP := func() bool {
					return strings.Index(Path, "_") == -1
				}()

				lSlash	:= func() int {
					return len(strings.Split(Path, "/"))
				}()

				if IsNotOP && lSlash <= 4 {
					switch Method := Connection.MyBody.Method; Method {
					case "POST":
						fallthrough
					case "PUT":
						return true
					}
				}

				return false
			},

			func(Delegate *HTTPConnection) ESRequestImpl {
				return ConvertToESRequestImpl(&ESIndexRequest{})
			}}}
)

func ConvertToESRequestImpl(Implement ESRequestImpl) ESRequestImpl {
	return Implement
}

func GetRecognizableRequests() ESRequestImplHelper {
	return RecognizableRequests
}
