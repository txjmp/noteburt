package data

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"lib"
	"log"
	"time"
)

type Note struct {
	Id       string
	Content  string
	When     time.Time
	Mono     bool
	Html     bool
	Markdown bool
}

type NoteParms struct {
	Bookid       string
	Tabid        string
	Previd       string // used for addNote, positionNote; noteid that newid will follow
	Note         *Note
	OrderChanges map[string]string
}

func (this *NoteParms) process(action string, resultChan chan Result) {
	var tab *Tab
	if tab = bookTabs[this.Bookid][this.Tabid]; tab == nil {
		log.Fatal("NoteParms.process Tab Not Found (bookid, tabid)", this.Bookid, this.Tabid)
	}
	// all NoteParm actions are changes, so ETag must be updated
	tab.ETag = fmt.Sprint(time.Now().UnixNano())
	tab.LastAccessTime = time.Now()

	switch action {

	case "addNote":
		this.Note.Id = getNextid()
		previdMap := tabPrevids[this.Tabid]
		this.OrderChanges = previdMap.add(this.Note.Id, this.Previd)
		dbWrite("saveNewNote", this, resultChan)

	case "saveNewNote":
		var err1, err2, err3 error
		var k, v []byte
		db.Update(func(tx *bolt.Tx) error {
			bumpSequentialid(tx) // required for all database adds
			bktBook := openBook(tx, this.Bookid)
			_, bktNotes, bktPrevids := openTabAll(bktBook, this.Tabid)
			k = bs(this.Note.Id)
			v, err1 = json.Marshal(this.Note)
			err2 = bktNotes.Put(k, v)
			for noteid, previd := range this.OrderChanges {
				if err3 = bktPrevids.Put(bs(noteid), bs(previd)); err3 != nil {
					break
				}
			}
			return lib.CheckErrs(errFatal, action, err1, err2, err3)
		})
		resultChan <- Result{DataOk, nil}

	case "changeNote":
		dbWrite("saveNoteChange", this, resultChan)

	case "saveNoteChange":
		db.Update(func(tx *bolt.Tx) error {
			bktBook := openBook(tx, this.Bookid)
			bktTab := openTab(bktBook, this.Tabid)
			bktNotes := openNotes(bktTab)
			k := bs(this.Note.Id)
			v, err1 := json.Marshal(this.Note)
			err2 := bktNotes.Put(k, v)
			return lib.CheckErrs(errFatal, action, err1, err2)
		})
		resultChan <- Result{DataOk, nil}

	case "positionNote": // delete, add back in new position
		previdMap := tabPrevids[this.Tabid]
		this.OrderChanges = previdMap.delete(this.Note.Id)
		addChanges := previdMap.add(this.Note.Id, this.Previd)
		for noteid, previd := range addChanges { // merge addChanges, replacing delete with add entry
			this.OrderChanges[noteid] = previd
		}
		dbWrite("savePositionNote", this, resultChan)

	case "savePositionNote":
		var err1 error
		db.Update(func(tx *bolt.Tx) error {
			bktBook := openBook(tx, this.Bookid)
			bktTab := openTab(bktBook, this.Tabid)
			bktPrevids := openPrevids(bktTab)
			for noteid, previd := range this.OrderChanges {
				if previd == "delete" {
					err1 = bktPrevids.Delete(bs(noteid))
				} else {
					err1 = bktPrevids.Put(bs(noteid), bs(previd))
				}
				if err1 != nil {
					break
				}
			}
			return lib.CheckErrs(errFatal, action, err1)
		})
		resultChan <- Result{DataOk, nil}

	case "deleteNote":
		previdMap := tabPrevids[this.Tabid]
		this.OrderChanges = previdMap.delete(this.Note.Id)
		dbWrite("saveDeleteNote", this, resultChan)

	case "saveDeleteNote":
		var err1, err2 error
		db.Update(func(tx *bolt.Tx) error {
			bktBook := openBook(tx, this.Bookid)
			_, bktNotes, bktPrevids := openTabAll(bktBook, this.Tabid)
			err1 = bktNotes.Delete(bs(this.Note.Id))
			for noteid, previd := range this.OrderChanges {
				if previd == "delete" {
					err2 = bktPrevids.Delete(bs(noteid))
				} else {
					err2 = bktPrevids.Put(bs(noteid), bs(previd))
				}
				if err2 != nil {
					break
				}
			}
			return lib.CheckErrs(errFatal, action, err1, err2)
		})
		resultChan <- Result{DataOk, nil}

	default:
		log.Fatal("data NoteParms.process() invalid action ", action)
	}
}
