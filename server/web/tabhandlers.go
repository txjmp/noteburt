package web

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"lib"
	d "nb/data"
	"net/http"
)

// 	handle POST(/tab/:token/:bookid)
func addTab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:addTab start")
	var result d.Result
	bookid := ps.ByName("bookid")
	login := &d.Login{Token: ps.ByName("token"), Bookid: bookid, Type: "edit"}
	resultChan := d.Data("auth", login)
	if result = <-resultChan; result.Status != d.DataOk {
		http.Error(w, result.Val.(string), http.StatusBadRequest)
		return
	}
	request := new(tabRequest)
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		http.Error(w, "Program Error: tabRequest json decode", http.StatusInternalServerError)
		return
	}
	// add tab
	tab := &d.Tab{
		Bookid:    bookid,
		TabNumber: request.TabNumber,
		TabName:   request.TabName,
		Hidden:    request.Hidden,
	}
	resultChan = d.Data("addTab", tab)
	result = <-resultChan
	w.Write([]byte(tab.Id))

	lib.Trace(2, "handler:addTab end")
}

// 	handle PUT(/tab/:token/:bookid/:tabid)
func changeTab(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:changeTab start")
	var result d.Result
	bookid := ps.ByName("bookid")
	login := &d.Login{Token: ps.ByName("token"), Bookid: bookid, Type: "edit"}
	resultChan := d.Data("auth", login)
	if result = <-resultChan; result.Status != d.DataOk {
		http.Error(w, result.Val.(string), http.StatusBadRequest)
		return
	}
	request := new(tabRequest)
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		http.Error(w, "Program Error: tabRequest json decode", http.StatusInternalServerError)
		return
	}
	// change tab
	tab := &d.Tab{
		Id:        ps.ByName("tabid"),
		Bookid:    bookid,
		TabNumber: request.TabNumber,
		TabName:   request.TabName,
		Hidden:    request.Hidden,
	}
	resultChan = d.Data("changeTab", tab)
	result = <-resultChan
	w.Write([]byte("Update Successful"))

	lib.Trace(2, "handler:changeTab end")
}

// handle GET(/tabnotes/:token/:bookid/:tabid)
func getTabNotes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:getTabNotes start")
	var result d.Result
	bookid := ps.ByName("bookid")
	login := &d.Login{Token: ps.ByName("token"), Bookid: bookid, Type: "view"}
	resultChan := d.Data("auth", login)
	if result = <-resultChan; result.Status != d.DataOk {
		http.Error(w, result.Val.(string), http.StatusInternalServerError)
		return
	}
	// get tab notes
	tabNotes := &d.GetTabNotes{
		Bookid: bookid,
		Tabid:  ps.ByName("tabid"),
		ETag:   r.Header.Get("If-None-Match"),
	} // see data/gettabnotes.go
	lib.Trace(1, "Request ETag =", tabNotes.ETag)
	resultChan = d.Data("getTabNotes", tabNotes)
	result = <-resultChan
	if result.Status == d.DataNotChanged {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.Header().Set("Etag", tabNotes.ETag)
		jsonResponseWriter(w, r, tabNotes.Notes) // Notes = map[string]*d.NoteResponseRec
	}
	lib.Trace(2, "handler:getTabNotes end", len(tabNotes.Notes))
}
