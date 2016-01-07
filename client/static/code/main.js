var UrlPrefix;  // set in main() func, used for all ajax calls

var OpenMode;
const EditMode = 1;
const ViewMode = 2;
const CreateMode = 3;

// global css values
var Green1 = "#118C4E";
var Orange1 = "#FF9009";
var Tan1 = "#F3FAB6";
var Gray1 = "#BCBCA7";    //585858

var NoteWidth = "1200px"; // used in editview, noteview

var BaseFont = "600 1.1em muli, sans-serif";
var NoteFont = "600 1.1em muli, sans-serif";
var MonoFont = "normal 1em consolas, monospace";

var WinHeight = 0;
var WinWidth = 0;

var Views = {};

function main() {
	// If user opens book directly from url, web service returns index.html with [index.html var LoadBookName] set.
	// 	  the url will end with "book/bookname" which must be removed for data requests
	//	  all data requests are prefixed with value of UrlPrefix
	UrlPrefix = window.location.href;   // values: http://localhost/, noteburt.com, www.noteburt.com
	urlIndexOfBook = UrlPrefix.indexOf("book");
	if(urlIndexOfBook > -1) {
		UrlPrefix = UrlPrefix.substring(0, urlIndexOfBook);
	}
	console.log(UrlPrefix);

	WinHeight = $(window).height();
	WinWidth = $(window).width();
	$(window).resize(function() {
		WinHeight = $(window).height();
		WinWidth = $(window).width();
	});
	Views.screen1 = new ViewScreen1();
	Views.tabMgr = new ViewTabMgr();
	Views.bookTabs = new ViewBookTabs();
	Views.noteEdit = new ViewNoteEdit();
	Views.noteView = new ViewNoteView();
	Views.position = new ViewPosition();
	for( view in Views ) {
		Views[view].build();
		console.log(view + " view build run");
	}	
	$(".view").hide();
	Views.screen1.display();
}

function trace(msg) {
  //console.log(msg);
}

var lastBroadcast = "";
function ShowBroadcast(msg) {  // used by login and getBroadcast
	if( msg == lastBroadcast ) return;
	lastBroadcast = msg;
	if (msg == "") return;
	Notice(msg);
}
function broadcastTimer() {
	console.log("broadcasttimer");
	if(OpenMode != undefined) {  // wait after login before starting
		getBroadcast();
	}
	setTimeout(broadcastTimer, 1000 * 60);
}
function getBroadcast() {
    var path = "broadcast/get";
	var settings = {
		url: UrlPrefix + path,
		method: "GET",
		dataType: "text",  // returned from server
	}
	$.ajax(settings)
		.done(function(data) {
			ShowBroadcast(data)
		})
		.fail(function(xhr) {
			Warning("Get Broadcast Message Failed: " + xhr.responseText);
		})
}
