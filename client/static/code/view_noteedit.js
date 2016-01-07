var noteChanged = false;  // boolean set to true when key pressed in note text, false on save

function ViewNoteEdit() {
	this.id = "#view_noteedit";
	this.html = [
		zdiv, {id:"top_div", w:"100%", flt:"left", bkg:Green1, mbottom:"0px",
			pad:"10px", b:"1px solid black", end:">"},
		zlbl, {id:"noteburt", text:"NoteBurt", flt:"left",
			f:"1.5em AlphaSlabOne", c:"white", textshadow:"1px 1px 2px black",
			mleft:"10px", end:"tag"},
		zlbl, {id:"tab_name", text:"No Tab Selected", flt:"left",
			f:"600 1.5em dosis sans-serif", c:"navy",
			mleft:"30px", pright:"15px", end:"tag"},
		zlbl, {id:"note_title", flt:"left",
			f:"600 1.5em dosis sans-serif", c:"navy",
			bleft:"2px solid black", pleft:"15px", end:"tag"},
		zend, zdiv,
	  
		// left side of screen containing buttons and options
		zdiv, {id:"options_div", w:"135px", h:"100%", flt:"left", bkg:"gray"},
		zbtn, {id:"close_btn", text:"Close", clas:"edit_btn", end:"tag"},
		zbtn, {id:"save_btn", text:"Save", clas:"edit_btn", end:"tag"},
		zbtn, {id:"delete_btn", text:"Delete", clas:"edit_btn", end:"tag"},
		zbtn, {id:"position_btn", text:"Position", clas:"edit_btn", end:"tag"},
		zditto, '<label><input type="checkbox" id="mono_chk"/>Monospace</label><br>',
		zditto, '<label><input type="checkbox" id="markdown_chk"/>Markdown</label><br>',
		zditto, '<label><input type="checkbox" id="html_chk"/>HTML</label><br>',
		zditto, '<label><input type="checkbox" id="view_chk"/>View</label><br>',
		zend, zdiv,

		// right side of screen containing note content
		zdiv, {id:"note_div", w:"80%", h:"95%", flt:"left", overflow:"auto", pad:"10px"},
		ztextarea, {id:"notetext", w:"95%", h:"95%", pad:"10px", tab:"4", bkg:"transparent", end:"tag"},
		zend, zdiv,
	 
		// right side of screen containing html content
		zdiv, {id:"html_div", w:"80%", h:"95%", disp:"none", flt:"left", position:"relative",
			f:NoteFont,	maxwidth:NoteWidth, overflow:"auto",
			m:"0px", p:"7px", end:"tag"},
	]
	this.css = [
		{ selector:".edit_btn",
		  settings: {f:BaseFont, c:"black", bkg:Green1,
			b:"1px solid black", bradius:"4px", w:"100px", mleft:"10px", mtop:"15px"}
		},
		{ selector:"#options_div label",
		  settings: {f:BaseFont, mtop:"20px", mleft:"5px", c:"black"}
		},
		{ selector:"#options_div input[type='checkbox']",
		  settings: {mtop:"15px", mleft:"5px"}
		},
	]
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
			if( DataNoteid ) {
				var note = DataNotes[DataNoteid];
				$(this.id + " #note_title").text(note.title);
				$(this.id + " #notetext").val(note.content);
				$(this.id + " #mono_chk").prop("checked", note.mono);
				$(this.id + " #html_chk").prop("checked", note.html);
				$(this.id + " #markdown_chk").prop("checked", note.markdown);
			} else {  // new note
				$(this.id + " #note_title").text("");
				$(this.id + " #notetext").val("");
				$(this.id + " #mono_chk").prop("checked", false);
				$(this.id + " #html_chk").prop("checked", false);
				$(this.id + " #markdown_chk").prop("checked", false);
			}
			if( $(this.id + " #mono_chk").prop("checked") ) {
				$(this.id + " #notetext").css("font", MonoFont);
			} else {
				$(this.id + " #notetext").css("font", NoteFont);      
			}
			$(this.id + " #view_chk").prop("checked",false);
			$(this.id + " #html_div").hide();      
			$(this.id + " #note_div").show(); 
		}
		$(this.id).show();
	}
	this.events = function() {
		var viewid = this.id;
		var thisView = this;
		$(viewid + " #close_btn").click(function() {
			if( noteChanged ) {
				var okFunc = function() {
					noteChanged = false;
					$(viewid).hide();
					Hub.viewNoteEditClose();
				};
				Confirm("Confirm", "Changes Not Saved, Close Anyway ?", okFunc);
			} else {
				$(viewid).hide();
				Hub.viewNoteEditClose();
			}
		});
		// --------------------------------------------
		$(viewid + " #position_btn").click(function() {
			if( DataNoteid == undefined ) {
				Notice("note must be saved before changing its position");
				return;
			}
			Views.position.display();
		});
		// --------------------------------------------
		$(viewid + " #save_btn").click(function() {
			thisView.saveNote();
		});
		// --------------------------------------------
		$(viewid + " #delete_btn").click(function() {
			if( DataNoteid == undefined ) {
				Notice("Note has not been saved.");
				return;
			}
			Confirm("Confirm", "Delete This Note", function() { thisView.deleteNote(DataNoteid); })
		});
		// --------------------------------------------
		$(viewid + " #mono_chk").change(function() {
			if($(this).prop("checked")) {
				$(viewid + " #notetext").css("font", MonoFont);
				$(viewid + " #html_div").css("font", MonoFont);
			} else {
				$(viewid + " #notetext").css("font", NoteFont);      
				$(viewid + " #html_div").css("font", NoteFont);      
			}
		});
		// --------------------------------------------
		$(viewid + " #html_chk").change(function() {
			if($(this).prop("checked")) {
				$(viewid + " #markdown_chk").prop("checked", false);
				thisView.setHtmlFirstLine();  // makes 1st line an html comment for use as note title
			} 
		});
		// --------------------------------------------
		$(viewid + " #markdown_chk").change(function() {
			if($(this).prop("checked")) {
				$(viewid + " #html_chk").prop("checked", false);
				thisView.setHtmlFirstLine();  // makes 1st line an html comment for use as note title
			} 
		});
		// --------------------------------------------
		$(viewid + " #view_chk").change(function() {
			if( !$(this).prop("checked") ) {
				$(viewid + " #html_div").hide();      
				$(viewid + " #note_div").show();      
				return;
			}
			var htmlChecked = $(viewid + " #html_chk").prop("checked");
			var markdownChecked = $(viewid + " #markdown_chk").prop("checked");
			if(htmlChecked) {
				var html = $(viewid + " #notetext").val();
				$(viewid + " #html_div").html(html);
				$(viewid + " #note_div").hide();      
				$(viewid + " #html_div").show();
			} else if(markdownChecked) {
				var markdown = $(viewid + " #notetext").val();
				var html = marked(markdown);
				$(viewid + " #html_div").html(html);
				$(viewid + " #note_div").hide();      
				$(viewid + " #html_div").show();
			} 
		});
		// --------------------------------------------
		//  handle tab keys in note text
		// --------------------------------------------
		$(document).delegate('#notetext', 'keydown', function(e) {
			noteChanged = true;
			var keyCode = e.keyCode || e.which;
			if ( keyCode == 9 ) {  // tab key
				e.preventDefault();
				var text = $(this).val();
				var domElement = $(this).get(0);
				var start = domElement.selectionStart;  // if no selection, position before cursor
				var end = domElement.selectionEnd;      // if no selection, position after cursor
				// set textarea value to: text before caret + tab + text after caret
				$(this).val(text.substring(0, start) + "\t" + text.substring(end));
				// put caret at right position
				domElement.selectionStart = domElement.selectionEnd = start + 1;
			}
		});
	}
	this.setHtmlFirstLine = function() {
		var content = $(this.id + " #notetext").val();
		var	lineFeed = content.indexOf("\n");
		if( lineFeed == -1 ) {
			var firstLine = content;
		} else {
			var firstLine = content.substring(0, lineFeed);
		}
		if( firstLine.indexOf("<!--") == -1 ) {
			$(this.id + " #notetext").val("<!-- Note Description Here -->\n" + content);
		}
	}
	this.saveNote = function() {
		var path, httpMethod;
		var requestData = {
			Content: 	$(this.id + " #notetext").val(),
			Mono: 		$(this.id + " #mono_chk").prop("checked"),
			Html: 		$(this.id + " #html_chk").prop("checked"),
			Markdown: 	$(this.id + " #markdown_chk").prop("checked"),
			Previd: 	'',    // only used for add, position
		}
		if(DataNoteid) {
			httpMethod = "PUT";		// change note
			path = "note/" + DataToken + "/" + DataBookid + "/" + DataTabid + "/" +  DataNoteid;
		} else {
			httpMethod = "POST";	// add note
			path = "note/" + DataToken + "/" + DataBookid + "/" + DataTabid;
			requestData.Previd = DataLastNoteid;  // for now add as last note, can be changed later
		}
		var settings = {
			url: UrlPrefix + path,
			method: httpMethod,
			data: JSON.stringify(requestData),
			dataType: "text",  // returned from server
		}
		$(this.id + " button").prop("disabled",true);

		var viewid = this.id;
		$.ajax(settings)
			.done(function(response) {
				noteChanged = false;
				if(httpMethod == "POST") {  // new note
					DataNoteid = response;
					DataNotes[DataNoteid] = {};
					DataAddPrevid(DataNoteid, requestData.Previd);  // see data.js
					DataBuildNoteOrder();
				}
				// save screen values to Data
				DataNotes[DataNoteid].content = requestData.Content;
				DataNotes[DataNoteid].when = new Date();
				DataNotes[DataNoteid].mono = requestData.Mono;
				DataNotes[DataNoteid].html = requestData.Html;
				DataNotes[DataNoteid].markdown = requestData.Markdown;
				DataNotes[DataNoteid].title = DataNoteTitle(DataNoteid);
				$(viewid + " #note_title").text(DataNotes[DataNoteid].title);

				Notice("Note Saved Successfully");
			})
			.fail(function(xhr) {
				Err("Note Save Failed \n" + xhr.responseText);
			})
			.always(function() {
				$(viewid + " button").prop("disabled",false);
			})
	}
	this.deleteNote = function(deleteid) {
		var viewid = this.id;
		var path = "note/" + DataToken + "/" + DataBookid + "/" + DataTabid + "/" +  deleteid;
		var settings = {
			url: UrlPrefix + path,
			method: "DELETE",
			dataType: "text",  // returned from server
		}
		$(viewid + " button").prop("disabled",true);

		$.ajax(settings)
			.done(function(response) {
				DataDeletePrevid(deleteid);  // changes previd on note following deleteid
				delete DataNotes[deleteid];
				DataBuildNoteOrder();
				Notice(response);
				Hub.noteDeleted();
			})
			.fail(function(xhr) {
				Err("Note Delete Failed \n" + xhr.responseText);
			})
			.always(function() {
				$(viewid + " button").prop("disabled",false);
			})
	}
}
