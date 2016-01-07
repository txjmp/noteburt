
const DataZeroid = "000000000000";

var DataToken;		// auth value included in data request url
var DataBookid;		// id of currently open book
var DataBookName;	// name of current book
var DataTabs;		// tab records for current book
var DataNotes;      // all note records for open tab
var DataTabid;		// id of open tab
var DataNoteid;     // id of open note
var DataNoteOrder;  // key=previd, val=noteid (linked list for ordering notes) 
var DataLastNoteid; // id of last note in sorted order, used when adding note (set in DataBuildNoteOrder below)

// Called on loadDataNotes (view_booktabs.js), note add/delete(view_noteedit.js), positionNote(view_position.js)
function DataBuildNoteOrder() {
	DataNoteOrder = {};  
	for( var id in DataNotes ) {
		DataNoteOrder[DataNotes[id].previd] = id;
	}
	DataLastNoteid = DataZeroid;  // if no notes for tab, this should work
	var nextid = DataZeroid;
	var cnt = Object.keys(DataNoteOrder).length;
	for( i = 0; i < cnt; i++ ) {
		DataLastNoteid = DataNoteOrder[nextid];
		nextid = DataLastNoteid;
	}
}
/*
var DataTabs = {
  "000000000001": {tabNumber:10, tabName:"Bubbles" },
  "000000000002": {tabNumber:20, tabName:"Trees" },
}
var DataNotes = {
  "000000000004": {previd: "00000000000", content:"1 Bubbles", when:new Date(), mono:false, html:false, markdown:false},
  "000000000005": {previd: "00000000006", content:"2 there you go", when:new Date(), mono:false, html:false, markdown:false},
  "000000000006": {previd: "00000000004", content:"3 please dont pop", when:new Date(), mono:false, html:false, markdown:false},
}
*/
// =============================================================
// set DataNotes.previd value on add, delete
// when changing position of note, delete then add are executed
// =============================================================
function DataAddPrevid(newid, newPrevid) {
	// newid is noteid of note being added
	// newPrevid is id of note that newid will follow in sorted order
	// set previd of note currently following newPrevid to newid
	for( id in DataNotes ) {
		if( DataNotes[id].previd == newPrevid ) { 
			DataNotes[id].previd = newid;
			break;
		}
	}
	DataNotes[newid].previd = newPrevid
}

// entry currently following deleteid will now follow entry before deleteid
// nothing is actually deleted
function DataDeletePrevid(deleteid) {
	for( id in DataNotes ) {
		if( DataNotes[id].previd == deleteid ) { 
			DataNotes[id].previd = DataNotes[deleteid].previd;
			break;
		}
	}
}
// construct note title from 1st line of note (an html comment for markup & html)
function DataNoteTitle(noteid) {
	var title;
	var content = DataNotes[noteid].content;
	var lineFeed = content.indexOf("\n");
	if( lineFeed == -1 ) {
		title = content;
	} else {
		title = content.substring(0, lineFeed);
	}
	if( title.substr(0,4) == "<!--" ) {
		var commentEnd = title.indexOf("-->");
		if( commentEnd < 0 )
			title = title.substring(4);
		else
			title = title.substring(4, commentEnd);
	}
	return title;
}

