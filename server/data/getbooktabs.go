package data

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type TabInfo struct { // tab data sent back to client
	TabNumber int
	TabName   string
	Hidden    bool
}
type GetBookTabs struct {
	Bookid     string
	TabInfoMap map[string]*TabInfo // key is tabid
}

func (this *GetBookTabs) process(action string, resultChan chan Result) {
	if bookTabs[this.Bookid] == nil { // tabs not yet loaded from db for this book
		bookTabs[this.Bookid] = make(tabMap)
		loadBookTabs(this.Bookid)
	}
	this.TabInfoMap = make(map[string]*TabInfo)
	for tabid, tab := range bookTabs[this.Bookid] {
		this.TabInfoMap[tabid] = &TabInfo{tab.TabNumber, tab.TabName, tab.Hidden}
	}
	resultChan <- Result{DataOk, nil}
}

// called on 1st request for tabs for bookid
//  load bookTabs[bookid] & tabPrevids[tabid] for each tab in book
//  bookTabs & tabPrevids are kept in sync with the database
func loadBookTabs(bookid string) {
	db.View(func(tx *bolt.Tx) error {
		bktBook := openBook(tx, bookid)
		bktTabs := openTabs(bktBook)
		cursor := bktTabs.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			tab := new(Tab)
			if err := json.Unmarshal(v, tab); err != nil {
				log.Fatal("loadBookTabs Unmarshal Tab Failed: ", err)
			}
			tab.LastAccessTime = time.Now()
			tab.ETag = fmt.Sprint(time.Now().UnixNano())
			bookTabs[bookid][tab.Id] = tab
			bktTab := openTab(bktBook, tab.Id)
			bktPrevids := openPrevids(bktTab)
			if bktPrevids == nil {
				log.Fatal("loadBookTabs previds bkt error (bookid, tabid)", bookid, tab.Id)
			}
			tabPrevids[tab.Id] = make(previdMap)
			previdsCursor := bktPrevids.Cursor()
			for noteid, previd := previdsCursor.First(); noteid != nil; noteid, previd = previdsCursor.Next() {
				tabPrevids[tab.Id][string(noteid)] = string(previd)
			}
		}
		return nil
	})
}
