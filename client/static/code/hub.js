var REFRESH = true;

Hub = {};

Hub.openBook = function() {
	$(Views.screen1.id).hide();
	Views.bookTabs.display(REFRESH);
}
Hub.showNote = function(noteid) {
	DataNoteid = noteid;
	$(Views.bookTabs.id).hide();
	if(OpenMode == ViewMode)
		Views.noteView.display(REFRESH);
	else {
		window.onbeforeunload = confirmPageExit; 
		Views.noteEdit.display(REFRESH);
	}
}
Hub.newNote = function() {
	DataNoteid = undefined;
	$(Views.bookTabs.id).hide();
	if(OpenMode == EditMode) {
		window.onbeforeunload = confirmPageExit; 
		Views.noteEdit.display(REFRESH);
	} else
		alert("Hub.newNote program error");
}
Hub.noteDeleted = function() {
	DataNoteid = undefined;
	$(Views.noteEdit.id).hide();
	Views.bookTabs.display(REFRESH);
}
Hub.viewTabMgrClose = function() {
	Views.bookTabs.display(REFRESH);
}
Hub.viewNoteViewClose = function() {
	Views.bookTabs.display(false);
}
Hub.viewNoteEditClose = function() {
	window.onbeforeunload = null;
	Views.bookTabs.display(REFRESH);
}
Hub.viewBookTabs_NoteburtClicked = function() {
	if( DirectOpen ) return;  // see main.js for explanation
	DataTabid = undefined;
	DataNoteid = undefined;
	$(Views.bookTabs.id).hide();
	Views.screen1.display();	
}
function confirmPageExit()  {
	return "You have attempted to leave this page. " +
	"If you have any unsaved changes they will be lost.";
}			

