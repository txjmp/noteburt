package data

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"lib"
	"log"
)

type Book struct {
	Id            string
	BookName      string
	PlainBookName string // used for matching
	AccessCode    string
	Email         string
}

func (this *Book) process(action string, resultChan chan Result) {
	switch action {

	// getBook options:
	//		if Id is loaded, match on Id
	//		if AccessCode is loaded, match on AccessCode & BookName
	//		if only BookName is loaded, match on BookName

	case "getBook": // all book docs loaded into bookMap at startup
		var book *Book
		if this.Id != "" { // match on book id
			book = bookMap[this.Id]
		} else if this.AccessCode != "" { // match on AccessCode and BookName
			for _, v := range bookMap {
				if v.AccessCode == this.AccessCode {
					if v.PlainBookName == lib.PlainString(this.BookName) {
						book = v
						break
					}
				}
			}
		} else { // match on BookName
			matchBookName := lib.PlainString(this.BookName)
			for _, v := range bookMap {
				if v.PlainBookName == matchBookName {
					book = v
					break
				}
			}
		}
		if book != nil {
			this.Id = book.Id
			this.BookName = book.BookName
			this.Email = book.Email
			this.AccessCode = book.AccessCode
			this.PlainBookName = book.PlainBookName
			resultChan <- Result{DataOk, nil}
		} else {
			resultChan <- Result{DataNotFound, "Book Not Found"}
		}

	case "addBook":
		// check for dupe book name
		newName := lib.PlainString(this.BookName)
		for _, v := range bookMap {
			if newName == v.PlainBookName {
				resultChan <- Result{DataDuplicate, "Duplicate - Notebook Name Already Exists"}
				lib.Trace(0, "dupe book name", newName)
				return
			}
		}
		this.Id = getNextid()
		this.PlainBookName = newName
		bookMap[this.Id] = this
		bookTabs[this.Id] = make(tabMap)
		dbWrite("saveNewBook", this, resultChan) // save new Book to db

	case "saveNewBook":
		var err, err1, err2, err3 error
		var k, v []byte

		tx, _ := db.Begin(true)
		bumpSequentialid(tx) // required for all db adds

		// --- save book to "books" bkt -----------------------
		bktBooks := openBooks(tx)
		k = bs(this.Id)
		v, err1 = json.Marshal(this)
		err2 = bktBooks.Put(k, v)

		// --- create bkt to hold this book's data & Commit ---
		tx.CreateBucket(bs("book_" + this.Id))
		err3 = tx.Commit()
		if err = lib.CheckErrs(errNotFatal, action, err1, err2, err3); err != nil {
			tx.Rollback()
			log.Fatal("Fatal DB Error-creating book bkts, Ending Program")
		}

		// --- create "tabs" bkt inside this book's bkt & Commit ---
		tx, _ = db.Begin(true)
		bktBook := openBook(tx, this.Id)
		_, err1 = bktBook.CreateBucket(bs("tabs"))
		err2 = tx.Commit()
		if err = lib.CheckErrs(errNotFatal, action, err1, err2); err != nil {
			tx.Rollback()
			log.Fatal("Fatal DB Error-creating tabs bkt, Ending Program")
		}
		resultChan <- Result{DataOk, this.Id}

	case "saveBookChange": // not currently used
		db.Update(func(tx *bolt.Tx) error {
			bktBooks := openBooks(tx)
			k := bs(this.Id)
			v, _ := json.Marshal(this)
			err := bktBooks.Put(k, v)
			if err != nil {
				log.Fatal("book not changed", err)
			}
			return nil
		})

	default:
		log.Fatal("Book.process action invalid", action)
	}
}
