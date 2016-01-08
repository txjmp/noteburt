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

  
