package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"nb/data"
	"nb/web"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type loginRequest struct {
	AccessCode string
}

const urlPrefix = "http://localhost/"

var wg sync.WaitGroup
var token, bookid string
var tabN int32 = 0

func main() {
	if len(os.Args) != 2 {
		log.Println("cmd missing book name")
		os.Exit(1)
	}
	bookName := os.Args[1]

	rand.Seed(time.Now().UnixNano()) // rand used by common.RandCode

	//	ADD BOOK
	requestData := loginRequest{AccessCode: "abc1def"}
	jsonBuf, _ := dataToJsonBuf(requestData)
	url := urlPrefix + "create/" + bookName
	resp, _ := http.Post(url, "application/json", jsonBuf)
	loginResponse := new(web.TestLoginResponse)
	json.NewDecoder(resp.Body).Decode(loginResponse)
	resp.Body.Close()
	token = loginResponse.Token
	bookid = loginResponse.Bookid

	wg.Add(2)
	go process()
	go process()
	wg.Wait()
}
func process() {
	defer wg.Done()

	var err error
	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		//	ADD TAB
		atomic.AddInt32(&tabN, 5)
		tabName := fmt.Sprintf("Tab %d", tabN)
		log.Println(tabName)
		tabData := web.TestTabRequest{TabNumber: int(tabN), TabName: tabName}
		jsonBuf, _ := dataToJsonBuf(tabData)
		url := urlPrefix + "tab/" + token + "/" + bookid
		resp, _ := http.Post(url, "application/json", jsonBuf)
		result, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		tabid := string(result)

		wait()

		// ADD NOTES
		noteData := web.TestNoteRequest{Content: "note 1", Previd: data.Zeroid}
		noteid1 := addNote(tabid, noteData)
		wait()
		noteData = web.TestNoteRequest{Content: "note 2", Previd: data.Zeroid}
		addNote(tabid, noteData)
		wait()
		noteData = web.TestNoteRequest{Content: "note 3", Previd: data.Zeroid}
		noteid3 := addNote(tabid, noteData)
		wait()

		// GET TAB NOTES
		url = urlPrefix + "tabnotes/" + token + "/" + bookid + "/" + tabid
		if resp, err = http.Get(url); err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		wait()

		// CHANGE NOTE
		noteData = web.TestNoteRequest{Content: "note 1 changed"}
		changeNote(tabid, noteid1, noteData)
		wait()

		//	POSITION NOTE
		noteData = web.TestNoteRequest{Previd: noteid1}
		positionNote(tabid, noteid3, noteData)
		wait()

		// GET TAB NOTES
		url = urlPrefix + "tabnotes/" + token + "/" + bookid + "/" + tabid
		if resp, err = http.Get(url); err != nil {
			log.Fatal(err)
		}
		notes := make(map[string]data.NoteResponseRec)
		json.NewDecoder(resp.Body).Decode(&notes)
		resp.Body.Close()
		noteOrder := make(map[string]string)
		for k, v := range notes {
			noteOrder[v.Previd] = k
		}
		nextid := data.Zeroid
		for i := 0; i < len(noteOrder); i++ {
			id := noteOrder[nextid]
			log.Println(notes[id].Content)
			nextid = id
		}
	}
}

func addNote(tabid string, noteData web.TestNoteRequest) (noteid string) {
	jsonBuf, _ := dataToJsonBuf(noteData)
	url := urlPrefix + "note/" + token + "/" + bookid + "/" + tabid
	resp, err := http.Post(url, "application/json", jsonBuf)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	noteid = string(result)
	return
}
func changeNote(tabid, noteid string, noteData web.TestNoteRequest) {
	jsonBuf, _ := dataToJsonBuf(noteData)
	url := urlPrefix + "note/" + token + "/" + bookid + "/" + tabid + "/" + noteid
	client := &http.Client{}
	request, _ := http.NewRequest("PUT", url, jsonBuf)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	//result, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
}
func positionNote(tabid, noteid string, noteData web.TestNoteRequest) {
	jsonBuf, _ := dataToJsonBuf(noteData)
	url := urlPrefix + "positionnote/" + token + "/" + bookid + "/" + tabid + "/" + noteid
	client := &http.Client{}
	request, _ := http.NewRequest("PUT", url, jsonBuf)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, 100)
	resp.Body.Read(result)
	resp.Body.Close()
}
func deleteNote(tabid, noteid string) {
	url := urlPrefix + "note/" + token + "/" + bookid + "/" + tabid + "/" + noteid
	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", url, nil)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]byte, 100)
	resp.Body.Read(result)
	resp.Body.Close()
}
func dataToJsonBuf(data interface{}) (*bytes.Buffer, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	jsonBuf := bytes.NewBuffer(jsonData)
	return jsonBuf, nil
}
func wait() {
	x := rand.Intn(50) + 20
	//log.Println("wait", x)
	time.Sleep(time.Duration(x) * time.Millisecond)
}
