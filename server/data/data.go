package data

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"lib"
	"log"
)

type DataStatus int

const (
	DataOk DataStatus = iota
	DataNotFound
	DataDuplicate
	DataError
	DataWait
	DataFail
	DataNotChanged
)
const errFatal = true
const errNotFatal = false

type Result struct {
	Status DataStatus
	Val    interface{}
}

var bookMap map[string]*Book // contains 1 entry for each book in db, key=bookid

type tabMap map[string]*Tab    // key = tabid, each tabMap contains all tabs for a book
var bookTabs map[string]tabMap // key = bookid

// sequentialid - used for recid for all rec types, converted to 12 digit string
//	 see getNextid() below
//   initialized in startupDataLoad() below
//   every database add will execute bumpSequentialid func (in db.go) to keep
//		"control" bkt sequence value in sync with sequentialid
var sequentialid uint64

type request struct {
	action     string
	data       processor
	resultChan chan Result
}

var requestChan chan request

type processor interface {
	process(action string, resultChan chan Result)
}

type SimpleRequest struct {
	RequestInput interface{}
}

func (this *SimpleRequest) process(action string, resultChan chan Result) {
	switch action {
	case "shutdown":
		dbWrite("shutdowndb", this, resultChan) // clear db write buffer

	case "shutdowndb":
		db.Close()
		lib.Trace(0, "db closed")
		resultChan <- Result{DataOk, nil}
		log.Println("shutdown done")
	}
}

func DataStart(dbName string) {
	lib.Trace(0, "Data Start Begin")

	dbStart(dbName) // see db.go

	bookMap = make(map[string]*Book)
	bookTabs = make(map[string]tabMap)

	startupDataLoad()

	requestChan = make(chan request, 100)

	go dataDispatch()

	lib.Trace(0, "Data Start Complete")
}

func dataDispatch() {
	var req request
	for {
		req = <-requestChan
		req.data.process(req.action, req.resultChan)
		lib.Trace(2, "data request done")
		if req.action == "shutdown" {
			lib.Trace(0, "data dispatch stopped")
			break
		}
	}
}

// Data - called by web handlers to access data
//	places data request on requestChan which is processed by dataDispatch goroutine (see above)
//  requester receives unique result channel
//	data.Processor method is sent this channel which places result on it
func Data(action string, data processor) chan Result {
	lib.Trace(2, "data request received", action)
	resultChan := make(chan Result)
	requestChan <- request{action, data, resultChan}
	return resultChan
}

func startupDataLoad() {
	lib.Trace(0, "startupDataLoad Begin")
	db.Update(func(tx *bolt.Tx) error {
		// load all book docs into bookMap
		bktBooks := openBooks(tx)
		cursor := bktBooks.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			book := new(Book)
			if err := json.Unmarshal(v, book); err != nil {
				log.Fatal("startupDataLoad json Unmarshal failed", err)
			}
			lib.Trace(0, book.BookName, book.AccessCode)
			bookMap[book.Id] = book
		}
		bktControl := openControl(tx)
		sequentialid, _ = bktControl.NextSequence()
		return nil
	})
	lib.Trace(0, "startupDataLoad End, sequentialid =", sequentialid)
}

// if multiple data goroutines are used, sync/atomic  will be needed to use getNextid() (uses pkg var)
func getNextid() string {
	sequentialid++
	return fmt.Sprintf("%012d", sequentialid)
}
