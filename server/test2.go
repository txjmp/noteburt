package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"nb/data"
	"nb/web"
	"net/http"
	"os"
)

type loginRequest struct {
	AccessCode string
}

const urlPrefix = "http://localhost/"

var token, bookid string
var tabNotes map[string][]string //k=tabid, val=slice of noteids

func main() {
	if len(os.Args) != 2 {
		log.Println("cmd missing book name")
		os.Exit(1)
	}
	bookName := os.Args[1]

	tabNotes = make(map[string][]string)

	//	ADD BOOK
	requestData := loginRequest{AccessCode: "abc1def"}
	jsonBuf, err := dataToJsonBuf(requestData)
	if err != nil {
		log.Fatalln(err)
	}
	url := urlPrefix + "create/" + bookName
	resp, err := http.Post(url, "application/json", jsonBuf)
	loginResponse := new(web.TestLoginResponse)
	json.NewDecoder(resp.Body).Decode(loginResponse)
	resp.Body.Close()
	token = loginResponse.Token
	bookid = loginResponse.Bookid
	var tabid string
	for k, v := range loginResponse.Tabs {
		log.Println(v.TabName, v.TabNumber)
		tabid = k
		break
	}
	tabNotes[tabid] = make([]string, 0, 10)

	//	ADD NOTE TO TAB 1
	noteData := web.TestNoteRequest{
		Content: "note 1",
		Previd:  data.Zeroid,
	}
	noteid := addNote(tabid, noteData)
	tabNotes[tabid] = append(tabNotes[tabid], noteid)

	//	ADD TAB 10
	tabData := web.TestTabRequest{
		TabNumber: 10,
		TabName:   "Tab Ten",
	}
	jsonBuf, _ = dataToJsonBuf(tabData)

	url = urlPrefix + "tab/" + token + "/" + bookid
	resp, err = http.Post(url, "application/json", jsonBuf)
	addTabResponse := make([]byte, 12)
	resp.Body.Read(addTabResponse)
	resp.Body.Close()
	tabid = string(addTabResponse)
	tabNotes[tabid] = make([]string, 0, 10)
	log.Println("new tab id", tabid)

	// ADD NOTE TO TAB 10
	noteData = web.TestNoteRequest{
		Content: "tab 10 note 1",
		Previd:  data.Zeroid,
	}
	noteid = addNote(tabid, noteData)
	tabNotes[tabid] = append(tabNotes[tabid], noteid)

	// CHANGE NOTE
	noteData = web.TestNoteRequest{
		Content: "tab 10 note 1 changed",
	}
	changeNote(tabid, noteid, noteData)

	// ADD NOTE 2 TO TAB 10
	noteData = web.TestNoteRequest{
		Content: "tab 10 note 2",
		Previd:  data.Zeroid,
	}
	noteid = addNote(tabid, noteData)
	tabNotes[tabid] = append(tabNotes[tabid], noteid)

	// DELETE NOTE 2
	deleteNote(tabid, noteid)

	// ADD NOTE 3 TO TAB 10
	noteData = web.TestNoteRequest{
		Content: "tab 10 note 3",
		Previd:  data.Zeroid,
	}
	noteid = addNote(tabid, noteData)
	tabNotes[tabid] = append(tabNotes[tabid], noteid)

	// ADD NOTE 4 TO TAB 10
	noteData = web.TestNoteRequest{
		Content: "tab 10 note 4",
		Previd:  data.Zeroid,
	}
	noteid = addNote(tabid, noteData)
	tabNotes[tabid] = append(tabNotes[tabid], noteid)

	//	POSITION NOTE 4 AFTER NOTE 3
	noteData = web.TestNoteRequest{
		Previd: tabNotes[tabid][2],
	}
	positionNote(tabid, tabNotes[tabid][3], noteData)

	// GET TAB 10 NOTES
	url = urlPrefix + "tabnotes/" + token + "/" + bookid + "/" + tabid
	resp, err = http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("get tab 10 notes", resp.Status)
	var notes = make(map[string]data.NoteResponseRec)
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
	log.Println("order should be 3, 4, 1")
}

func addNote(tabid string, noteData web.TestNoteRequest) (noteid string) {
	jsonBuf, _ := dataToJsonBuf(noteData)
	url := urlPrefix + "note/" + token + "/" + bookid + "/" + tabid
	resp, _ := http.Post(url, "application/json", jsonBuf)
	//addNoteResponse := make([]byte, 12)
	//resp.Body.Read(addNoteResponse)
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	noteid = string(result)
	log.Println("new note added", noteid)
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
	//result := make([]byte, 100)
	//resp.Body.Read(result)
	result, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Println("change note result:", string(result))
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
	log.Println("position note result:", string(result))
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
	log.Println("delete result:", string(result))
}
func dataToJsonBuf(data interface{}) (*bytes.Buffer, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	jsonBuf := bytes.NewBuffer(jsonData)
	return jsonBuf, nil
}

/*
func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

image/jpeg
application/json
text/plain
*/
