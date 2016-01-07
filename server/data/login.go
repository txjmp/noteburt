package data

import (
	"log"
	"nb/common"
	"time"
)

var loginMap map[string]*Login

func init() {
	loginMap = make(map[string]*Login)
}

type Login struct {
	Token  string
	When   time.Time
	Bookid string
	Type   string // edit, view
}

func (this *Login) process(action string, resultChan chan Result) {
	switch action {

	case "addLogin":
		this.Token = common.RandCode(10)
		this.When = time.Now()
		loginMap[this.Token] = this
		resultChan <- Result{DataOk, nil}

	case "auth":
		if login := loginMap[this.Token]; login == nil {
			resultChan <- Result{DataNotFound, "Invalid Login - Token Not Found"}
		} else if login.Bookid != this.Bookid {
			resultChan <- Result{DataError, "Invalid Token for Book"}
		} else if login.Type == "view" && this.Type == "edit" {
			resultChan <- Result{DataError, "Edit Not Allowed"}
		} else {
			resultChan <- Result{DataOk, nil}
		}

	case "removeLogin":
		delete(loginMap, this.Token)
		resultChan <- Result{DataOk, nil}

	default:
		log.Fatal("invalid login action", action)
	}
}
