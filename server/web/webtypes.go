package web

import d "nb/data"

// Note: many responses are simple text (not special types)
// The request URL contains token, bookid, tabid, noteid, bookname

type loginRequest struct {
	AccessCode string
}
type loginResponse struct {
	Token     string
	Bookid    string
	BookName  string
	Broadcast string
	Tabs      map[string]*d.TabInfo
}
type TestLoginResponse struct {
	Token     string
	Bookid    string
	BookName  string
	Broadcast string
	Tabs      map[string]*d.TabInfo // TabInfo defined in data/getbooktabs.go
}
type tabRequest struct {
	TabNumber int
	TabName   string
	Hidden    bool
}
type TestTabRequest struct {
	TabNumber int
	TabName   string
	Hidden    bool
}
type noteRequest struct {
	Content  string
	Mono     bool
	Html     bool
	Markdown bool
	Previd   string
}
type TestNoteRequest struct {
	Content  string
	Mono     bool
	Html     bool
	Markdown bool
	Previd   string
}
