package data

/*
--- Database Structure ------------------------------------------------------------------------------
1.0 books - root bkt, all Book docs, key:seqid, val:Book
2.0 book_xxx - root bkt, no docs, sub bkts only
	2.1 tabs - sub bkt, all tab docs for book, key:seqid, val:Tab
	2.2 tab_xxx - sub bkt, no docs, 1 bkt for each tab
		2.2.1 notes - all notes for tab_xxx, key:seqid, val:Note
		2.2.2 previds - key:noteid, val:previd  (id of note preceding this note in sorted order)
_____________________________________________________________________________________________________
*/

import (
	"github.com/boltdb/bolt"
	"lib"
	"log"
	"strconv"
)

const Zeroid = "000000000000"

type bs []byte // byte slice type used for value conversions

var dbWriteChan chan request

var db *bolt.DB

func dbStart(dbFile string) {
	lib.Trace(0, "dbStart Begin")
	var err error
	db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal("db open fail", err)
	}
	dbWriteChan = make(chan request, 100)

	go dbWriteDispatch()

	// if 1st run for db, following will create bkts for startup
	db.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists(bs("books"))
		_, err2 := tx.CreateBucketIfNotExists(bs("control"))
		return lib.CheckErrs(errFatal, "db startup error", err1, err2)
	})
	lib.Trace(0, "dbStart Complete")
}

func dbWriteDispatch() {
	var req request
	for {
		req = <-dbWriteChan
		lib.Trace(2, "dbWrite Request", req.action)
		req.data.process(req.action, req.resultChan)
		lib.Trace(2, "dbWrite Request done")
		if req.action == "shutdowndb" {
			lib.Trace(0, "dbWrite stopped")
			break
		}
	}
}

// dbWrite func - places request on dbWriteChan
//  resultChan is included in parms (unlike Data(...)) because a Data request
//   may be invoking dbWrite and it already has a resultChan that should be used
//   to send response back to caller
func dbWrite(action string, data processor, resultChan chan Result) {
	dbWriteChan <- request{action, data, resultChan}
}

func bumpSequentialid(tx *bolt.Tx) {
	bktControl := openControl(tx)
	bktControl.NextSequence()
}

func DBBkup(path string) {
	db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile(path, 0600)
		if err != nil {
			lib.Trace(0, "*** ERROR *** DB Bkup Error: ", err.Error())
		}
		return err
	})
}

func bytesToInt(in []byte) int {
	a := string(in)
	i, _ := strconv.Atoi(a)
	return i
}
func bytesToInt64(in []byte) int64 {
	a := string(in)
	i, _ := strconv.ParseInt(a, 10, 64)
	return i
}
func intToBytes(in int) []byte {
	a := strconv.Itoa(in)
	return []byte(a)
}
func int64ToBytes(in int64) []byte {
	a := strconv.FormatInt(in, 10)
	return []byte(a)
}

func openControl(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(bs("control"))
}
func openBooks(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket(bs("books"))
}
func openBook(tx *bolt.Tx, bookid string) *bolt.Bucket {
	return tx.Bucket(bs("book_" + bookid))
}
func openTabs(bktBook *bolt.Bucket) *bolt.Bucket {
	return bktBook.Bucket(bs("tabs"))
}
func openTab(bktBook *bolt.Bucket, tabid string) *bolt.Bucket {
	return bktBook.Bucket(bs("tab_" + tabid))
}
func openNotes(bktTab *bolt.Bucket) *bolt.Bucket {
	return bktTab.Bucket(bs("notes"))
}
func openPrevids(bktTab *bolt.Bucket) *bolt.Bucket {
	return bktTab.Bucket(bs("previds"))
}
func openTabAll(bktBook *bolt.Bucket, tabid string) (*bolt.Bucket, *bolt.Bucket, *bolt.Bucket) {
	bktTab := bktBook.Bucket(bs("tab_" + tabid))
	bktNotes := bktTab.Bucket(bs("notes"))
	bktPrevids := bktTab.Bucket(bs("previds"))
	return bktTab, bktNotes, bktPrevids
}
