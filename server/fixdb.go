package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type Book struct {
	Id            string
	BookName      string
	PlainBookName string // used for matching
	AccessCode    string
	Email         string
}
type Tab struct {
	Id             string
	Bookid         string
	TabNumber      int
	TabName        string
	LastAccessTime time.Time // updated on any access (read,write) to any note in tab
	ETag           string    // updated on any change to a note
}
type Note struct {
	Id       string
	Content  string
	When     time.Time
	Mono     bool
	Html     bool
	Markdown bool
}

//type previdMap map[string]string    // key=noteid, val=previd (in sorted order)
//var tabPrevids map[string]previdMap // key=tabid

type bs []byte // byte slice type used for value conversions

var db *bolt.DB

func main() {
	fix := flag.Bool("fix", false, "include -fix to run fix() func")
	prod := flag.Bool("prod", false, "include -prod to open prod db")
	flag.Parse()

	var err error
	var dbName string
	if *prod {
		dbName = "/webapps/nb/db/prod.db"
	} else {
		dbName = "/webapps/nb/db/test.db"
	}

	if db, err = bolt.Open(dbName, 0600, nil); err != nil {
		log.Fatal("db open fail", err)
	}
	defer db.Close()

	if *fix {
		fixDB()
		log.Println("fix run")
	} else {
		log.Println("fix NOT run")
	}

	db.View(func(tx *bolt.Tx) error {
		bktBooks := openBooks(tx)
		booksCursor := bktBooks.Cursor()
		// --- Books ---
		for k, v := booksCursor.First(); k != nil; k, v = booksCursor.Next() {
			//fmt.Println("book k, v", string(k), string(v))
			book := new(Book)
			if err := json.Unmarshal(v, book); err != nil {
				log.Fatal("book json Unmarshal failed", err)
			}
			fmt.Println("Book", book.Id, book.BookName, book.AccessCode)
			bktBook := openBook(tx, book.Id)
			bktTabs := openTabs(bktBook)
			tabsCursor := bktTabs.Cursor()
			// --- Tabs In Book ---
			for k, v := tabsCursor.First(); k != nil; k, v = tabsCursor.Next() {
				tab := new(Tab)
				if err := json.Unmarshal(v, tab); err != nil {
					log.Fatal("tab json Unmarshal failed", err)
				}
				fmt.Println("tab info", tab.Id, tab.TabNumber, tab.TabName)
				bktNotes, bktPrevids := openTabAll(bktBook, tab.Id)
				// --- Notes in Tab ---
				notesCursor := bktNotes.Cursor()
				for k, v := notesCursor.First(); k != nil; k, v = notesCursor.Next() {
					note := new(Note)
					if err := json.Unmarshal(v, note); err != nil {
						log.Fatal("Note json Unmarshal failed", err)
					}
					previd := bktPrevids.Get(k)
					if previd == nil {
						fmt.Println("no previd found")
					}
					x := len(note.Content) - 1
					if x > 10 {
						x = 10
					}
					fmt.Println("noteid, content, previd", note.Id, note.Content[0:x], string(previd))
				}
				// --- Previds in Tab ---
				previdsCursor := bktPrevids.Cursor()
				for k, v := previdsCursor.First(); k != nil; k, v = previdsCursor.Next() {
					noteid := string(k)
					previd := string(v)
					fmt.Println("tabid, noteid, previd=", tab.Id, noteid, previd)
				}
			}
		}
		return nil
	})
}

func fixDB() {
	db.Update(func(tx *bolt.Tx) error {
		bookid := 2
		tabid := 9
		noteid := 13
		previd := 10
		bktBook := openBook(tx, recid(bookid))
		bktTab := openTab(bktBook, recid(tabid))
		bktPrevids := openPrevids(bktTab)
		k := recid(noteid)
		v := recid(previd)
		bktPrevids.Put(bs(k), bs(v))
		return nil
	})
}

func recid(id int) string {
	return fmt.Sprintf("%012d", id)
}

// ================================================

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
func openTabAll(bktBook *bolt.Bucket, tabid string) (*bolt.Bucket, *bolt.Bucket) {
	bktTab := bktBook.Bucket(bs("tab_" + tabid))
	bktNotes := bktTab.Bucket(bs("notes"))
	bktPrevids := bktTab.Bucket(bs("previds"))
	return bktNotes, bktPrevids
}
