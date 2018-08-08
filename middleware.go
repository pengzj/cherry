package cherry

import "net/http"

type MiddleWareFunc func(http.ResponseWriter, *http.Request) bool

var middleWares []MiddleWareFunc

func Middleware(middleWare MiddleWareFunc)  {
	middleWares = append(middleWares, middleWare)
}

func getMiddles() []MiddleWareFunc {
	return middleWares
}
