package main

import (
	"lib"
	"log"
	"math/rand"
	d "nb/data"
	"os"
	"strconv"
	"time"
)

var doneChan chan string

func main() {
	var err error

	lib.TraceStart("stdout")
	dbName := "datafiles/test1.db"
	if err = os.Remove(dbName); err != nil {
		log.Fatal(err)
	}
	d.DataStart(dbName)

	bk1 := d.Book{BookName: "book1", AccessCode: "ab1cd"}
	resultChan := d.Data("addBook", &bk1)
	result := <-resultChan
	bookid := result.Val.(string)
	log.Println("bookid=", bookid)

	tb1 := d.Tab{Bookid: bookid, TabName: "tab1", TabNumber: 10}
	resultChan = d.Data("addTab", &tb1)
	result = <-resultChan
	tabid := result.Val.(string)
	log.Println("tabid=", tabid)

	noteids := make([]string, 0, 10)
	for i := 0; i < 3; i++ {
		note := new(d.Note)
		note.Content = "note content " + strconv.Itoa(i)
		note.When = time.Now()
		note.Mono = true
		noteParms := d.NoteParms{Bookid: bookid, Tabid: tabid, Previd: d.Zeroid, Note: note}
		resultChan = d.Data("addNote", &noteParms)
		if result = <-resultChan; result.Status != d.DataOk {
			log.Fatal(result.Val.(string))
		}
		noteids = append(noteids, result.Val.(string))
	}

	rand.Seed(time.Now().UnixNano()) // rand used by common.RandCode

	doneChan = make(chan string)

	//lib.TraceLevel(9)

	go writeLoop(bookid, tabid, noteids)

	//time.Sleep(1 * time.Second)

	go readLoop(bookid, tabid)

	done1 := <-doneChan
	log.Println(done1)

	done2 := <-doneChan
	log.Println(done2)

	log.Println("main end")
}

func readLoop(bookid, tabid string) {
	for i := 0; i < 200; i++ {
		dataParms := new(d.TabNotes)
		dataParms.Bookid = bookid
		dataParms.Tabid = tabid
		resultChan := d.Data("getTabNotes", dataParms)
		result := <-resultChan
		if result.Status != d.DataOk {
			log.Println(result.Val.(string))
			break
		}
		/*
			for noteid, note := range dataParms.Notes {
				log.Println(noteid, note.Content, note.When, note.Mono, note.Html)
			}
			for noteid, previd := range dataParms.Previds {
				log.Println(noteid, previd)
			}
		*/
	}
	doneChan <- "readLoop done"
}
func writeLoop(bookid, tabid string, noteids []string) {
	for i := 0; i < 10; i++ {
		note := new(d.Note)
		note.Id = noteids[rand.Intn(len(noteids))]
		note.Content = "changed Note " + note.Id
		note.When = time.Now()
		note.Html = true
		note.Mono = false
		noteParms := d.NoteParms{Bookid: bookid, Tabid: tabid, Note: note}
		resultChan := d.Data("changeNote", &noteParms)
		if result := <-resultChan; result.Status != d.DataOk {
			log.Fatal(result.Val.(string))
		}
	}
	doneChan <- "writeLoop done"
}
