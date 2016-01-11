var UrlPrefix;  // set in main() func, used for all ajax calls
var DirectOpen;  // see comments below

var OpenMode;
const EditMode = 1;
const ViewMode = 2;
const CreateMode = 3;

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
	
	if(LoadBookName == "BOOK_NAME")
		DirectOpen = false; 
	else
		DirectOpen = true;  // prevents user from accessing screen1, opens book directly

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
  console.log(msg);
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
