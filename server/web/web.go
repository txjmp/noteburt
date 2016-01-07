package web

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"lib"
	"log"
	"nb/common"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var indexHtml []byte // contents of file index.html
// var indexBookName int   // index of literal "BOOK_NAME" in indexHtml
// var indexHtmlLock sync.RWMutex

var codeJSFileName = "codejs.gzip"
var codeJSModTime time.Time

var broadcastMsg string = "broadcast msg"
var broadcastLock sync.RWMutex

// WebStart - setup handlers and start ListenAndServ
func WebStart(httpPort string) {
	lib.Trace(0, "WebStart")
	loadIndexHtml()
	gzipCodeJS()

	router := httprouter.New()

	router.GET("/book/:bookname", returnIndexHtml) // user specified URL
	router.GET("/code.js", returnCodeJS)           // all js code, gzipped

	// handlers in bookhandlers.go
	router.POST("/open/:bookname", openBook)
	router.GET("/view/:bookname", viewBook)
	router.POST("/create/:bookname", createBook)

	// handlers in tabhandlers.go
	router.POST("/tab/:token/:bookid", addTab)
	router.PUT("/tab/:token/:bookid/:tabid", changeTab)
	router.GET("/tabnotes/:token/:bookid/:tabid", getTabNotes)

	// handlers in notehandlers.go
	router.POST("/note/:token/:bookid/:tabid", addNote)
	router.PUT("/note/:token/:bookid/:tabid/:noteid", changeNote)
	router.PUT("/positionnote/:token/:bookid/:tabid/:noteid", positionNote)
	router.DELETE("/note/:token/:bookid/:tabid/:noteid", deleteNote)

	// handlers in specialhandlers.go
	router.GET("/accesscode", getAccessCode)
	router.GET("/shutdown/:keyword", shutdown)

	router.NotFound = http.FileServer(http.Dir("static"))

	lib.Trace(0, "handlers assigned - listening on ", httpPort)

	if httpPort == ":https" {
		log.Fatal(http.ListenAndServeTLS(httpPort, "tls/cert.pem", "tls/key.pem", router))
	} else {
		log.Fatal(http.ListenAndServe(httpPort, router))
	}
}

// handler for /:bookname, returns indexHtml with requested book name included
func returnIndexHtml(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookName := []byte(ps.ByName("bookname"))
	response := bytes.Replace(indexHtml, []byte("BOOK_NAME"), bookName, 1)
	w.Write(response)
}

// handler for /code.js, returns gzipped js code file
func returnCodeJS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Content-Encoding", "gzip")
	codeFile, _ := os.Open(codeJSFileName)
	http.ServeContent(w, r, "", codeJSModTime, codeFile)
	codeFile.Close()
}

// jsonResponseWriter - convert response data to json, gzip if allowed, write to http.ResponseWriter
func jsonResponseWriter(w http.ResponseWriter, r *http.Request, data interface{}) error {
	var err error
	var jsonData []byte
	if jsonData, err = json.Marshal(data); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")

	gzipOk := strings.Contains(fmt.Sprint(r.Header["Accept-Encoding"]), "gzip")
	if gzipOk {
		w.Header().Set("Content-Encoding", "gzip")
		compressor := common.GzipWriterPool.Get().(*gzip.Writer)
		compressor.Reset(w)
		_, err = compressor.Write(jsonData)
		compressor.Close()
		common.GzipWriterPool.Put(compressor)
		lib.Trace(3, "json response gzip")
	} else {
		_, err = w.Write(jsonData)
	}
	if err != nil {
		lib.Trace(0, "jsonResponseWriter error: ", err.Error())
	}
	return err
}

func etagMatch(r *http.Request, serverETag string) bool {
	requestETag := r.Header.Get("If-None-Match")
	lib.Trace(2, "requestEtag:", requestETag, ",  serverETag:", serverETag)
	if requestETag == serverETag {
		return true
	}
	return false
}

/*
func badLogin(w http.ResponseWriter) {
	http.Error(w, "Connect Error: Login Invalid or Expired", http.StatusBadRequest)
}

func badURL(handler string, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Program Error: invalid URL", http.StatusInternalServerError)
	lib.Trace(0, handler, "invalid url: "+r.URL.Path)
}
*/
func getBroadcast() string {
	broadcastLock.RLock()
	msg := broadcastMsg
	broadcastLock.RUnlock()
	return msg
}
func loadIndexHtml() {
	var err error
	indexHtml, err = ioutil.ReadFile("static/index.html")
	if err != nil {
		log.Fatal("web.loadIndexHtml Failed", err)
	}
	/*
		indexBookName = bytes.Index(indexHtml, []byte("BOOK_NAME"))
		if indexBookName < 0 {
			log.Fatal("web.loadIndexHtml BOOK_NAME not found")
		}
	*/
}

// merge and gzip app javascript files into a single file, run once at startup
func gzipCodeJS() {
	codeJSModTime = time.Now()

	fileOut, _ := os.Create(codeJSFileName)
	compressor, err := gzip.NewWriterLevel(fileOut, gzip.BestCompression)
	if err != nil {
		log.Fatal("gzipCodeJS", err)
	}
	files := []string{
		"main.js",
		"hub.js",
		"data.js",
		"screen1.js",
		"view_booktabs.js",
		"view_noteedit.js",
		"view_noteview.js",
		"view_tabmgr.js",
		"view_position.js",
		"lib.js",
	}
	var fileIn *os.File
	for _, fileName := range files {
		fileIn, err = os.Open("static/code/" + fileName)
		if err != nil {
			log.Fatalln("gzipCodeJS error: ", fileName, err)
		}
		io.Copy(compressor, fileIn)
		fileIn.Close()
	}
	compressor.Close()
	fileOut.Close()
}
