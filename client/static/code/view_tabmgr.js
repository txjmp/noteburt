
function ViewTabMgr() {
	this.id = "#view_tabmgr";
	this.html = [
		zdiv, {id:"top_div", w:"100%", h:"50px", bbottom:"1px solid black"},
		zh2, {id:"book_name", text:"Notebook Name", w:"70%", flt:"left",
			f:HeadFont, c:HeadColor,
			mleft:"30px", mtop:"10px", end:"tag"},
		zbtn, {id:"close_btn", text:"Close", flt:"right", m:"8px", end:"tag"},
		zend, zdiv,

		zdiv, {id:"options_div", w:"320px", h:"300px", flt:"left"},
		zditto, '<label><input type="checkbox" id="newtab_chk" checked>Create New Tab</label>',
		zbr, zbr,
		zlbl, {id:"curtabs_lbl", text:"Pick A Tab To Change :", mleft:"30px", mtop:"10px", end:"tag"},
		zselect, {id:"select_tab", mtop:"10px", mleft:"30px", w:"250px",
			size:"1", bkg:"white", f:"1.1em sans-serif", end:"tag"},
		zend, zdiv,

		zdiv, {id:"info_div", w:"390px", h:"300px", flt:"left"},
		zh2, {id:"info_head", text:"Enter New Tab Info", w:"100%",
			f:"600 1.4em muli, sans-serif", textalign:"center",
			mtop:"20px", end:"tag"},
		zlbl, {text:"Tab Number", end:"tag"},
		zinp, {id:"tab_number", typ:"text", w:"50px", mright:"100px", textalign:"center"},
		zlbl, {text:"Tab Name", end:"tag"},
		zinp, {id:"tab_name", typ:"text", w:"240px", pleft:"10px", end:"tag"},
		zlbl, {text:"Hidden", end:"tag"},
		zinp, {id:"hidden", typ:"checkbox", mright:"100px",},
		zbtn, {id:"save_btn", text:"Save New Tab", val:"new",
			w:"50%", mleft:"25%", mtop:"30px", end:"tag"},
		zend, zdiv,
	];
	this.css = [
		{ selector:"label",
			settings:{f:"600 1.2em muli, sans-serif", c:"black"}
		},
		{ selector:"button",
			settings: {f:"600 1.2em muli, sans-serif", c:"black", bkg:"transparent",
			b:"1px solid black", bradius:"4px", pad:"5px"}
		},
		{ selector:"#options_div input[type='checkbox']",
			settings: {mtop:"25px", mright:"7px", mleft:"30px"}
		},
		{ selector:"#info_div label",
			settings: {mtop:"15px", flt:"left", w:"125px", textalign:"right", mright:"10px"}
		},
		{ selector:"#info_div input",
			settings: {f:"700 1.2em muli, sans-serif", c:"black", mtop:"15px", flt:"left"}
		},
	]
	this.build = function() {
		var viewName = this.id.substring(1);
		$('body').append('<div id="' + viewName + '" class="view"></div>');

		$(this.id).css({
			"width":"740px", "height":"350px", 
			"z-index":"2", "position":"absolute", "top":"100px", "left":"20%"
		});
		var html = GenHtml(this.html);
		$(this.id).html(html);

		var css = GenCss(this.css);
		ApplyCss(this.id, css);

		this.events();
	}
	this.display = function(refresh) {
		$(this.id).css("background", Green1);
		$(this.id).css("border", "1px solid black");

		if(refresh)
			this.loadSelectTab();

		$(this.id + " #book_name").text(DataBookName);
		$(this.id + " #newtab_chk").prop("checked",true);
		$(this.id + " #select_tab").attr("disabled",true);
		$(this.id + " #info_head").text("Enter New Tab Info");
		$(this.id + " #save_btn").text("Save New Tab");
		$(this.id + " #tab_number").val("");
		$(this.id + " #tab_name").val("");
		$(this.id + " #save_btn").val("new");
		$(this.id + " #hidden").prop("checked",false);

		$(this.id).show();
	}
	this.loadSelectTab = function() {
		var selectTab = $(this.id + " #select_tab");
		var tabOrder = Object.keys(DataTabs).sort(function(a,b) {
			return DataTabs[a].tabNumber - DataTabs[b].tabNumber;
		});
		var html = "";
		var tab, optText;
		tabOrder.forEach(function(id) {
			tab = DataTabs[id];
			optText = tab.tabNumber + " " + tab.tabName;
			if(tab.hidden) 
				optText = '(x)' + optText;
			html += '<option value="' + id + '">' + optText + '</option>\n'; 
		});
		selectTab.html(html);
	}
	this.events = function() {
		var viewid = this.id;
		var thisView = this;
		$(viewid + " #newtab_chk").click(function() {
			if($(this).prop("checked")) {
				$(viewid + " #info_head").text("Enter New Tab Info");
				$(viewid + " #save_btn").text("Save New Tab");
				$(viewid + " #tab_number").val("");
				$(viewid + " #tab_name").val("");
				$(viewid + " #save_btn").val("new");
				$(viewid + " #hidden").prop("checked", false);
				$(viewid + " #select_tab").attr("disabled",true);
			} else if(Object.keys(DataTabs).length > 0) {
				var tabid = $(viewid + " #select_tab").val();
				$(viewid + " #info_head").text("Enter Tab Changes");
				$(viewid + " #save_btn").text("Save Changes");
				$(viewid + " #save_btn").val("change");
				$(viewid + " #tab_number").val(DataTabs[tabid].tabNumber);
				$(viewid + " #tab_name").val(DataTabs[tabid].tabName);
				$(viewid + " #hidden").prop("checked", DataTabs[tabid].hidden);
				$(viewid + " #select_tab").attr("disabled",false);
			}
		});
		$(viewid + " #select_tab").change(function() {
			if( $(viewid + " #newtab_chk").prop("checked") )
				return;
			var tabid = $(this).val();
			$(viewid + " #tab_number").val(DataTabs[tabid].tabNumber);
			$(viewid + " #tab_name").val(DataTabs[tabid].tabName);
			$(viewid + " #hidden").prop("checked", DataTabs[tabid].hidden);
		});
		$(viewid + " #save_btn").click(function(e) {
			var saveMode = $(e.target).val();
			thisView.saveTab(saveMode);
		});
		$(viewid + " #close_btn").click(function() {
			$(viewid).hide();
			Hub.viewTabMgrClose();
		});
	}
	this.saveTab = function(saveMode) {
		var viewid = this.id;
		var tabNumber = $(viewid + " #tab_number").val();
		if(isNaN(parseInt(tabNumber))) {
			Warning("Invalid Tab Number");
			return;
		}
		var tabName = $(viewid + " #tab_name").val();
		if(tabName.length == 0) {
			Notice("Tab Name is Required");
			return;
		}
		var requestData = {
			TabNumber: 	parseInt(tabNumber),
			TabName: 	tabName,
			Hidden:		$(viewid + " #hidden").prop("checked")
		}
		var path, httpMethod, tabid;
		if(saveMode == "new") {
			httpMethod = "POST";
			path = "tab/" + DataToken + "/" + DataBookid;
		} else {
			httpMethod = "PUT";
			tabid = $(viewid + " #select_tab").val();
			path = "tab/" + DataToken + "/" + DataBookid + "/" + tabid;
		}
		var settings = {
			url: UrlPrefix + path,
			method: httpMethod,
			data: JSON.stringify(requestData),
			dataType: "text",  // returned from server
		}
		$(viewid + " button").prop("disabled",true);
		$.ajax(settings)
			.done(function(response) {
				if(httpMethod == "POST") {  // for new tab, otherwise tabid is set above where saveMode != "new"
					tabid = response;
				} 
				DataTabs[tabid] = {
					"tabNumber":requestData.TabNumber,
					"tabName":	requestData.TabName,
					"hidden":	requestData.Hidden
				};
				Notice("Tab Update Successful");
			})
			.fail(function(xhr) {
				Err("SaveTab Failed \n" + xhr.responseText);
			})
			.always(function() {
				$(viewid + " button").prop("disabled",false);
			})
	}
}