package cherry

import (
	"log"
	"net/http"
)

func Run(addr string)  {
	startupRoute()


	log.Fatal(http.ListenAndServe(addr, nil))
}

