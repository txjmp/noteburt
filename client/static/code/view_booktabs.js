var noteListFont = "600 1.2em muli, sans-serif";
function ViewBookTabs() {
	this.id = "#view_booktabs";
	this.html = [
		zdiv, {id:"top_div", w:"100%", flt:"left", bkg:Green1, mbottom:"0px",
			pad:"10px", b:"1px solid black", end:">"},
		zlbl, {id:"noteburt", text:"NoteBurt", flt:"left", 
			f:"1.6em AlphaSlabOne",	c:"white", textshadow:"1px 1px 2px black",
			mleft:"45px", end:"tag"},
		zlbl, {id:"notebook_name", bright:"3px solid black", flt:"left",
			f:"600 1.5em dosis, sans-serif", c:"navy",
			pright:"20px", mleft:"80px", end:"tag"},
		zlbl, {id:"tab_name", text:"No Tab Selected", flt:"left",
			f:"600 1.5em dosis, sans-serif", c:"navy",
			pleft:"20px", end:"tag"},
		zend, zdiv,
	// ------------------------------------------------------------------------  
	// left side of screen containing tabmgr_btn and tab_btn list
	// ------------------------------------------------------------------------
		zdiv, {id:"tab_div", w:"240px", h:"100%", mtop:"0px", flt:"left"},
		zbtn, {id:"addnote_btn", text:"Add Note", w:"100%", flt:"left",
			f:"600 1.3em muli, sans-serif", c:"black", textalign:"center",
			b:"1px solid black", bkg:"white",
			p:"10px", mbottom:"0px", mtop:"0px", end:"tag"},
		zbtn, {id:"tabmgr_btn", text:"Tabs ...", w:"100%", flt:"left",
			f:"600 1.3em muli, sans-serif", c:"black", textalign:"center",
			b:"1px solid black", bkg:"white",
			p:"10px", mbottom:"0px", mtop:"0px", end:"tag"},
		// tablist_div containing 1 btn for each tab
		zdiv, {id:"tablist_div", h:"90%", flt:"left", overflow:"auto", end:"tag"},
		//  ... tab list items inserted here by display method
	
		zend, zdiv,  // end tab_div

		// right side of screen containing note list for selected tab
		zdiv, {id:"notelist_div", w:"400px", h:"95%", flt:"left", m:"15px", overflow:"auto", end:"tag"},
		//  ... note list items inserted here by display method
	]

	this.build = function() {
		var viewName = this.id.substring(1);
		$('body').append('<div id="' + viewName + '" class="view"></div>');

		$(this.id).css("height", WinHeight-50);
		$(this.id).css("width", WinWidth-50);

		var html = GenHtml(this.html);
		$(this.id).html(html);

		$(this.id + " #noteburt").css("cursor", "pointer");  // logo is btn

		this.events();

		// css for class .tab_btn, to be applied when elements are reloaded
		var shortCss = [
			{ selector: ".tab_btn", 
			  settings:{ f:BaseFont, c:"black", w:"240px",
			  pad:"10px", b:"none", btop:"1px solid black",
			  mbottom:"0px", textalign:"left"}
			}
		]
		this.cssTabBtn = GenCss(shortCss);

		// css for class .notelist_item, to be applied when elements are reloaded
		shortCss = [
			{ selector: ".notelist_item", 
			  settings:{ f:noteListFont, mtop:"7px", w:"95%",
			  c:'black',
			  textalign:"left", flt:"left", bkg:"transparent", b:"none"}
			},
		]
		this.cssNoteListItem = GenCss(shortCss);
	}
	this.display = function(refresh) {
		$("body").css("background", Gray1);

		if(refresh) {
			$(this.id + " #notebook_name").text(DataBookName);
			this.loadTabList();
			$(this.id + " #notelist_div").empty();
			if( DataTabid ) {
				this.loadNoteList();
			}
		}
		$(this.id).show();

		if( OpenMode == ViewMode ) {
			$(this.id + " #addnote_btn, #tabmgr_btn").hide();
			$(this.id + " #tablist_div").css("margin-top", "10px");
		}
	}
	this.events = function() {
		var viewid = this.id;
		$(viewid + " #addnote_btn").click(function() {
			if(!DataTabid) {
				Notice("A Tab Must Be Selected First, Thanks");
				return;
			}
			Hub.newNote();
		});
		$(viewid + " #tabmgr_btn").click(function() {
			Views.tabMgr.display(true);
		});
		$(viewid + " #noteburt").click(function() {
			Hub.viewBookTabs_NoteburtClicked();
		});
	}
//	this.tabColors = ["rgba(255,0,0,.4", "rgba(0,255,0,.4)", "rgba(0,0,255,.4)", "rgba(247,239,12,.7)", "rgba(170,10,250,.4)", "rgba(250,100,10,.7)"];
	this.tabColors = ["rgba(255,0,0,.7", "rgba(0,255,0,.7)", "rgba(0,0,255,.7)", "rgba(247,239,12,.7)"];
	this.loadTabList = function() {
		var viewid = this.id;
		var thisView = this;
		var html = "";
		var tabText, tab;
		var tabOrder = Object.keys(DataTabs).sort(function(a,b) {
			return DataTabs[a].tabNumber - DataTabs[b].tabNumber;
		});
		tabOrder.forEach(function(id) {
			tab = DataTabs[id];
			if( tab.tabNumber < 10 )
				tabText = "0" + tab.tabNumber + " " + tab.tabName;
			else
				tabText = tab.tabNumber + " " + tab.tabName;
			html += '<button class="tab_btn" id="' + id + '">' + tabText + '</button>\n';
		});
		$(viewid + " #tablist_div").html(html);

		ApplyCss(viewid, this.cssTabBtn);

		$(viewid + " .tab_btn").click(function(e) {
			var target = $(e.target);
			DataTabid = target.attr("id");
			$(viewid + " #tab_name").text(DataTabs[DataTabid].tabName);
			thisView.getTabNotes();
		});
		var tabColors = this.tabColors;
		var color = 0;
		$(viewid + " .tab_btn").each( function(index) {
			$(this).css("background", tabColors[color]);
			color++;
			if( color == tabColors.length )
				color = 0;
		})		
	}
	this.loadNoteList = function() {
		var html = "";
		var noteCount = Object.keys(DataNotes).length;
		var id;
		var nextid = DataZeroid;
		
		//if( byDate ) {
		//	var tabOrder = Object.keys(DataTabs).sort(function(a,b) {
		//		return DataTabs[a].tabNumber - DataTabs[b].tabNumber;
		//	});

		//	for(var i=0; i < noteCount; i++) {
		//
		//}
		
		for(var i=0; i < noteCount; i++) {
			id = DataNoteOrder[nextid];
			nextid = id;
			html += '<button class="notelist_item" value="' + id + '">' + DataNotes[id].title + '</button>\n';
		}
		$(this.id + " #notelist_div").html(html);

		ApplyCss(this.id, this.cssNoteListItem);

		$(this.id + " .notelist_item").click(function(e) {
			var noteid = $(e.target).val();
			Hub.showNote(noteid);
		});
	}
	this.loadDataNotes = function(response) {
		DataNotes = {};
		//if(Object.keys(response).length == 0) {
		//	Notice("No Notes Found For Tab. Click Add Note.");
		//	return;
		//}
		for(id in response) {
			DataNotes[id] = {
				"previd" : 	response[id].Previd,
				"content": 	response[id].Content,
				"when": 	ArrayToDate(response[id].When),   // convert [yr, mo, da, hr, min, sec]
				"mono": 	response[id].Mono,
				"html": 	response[id].Html,
				"markdown": response[id].Markdown,
			}
			DataNotes[id].title = DataNoteTitle(id);
		}
		DataBuildNoteOrder();
	}
	this.getTabNotes = function() {
		var viewid = this.id;
		var thisView = this;
		var path = "tabnotes/" + DataToken + "/" + DataBookid + "/" + DataTabid;
		var settings = {
			url: UrlPrefix + path,
			type: "GET",
			dataType: "json",  // returned from server
		}
		$(viewid + " button").prop("disabled",true);
		var sent = new Date();
		$.ajax(settings)
			.done(function(response) {
				var received = new Date();
				var elapsed = received - sent;
				trace("getNotes Response Received in " + elapsed + " milliseconds");
				thisView.loadDataNotes(response);
				thisView.loadNoteList();
			})
			.fail(function(xhr) {
				Err("Get Notes Failed \n" + xhr.responseText);
			})
			.always(function() {
				$(viewid + " button").prop("disabled",false);
			})
	}
}

/* -- getTabNotes response ---------------------------
type NoteResponseRec struct {
	Previd   string
	Content  string
	When     [6]int //yr,mo,da,hr,min,sec
	Mono     bool
	Html     bool
	Markdown bool
}
response  map[noteid]*NoteResponseRec
*/
