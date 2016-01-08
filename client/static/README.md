# NoteBurt Client Modules
1. main.js
  * contains main function, starting point for app
  * global values used by multiple modules
  * Views obj that contains all views
2. hub.js
  * manages interactions betweens views
3. view_?????.js
  * user view modules
  * each contains 1 constructor function that defines view
  * all views have html code, methods: build, display, events
  * display method repopulates view data elements from values in data.js
4. data.js
  * contains data values used by all views
5. lib.js
  * constants & abbreviated codes for view html/css definitions
  * GenHtml function
  * GenCss function
  * Notice, Warning functions
  * Misc utility functions
  
