package web

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"lib"
	d "nb/data"
	"net/http"
)

//
func createBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:createBook start")
	request := new(loginRequest)
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		http.Error(w, "Program Error: loginRequest json decode", http.StatusInternalServerError)
		return
	}
	// add book
	book := &d.Book{
		BookName:   ps.ByName("bookname"),
		AccessCode: request.AccessCode,
	}
	resultChan := d.Data("addBook", book)
	result := <-resultChan
	if result.Status != d.DataOk { // resultVal contains errMsg
		http.Error(w, result.Val.(string), http.StatusBadRequest)
		return
	}
	// add tab
	tab := &d.Tab{
		Bookid:    book.Id,
		TabNumber: 1,
		TabName:   "About",
		Hidden:    false,
	}
	resultChan = d.Data("addTab", tab)
	result = <-resultChan

	// get book tabs (will only be the 1 just added)
	bookTabs := &d.GetBookTabs{Bookid: book.Id} // see data/getbooktabs.go
	resultChan = d.Data("getBookTabs", bookTabs)
	result = <-resultChan

	// build response
	response := &loginResponse{
		Bookid:    book.Id,
		BookName:  book.BookName,
		Tabs:      bookTabs.TabInfoMap,
		Broadcast: getBroadcast(),
	}
	// add Login
	login := &d.Login{Bookid: book.Id, Type: "edit"}
	resultChan = d.Data("addLogin", login)
	result = <-resultChan
	response.Token = login.Token

	jsonResponseWriter(w, r, response)

	lib.Trace(2, "handler:createBook end")
}

func openBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:openBook start")
	request := new(loginRequest)
	if decodeErr := json.NewDecoder(r.Body).Decode(request); decodeErr != nil {
		http.Error(w, "Program Error: loginRequest json decode", http.StatusInternalServerError)
		return
	}
	// get book
	book := &d.Book{
		BookName:   ps.ByName("bookname"),
		AccessCode: request.AccessCode,
	}
	resultChan := d.Data("getBook", book)
	result := <-resultChan
	if result.Status != d.DataOk { // resultVal contains errMsg
		http.Error(w, result.Val.(string), http.StatusBadRequest)
		return
	}
	// get book tabs
	bookTabs := &d.GetBookTabs{Bookid: book.Id} // see data/getbooktabs.go
	resultChan = d.Data("getBookTabs", bookTabs)
	result = <-resultChan

	// build response
	response := &loginResponse{
		Bookid:    book.Id,
		BookName:  book.BookName,
		Tabs:      bookTabs.TabInfoMap,
		Broadcast: getBroadcast(),
	}
	// add Login
	login := &d.Login{Bookid: book.Id, Type: "edit"}
	resultChan = d.Data("addLogin", login)
	result = <-resultChan
	response.Token = login.Token

	jsonResponseWriter(w, r, response)

	lib.Trace(2, "handler:openBook end")
}

func viewBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Trace(2, "handler:viewBook start")
	// get book
	book := &d.Book{BookName: ps.ByName("bookname")}
	resultChan := d.Data("getBook", book)
	result := <-resultChan
	if result.Status != d.DataOk { // resultVal contains errMsg
		http.Error(w, result.Val.(string), http.StatusBadRequest)
		return
	}
	// get book tabs
	bookTabs := &d.GetBookTabs{Bookid: book.Id} // see data/getbooktabs.go
	resultChan = d.Data("getBookTabs", bookTabs)
	result = <-resultChan

	// build response
	response := &loginResponse{
		Bookid:    book.Id,
		BookName:  book.BookName,
		Tabs:      bookTabs.TabInfoMap,
		Broadcast: getBroadcast(),
	}
	// add Login
	login := &d.Login{Bookid: book.Id, Type: "view"}
	resultChan = d.Data("addLogin", login)
	result = <-resultChan
	response.Token = login.Token

	jsonResponseWriter(w, r, response)

	lib.Trace(2, "handler:viewBook end")
}
