package data

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"lib"
	"log"
	"time"
)

type Tab struct {
	Id             string
	Bookid         string
	TabNumber      int
	TabName        string
	LastAccessTime time.Time // updated on any access (read,write) to any note in tab
	ETag           string    // updated on any change to a note
}

func (this *Tab) process(action string, resultChan chan Result) {
	switch action {

	case "addTab":
		this.Id = getNextid()
		this.LastAccessTime = time.Now()
		this.ETag = fmt.Sprint(time.Now().UnixNano())
		if bookTabs[this.Bookid] == nil {
			log.Fatal("addTab - book tabMap not found for bookid =", this.Bookid)
		}
		bookTabs[this.Bookid][this.Id] = this
		tabPrevids[this.Id] = make(previdMap)   // see previd.go
		dbWrite("saveNewTab", this, resultChan) // save new Tab to db

	case "saveNewTab":
		var err, err1, err2, err3 error
		var k, v []byte
		var bktBook, bktTabs, bktTab *bolt.Bucket

		tx, _ := db.Begin(true)
		bumpSequentialid(tx) // required for all db doc adds

		// --- add new tab to "tabs" bkt inside parent book bkt ---
		bktBook = openBook(tx, this.Bookid)
		bktTabs = openTabs(bktBook)
		k = bs(this.Id)
		v, err1 = json.Marshal(this)
		err2 = bktTabs.Put(k, v)

		// --- create bkt for this tab to hold its notes & Commit ---
		bktBook.CreateBucket(bs("tab_" + this.Id))
		err3 = tx.Commit()
		if err = lib.CheckErrs(errNotFatal, action, err1, err2, err3); err != nil {
			tx.Rollback()
			log.Fatal("Fatal DB Error-creating new tab bkts, Ending Program")
		}

		// --- create "notes","previds" bkts inside new tab's bkt & Commit ---
		tx, _ = db.Begin(true)
		bktBook = openBook(tx, this.Bookid)
		bktTab = openTab(bktBook, this.Id)
		_, err1 = bktTab.CreateBucket(bs("notes"))
		bktTab.CreateBucket(bs("previds"))
		err2 = tx.Commit()
		if err = lib.CheckErrs(errNotFatal, action, err1, err2); err != nil {
			tx.Rollback()
			log.Fatal("Fatal DB Error-creating notes bkt, Ending Program")
		}
		resultChan <- Result{DataOk, this.Id}

	case "changeTab":
		var tab *Tab
		if tab = bookTabs[this.Bookid][this.Id]; tab == nil {
			log.Fatal("changeTab, tab not found (bookid, tabid)", this.Bookid, this.Id)
		}
		tab.TabName = this.TabName
		tab.TabNumber = this.TabNumber
		tab.LastAccessTime = time.Now()
		tab.ETag = fmt.Sprint(time.Now().UnixNano())
		dbWrite("saveTabChange", this, resultChan) // save Tab to db

	case "saveTabChange":
		db.Update(func(tx *bolt.Tx) error {
			bktBook := openBook(tx, this.Bookid)
			bktTabs := openTabs(bktBook)
			k := bs(this.Id)
			v, err1 := json.Marshal(this)
			err2 := bktTabs.Put(k, v)
			return lib.CheckErrs(errFatal, action, err1, err2)
		})
		resultChan <- Result{DataOk, nil}

	default:
		log.Fatal("data Tab.process() invalid action ", action)
	}
}
