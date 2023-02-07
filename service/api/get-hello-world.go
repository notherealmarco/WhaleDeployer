package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Ciao, sono il backend di Whale Deployer! Che si dice?"))
}
