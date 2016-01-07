package data

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"lib"
	"lib/date"
	"log"
)

type NoteResponseRec struct {
	Previd   string
	Content  string
	When     [6]int //yr,mo,da,hr,min,sec
	Mono     bool
	Html     bool
	Markdown bool
}
type GetTabNotes struct {
	Bookid string
	Tabid  string
	ETag   string
	Notes  map[string]*NoteResponseRec
}

func (this *GetTabNotes) process(action string, resultChan chan Result) {
	tab := bookTabs[this.Bookid][this.Tabid]
	lib.Trace(1, "Tab ETag =", tab.ETag)
	if tab.ETag == this.ETag {
		resultChan <- Result{DataNotChanged, nil}
		return
	}
	this.ETag = tab.ETag // loaded into response header
	this.Notes = make(map[string]*NoteResponseRec)
	previds := tabPrevids[this.Tabid] // tabPrevids map kept updated, so db is not accessed
	db.View(func(tx *bolt.Tx) error {
		bktBook := openBook(tx, this.Bookid)
		bktTab := openTab(bktBook, this.Tabid)
		bktNotes := openNotes(bktTab)
		if bktNotes == nil {
			log.Fatal("GetTabNotes bkt problem (bookid, tabid)", this.Bookid, this.Tabid)
		}
		cursor := bktNotes.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			note := new(Note)
			if err := json.Unmarshal(v, note); err != nil {
				log.Fatal("GetTabNotes json Unmarshal failed", err)
			}
			dte := date.TimeToSimpleDate(note.When)
			this.Notes[note.Id] = &NoteResponseRec{
				Content:  note.Content,
				When:     date.SimpleDateToArray(dte),
				Mono:     note.Mono,
				Html:     note.Html,
				Markdown: note.Markdown,
				Previd:   previds[note.Id],
			}
		}
		return nil
	})
	resultChan <- Result{DataOk, nil}
}
