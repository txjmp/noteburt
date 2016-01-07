package web

import (
	"github.com/julienschmidt/httprouter"
	"lib"
	"nb/common"
	d "nb/data"
	"net/http"
	"os"
	"time"
)

// getAccessCodeHandler - GET(/accesscode)
func getAccessCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accessCode := common.AccessCode()
	w.Write([]byte(accessCode))
}

/*
// broadcastHandler - handles requests with URL = /broadcast/msgcontent
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	urlTypes := []string{"string", "string"}
	var urlVals []interface{}
	var ok bool
	if urlVals, ok = splitURL(r, urlTypes); !ok {
		badURL("broadcastHandler", w, r)
		return
	}
	msg := urlVals[1].(string)
	if msg == "get" {
		broadcastLock.RLock()
		msg = broadcastMsg
		broadcastLock.RUnlock()
		w.Write([]byte(msg))
		return
	}
	Trace(1, "broacast msg =", msg)
	broadcastLock.Lock()
	if msg == "off" {
		broadcastMsg = ""
	} else {
		broadcastMsg = strings.Replace(msg, "_", " ", -1)
	}
	broadcastLock.Unlock()
	w.Write([]byte("broadcast action complete"))
}
*/

// handles requests with URL = /shutdown/:keyword
func shutdown(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if ps.ByName("keyword") != "downboy" {
		http.Error(w, "Sorry Charlie", http.StatusInternalServerError)
		lib.Trace(0, "shutdown - invalid keyword")
		return
	}
	w.Write([]byte("shutdown request received"))
	lib.Trace(0, "shutdown request received")
	resultChan := d.Data("shutdown", &d.SimpleRequest{nil})
	<-resultChan
	time.Sleep(3 * time.Second) // wait for Trace buffer to clear
	os.Exit(0)
}
