package web

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"lib"
	d "nb/data"
	"net/http"
	"time"
)

// 	handle POST(/note/:token/:bookid/:tabid)
func addNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:addNote start")
	var result d.Result
	bookid := ps.ByName("bookid")
	login := &d.Login{Token: ps.ByName("token"), Bookid: bookid, Type: "edit"}
	resultChan := d.Data("auth", login)
	if result = <-resultChan; result.Status != d.DataOk {
		http.Error(w, result.Val.(string), http.StatusInternalServerError)
		return
	}
	request := new(noteRequest)
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		http.Error(w, "Program Error: noteRequest json decode", http.StatusInternalServerError)
		return
	}
	// add note
	note := &d.Note{
		Content:  request.Content,
		When:     time.Now(),
		Html:     request.Html,
		Markdown: request.Markdown,
		Mono:     request.Mono,
	}
	noteParms := &d.NoteParms{
		Bookid: bookid,
		Tabid:  ps.ByName("tabid"),
		Previd: request.Previd,
		Note:   note,
	}
	resultChan = d.Data("addNote", noteParms)
	result = <-resultChan
	w.Write([]byte(note.Id)) // loaded by addNote
	lib.Trace(2, "handler:addNote end")
}

// 	handle PUT(/note/:token/:bookid/:tabid/:noteid)
func changeNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:changeNote start")
	var result d.Result
	bookid := ps.ByName("bookid")
	login := &d.Login{Token: ps.ByName("token"), Bookid: bookid, Type: "edit"}
	resultChan := d.Data("auth", login)
	if result = <-resultChan; result.Status != d.DataOk {
		http.Error(w, result.Val.(string), http.StatusInternalServerError)
		return
	}
	request := new(noteRequest)
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		http.Error(w, "Program Error: noteRequest json decode", http.StatusInternalServerError)
		return
	}
	note := &d.Note{
		Id:       ps.ByName("noteid"),
		Content:  request.Content,
		When:     time.Now(),
		Html:     request.Html,
		Markdown: request.Markdown,
		Mono:     request.Mono,
	}
	noteParms := &d.NoteParms{
		Bookid: bookid,
		Tabid:  ps.ByName("tabid"),
		Note:   note,
	}
	resultChan = d.Data("changeNote", noteParms)
	result = <-resultChan
	w.Write([]byte("Update Successful"))

	lib.Trace(2, "handler:changeNote end")
}

// 	handle DELETE(/note/:token/:bookid/:tabid/:noteid)
func deleteNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:deleteNote start")
	var result d.Result
	bookid := ps.ByName("bookid")
	login := &d.Login{Token: ps.ByName("token"), Bookid: bookid, Type: "edit"}
	resultChan := d.Data("auth", login)
	if result = <-resultChan; result.Status != d.DataOk {
		http.Error(w, result.Val.(string), http.StatusInternalServerError)
		return
	}
	note := &d.Note{Id: ps.ByName("noteid")}
	noteParms := &d.NoteParms{
		Bookid: bookid,
		Tabid:  ps.ByName("tabid"),
		Note:   note,
	}
	resultChan = d.Data("deleteNote", noteParms)
	result = <-resultChan
	w.Write([]byte("Update Successful"))

	lib.Trace(2, "handler:deleteNote end")
}

// handle PUT(/positionnote/:token/:bookid/:tabid/:noteid)
func positionNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:positionNote start")
	var result d.Result
	bookid := ps.ByName("bookid")
	login := &d.Login{Token: ps.ByName("token"), Bookid: bookid, Type: "edit"}
	resultChan := d.Data("auth", login)
	if result = <-resultChan; result.Status != d.DataOk {
		http.Error(w, result.Val.(string), http.StatusInternalServerError)
		return
	}
	request := new(noteRequest)
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		http.Error(w, "Program Error: noteRequest json decode", http.StatusInternalServerError)
		return
	}
	noteParms := &d.NoteParms{
		Bookid: bookid,
		Tabid:  ps.ByName("tabid"),
		Note:   &d.Note{Id: ps.ByName("noteid")},
		Previd: request.Previd,
	}
	resultChan = d.Data("positionNote", noteParms)
	result = <-resultChan
	w.Write([]byte("Update Successful"))

	lib.Trace(2, "handler:positionNote end")
}
