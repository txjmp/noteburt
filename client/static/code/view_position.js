function ViewPosition() {
	this.id = "#view_position";
	this.html = [
		zdiv, {id:"top_div", w:"100%", h:"50px", bbottom:"1px solid black", flt:"left"},
		zh2, {id:"note_title", text:"Note Title Here", w:"70%", flt:"left",
			f:"600 1.2em muli, sans-serif", c:"black",
			mleft:"20px", mtop:"12px", end:"tag"},
		zbtn, {id:"close_btn", text:"Close", flt:"right", m:"8px", end:"tag"},
		zend, zdiv,

		zdiv, {id:"select_div", w:"380px", h:"220px", flt:"left"},
		zlbl, {id:"select_head", text:"Current Note List",
			w:"100%", textalign:"center", end:"tag"},
		zselect, {id:"select_note", mtop:"10px", mleft:"10px", w:"90%",
			size:"1", bkg:"white", f:"1.1em sans-serif", end:"tag"},
		zend, zdiv,

		zdiv, {id:"move_div", w:"330px", h:"220px", flt:"left"},
		zlbl, {id:"move_head", text:"Select Note to Position This Note After",
			w:"100%", textalign:"center", end:"tag"},
		zditto, '<label><input type="checkbox" id="firstone_chk"/>Move To Top Of List</label><br>',
		zbtn, {id:"move_btn", text:"Move Note",	w:"50%", mleft:"25%", mtop:"30px", end:"tag"},
		zend, zdiv,
	];
	this.css = [
		{ selector:"button",
		  settings: {f:BaseFont, c:"black", bkg:"transparent", b:"1px solid black", bradius:"4px", pad:"5px"}
		},
		{ selector:"label",
		  settings: {f:BaseFont, c:"black", mtop:"15px", disp:"block"}
		},
		{ selector:"#move_div input[type='checkbox']",
		  settings: {mtop:"15px", mleft:"70px", mright:"7px"}
		}
	]
	this.build = function() {
		var viewName = this.id.substring(1);
		$('body').append('<div id="' + viewName + '" class="view"></div>');

		$(this.id).css({
			"width":"750px", "height":"300px", 
			"z-index":"2", "position":"absolute", "top":"100px", "left":"15%"
		});
		var html = GenHtml(this.html);
		$(this.id).html(html);

		var css = GenCss(this.css);
		ApplyCss(this.id, css);
		
		this.events();
	}
	this.events = function() {
		var viewid = this.id;
		var thisView = this;
		$(viewid + " #move_btn").click(function() {
			var previd;
			if( $(viewid + " #firstone_chk").prop("checked") )
				previd = DataZeroid;
			else 
				previd = $(viewid + " #select_note").val();  // id of note this note (id=DataNoteid) is being moved after  
			if( previd == undefined || previd == DataNoteid ) {
				Notice("Invalid Choice - No Affect On Order");
				return;        // can't position note after itself
			}
			thisView.positionNote(previd);
		});

		$(viewid + " #close_btn").click(function() {
			$(viewid).hide();
		});
	}
	this.display = function() {
		$(this.id).css("background", Green1);
		$(this.id).css("border", "1px solid black");

		$(this.id + " #note_title").text(DataNotes[DataNoteid].title);

		var selectNote = $(this.id + " #select_note");
		var noteCount = Object.keys(DataNotes).length;
		var id;
		var html = "";
		var nextid = DataZeroid;
		for(var i=0; i < noteCount; i++) {
			id = DataNoteOrder[nextid];
			html += '<option value="' + id + '">' + DataNotes[id].title + '</option>\n'; 
			nextid = id;
		}
		selectNote.html(html);

		$(this.id + " #book_name").text(DataBookName);
		$(this.id).show();
	} 
	this.positionNote = function(previd) {
		var viewid = this.id;
		DataDeletePrevid(DataNoteid);		// changes previd of note following this note
		DataAddPrevid(DataNoteid, previd);
		DataBuildNoteOrder();

		if( NoServerTesting ) return;

		var requestData = {
			Previd: previd,
		}
		var path = "positionnote/" + DataToken + "/" + DataBookid + "/" + DataTabid + "/" + DataNoteid;
		var settings = {
			url: UrlPrefix + path,
			method: "PUT",
			data: JSON.stringify(requestData),
			dataType: "text",  // returned from server
		}
		$(viewid + " button").prop("disabled",true);
		$.ajax(settings)
			.done(function(data) {
				$(viewid).hide();
				Notice(data);
			})
			.fail(function(xhr) {
				Err("Position Change Failed \n" + xhr.responseText);
			})
			.always(function() {
				$(viewid + " button").prop("disabled",false);
			})
	}
}
