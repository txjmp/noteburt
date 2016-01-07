function ViewNoteView() {
	this.id = "#view_noteview";
	this.html = [
		zdiv, {id:"top_div", w:"100%", flt:"left", bkg:Green1, mbottom:"0px",
			pad:"10px", b:"1px solid black", end:">"},
		zlbl, {id:"noteburt", text:"NoteBurt", flt:"left", 
			f:"1.5em AlphaSlabOne", c:"white", textshadow:"1px 1px 2px black",
			mleft:"10px", end:"tag"}, 
		zlbl, {id:"tab_name", text:"No Tab Selected", flt:"left",
			f:"600 1.5em dosis, sans-serif", c:"navy",
			mleft:"30px", pright:"15px", end:"tag"},
		zlbl, {id:"note_title", flt:"left",
			f:"600 1.5em dosis, sans-serif", c:"navy",
			bleft:"2px solid black", pleft:"15px", end:"tag"},
		zend, zdiv,
	  
		// left side of screen containing buttons and options
		zdiv, {id:"options_div", w:"90px", h:"100%", flt:"left", bkg:"gray"},
		zbtn, {id:"close_btn", text:"Close", clas:"option_btn", end:"tag"},
		zbtn, {id:"next_btn", text:"Next", clas:"option_btn", end:"tag"},
		zbtn, {id:"prev_btn", text:"Prev", clas:"option_btn", end:"tag"},
		zend, zdiv,
	  
		// right side of screen containing note content
		zdiv, {id:"note_div", w:"80%", maxwidth:NoteWidth, h:"95%",	flt:"left", overflow:"auto", pad:"10px",},
		ztextarea, {id:"notetext", w:"95%", h:"95%", pad:"10px", readonly:"readonly", bleft:"1px solid black",
			bkg:"transparent", f:NoteFont, end:"tag"},
		zend, zdiv,

		// right side of screen containing html content
		zdiv, {id:"html_div", w:"80%", maxwidth:NoteWidth, h:"95%", disp:"none", position:"relative", f:NoteFont,
			flt:"left", overflow:"auto", bleft:"1px solid black", pad:"7px", m:"0px", end:"tag"},
	];
	this.css = [ 
		{ selector:".option_btn",
			settings: {f:BaseFont, c:"black", bkg:Green1,
			b:"1px solid black",
			bradius:"4px", w:"70px", mleft:"10px", mtop:"15px"}
		},
	];
	this.build = function() {
		var viewName = this.id.substring(1);
		$('body').append('<div id="' + viewName + '" class="view"></div>');

		$(this.id).css("height", WinHeight-50);
		$(this.id).css("width", WinWidth-50);

		var html = GenHtml(this.html);
		$(this.id).html(html);

		var css = GenCss(this.css);
		ApplyCss(this.id, css);

		this.events();
	}
	this.display = function(refresh) {
		$("body").css("background", Gray1);

		if(refresh) {
			$(this.id + " #notebook_name").text(DataBookName);
			$(this.id + " #tab_name").text(DataTabs[DataTabid].tabName);
			$(this.id + " #note_title").text(DataNotes[DataNoteid].title);
			var note = DataNotes[DataNoteid];
			if(note.html) {
				var html = note.content;
				$(this.id + " #note_div").hide();      
				$(this.id + " #html_div").html(html);
				$(this.id + " #html_div").show();
			} else if(note.markdown) {
				var html = marked(note.content);
				$(this.id + " #note_div").hide();      
				$(this.id + " #html_div").html(html);
				$(this.id + " #html_div").show();
			} else {
				$(this.id + " #notetext").val(note.content);
				$(this.id + " #html_div").hide();
				$(this.id + " #note_div").show();
			}
		}
		$(this.id).show();
	}
	this.events = function() {
		var viewid = this.id;
		var thisView = this;
		$(this.id + " #close_btn").click(function() {
			$(viewid).hide();
			Hub.viewNoteViewClose();
		});
		$(this.id + " #next_btn").click(function() {
			var nextid = DataNoteOrder[DataNoteid];
			if( nextid == undefined ) {
				Notice("No More Notes");
			} else {
				DataNoteid = nextid;
				thisView.display(true);
			}
		});
		$(this.id + " #prev_btn").click(function() {
			var previd = DataNotes[DataNoteid].previd;
			if( previd == DataZeroid ) {
				Notice("This is The First Note");
			} else {
				DataNoteid = previd;
				thisView.display(true);
			}
		});
	}
}
