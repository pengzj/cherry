package cherry

import (
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

type handlerInfo struct {
	method string
	pattern string
	handler handlerFunc
	afterFunctions []handlerFunc
}

var routes = make(map[string][]handlerInfo)


func Get(pattern string, handler func(http.ResponseWriter,  *http.Request), afterFunctions...func(http.ResponseWriter, *http.Request))  {
	handleMethod("GET", pattern, handler, afterFunctions...)
}

func Post(pattern string, handler func(http.ResponseWriter,  *http.Request), afterFunctions...func(http.ResponseWriter, *http.Request))  {
	handleMethod("POST", pattern, handler, afterFunctions...)
}

func Put(pattern string, handler func(http.ResponseWriter,  *http.Request), afterFunctions...func(http.ResponseWriter, *http.Request))  {
	handleMethod("PUT", pattern, handler, afterFunctions...)
}

func Delete(pattern string, handler func(http.ResponseWriter,  *http.Request), afterFunctions...func(http.ResponseWriter, *http.Request))  {
	handleMethod("DELETE", pattern, handler, afterFunctions...)
}

func Patch(pattern string, handler func(http.ResponseWriter,  *http.Request), afterFunctions...func(http.ResponseWriter, *http.Request))  {
	handleMethod("PATCH", pattern, handler, afterFunctions...)
}


func handleMethod(method string, pattern string, handler func(http.ResponseWriter,  *http.Request), afterFunctions...func(http.ResponseWriter, *http.Request))  {
	var afterFuncs []handlerFunc
	for _, a := range afterFunctions {
		afterFuncs = append(afterFuncs, a)
	}

	routes[pattern] = append(routes[pattern], handlerInfo{
		method,
		pattern,
		handler,
		afterFuncs,
	})
}

func startupRoute()  {
	for pattern, v :=range routes {
		http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
			isMatch := false
			for _, info := range v {
				if request.Method == info.method {
					isMatch = true
					middles := getMiddles()
					var hasFinished = false
					for _, middle := range middles {
						if middle(writer, request) == false {
							hasFinished = true
							break
						}
					}
					if hasFinished == true {
						return
					}
					info.handler(writer, request)

					for _, handle := range info.afterFunctions {
						handle(writer, request)
					}
				}
			}
			if isMatch == false {
				writer.WriteHeader(http.StatusMethodNotAllowed)
			}
		})
	}
}



