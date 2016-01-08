# NoteBurt Server Packages
1. main (top level directory, other pkgs are in a sub-directory with same name as pkg)
  * nb.go - contains main func, performs app startup actions
2. web
  * web.go - assigns http handlers, some init functions, starts ListenAndServe
  * ???handlers.go - each contains a related set of http request handlers
3. data
  * data.go
    * code to handle all data requests from handlers
    * startup func that preloads some data from the database, starts dataDispatch goroutine
    * some pkg level data stores, used by various pkg funcs
    * request, result, processor (interface) type definitions
    * dataDispatch func which processes all data requests from the request channel
  * db.go
    * database init, and starts dbWrite goroutine
    * dbWrite func which processes all db write requests (posted by dataDispatch)
    * design prevents subsequent data requests from waiting on db writes 
  * book.go, tab.go, note.go - data type def and processors for each data type
  * login.go - handles login, authentication
  * previd.go - hanldes note ordering via a linked list  
    maintains the id of the note preceding each note (previd) in sorted order  
    the 1st note's previd value is the zeroid (all zeros value)  
  * getbooktabs.go, gettabnotes.go - each handles a specific request
4. schedule - process scheduled actions in separate goroutine  
