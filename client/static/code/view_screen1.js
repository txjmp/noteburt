function ViewScreen1() {
	this.id = "#view_screen1";
	this.html = [
		zdiv, {id:"div1", w:"100%", h:"200px", flt:"left"},
		zh1, {text:"NoteBurt", flt:"left", mleft:"30px",
			f:"5em " + LogoFont, c:"white", textshadow:"3px 3px 5px black",
			end:"tag"},
		zlbl, {id:"release", text:"Test Version",
			f:BaseFont, c:"red", textshadow:"1px 1px 2px black",
			position:"absolute", left:"440px", top:"155px", end:"tag"},
		zend, zdiv,

		zdiv, {id:"div2", w:"150px", h:"300px", mleft:"30px", flt:"left"},
		zbtn, {id:"edit_btn", clas:"modebtn", text:"Edit", end:"tag"},
		zbtn, {id:"view_btn", clas:"modebtn", text:"View", end:"tag"},
		zbtn, {id:"create_btn", clas:"modebtn", text:"Create", end:"tag"},
		zbtn, {id:"info_btn", text:"Info", end:"tag"},
		zend, zdiv,

		zdiv, {id:"start_div", clas:"input_div"},
		zh1, {id:"start_text", text:"Share What You Know", mtop:"70px",
			f:"3em muli, sans-serif", c:Gray1, textshadow:"2px 2px 3px black",
			end:"tag"},
		zend, zdiv,

		zdiv, {id:"edit_div", clas:"input_div", disp:"none"},
		zh2, {text:"Enter Notebook Name & Access Code", end:"tag"},
		zlbl, {text:"Notebook Name", end:"tag"},
		zinp, {id:"edit_nb_name", typ:"text", w:"300px"},
		zlbl, {text:"Access Code", end:"tag"},
		zinp, {id:"edit_access_code", typ:"text", w:"120px"},
		zbr,
		zbtn, {text:"Open", clas:"open_btn", end:"tag"},
		zend, zdiv,

		zdiv, {id:"view_div", clas:"input_div", disp:"none"},
		zh2,  {text:"Enter Notebook Name", end:"tag"},
		zlbl, {text:"Notebook Name", end:"tag"},
		zinp, {id:"view_nb_name", typ:"text", w:"300px"},
		zbtn, {text:"Open", clas:"open_btn", end:"tag"},
		zend, zdiv, 

		zdiv, {id:"create_div", clas:"input_div", disp:"none"},
		zh2,  {text:"Create Notebook", end:"tag"},
		zlbl, {text:"Notebook Name", end:"tag"},
		zinp, {id:"create_nb_name", typ:"text", w:"300px"},
		zlbl, {text:"Access Code", end:"tag"},
		zinp, {id:"create_access_code", typ:"text", w:"140px", disabled:"disabled",},
		zbtn, {id:"get_accesscode_btn", text:"Get One", w:"120px", mleft:"15px", end:"tag"},
		zbtn, {text:"Open", clas:"open_btn", end:"tag"},
		zend, zdiv,
	
		zp, {id:"copyright", position:"absolute", bottom:"5px", left:"30px",
			text:"Copyright &copy 2015 Jay Poss", end:"tag"}
	]
	this.css = [
	  { selector:"#div2 button",
		settings: {f:"500 1.5em muli, sans-serif", c:"white", bkg:"transparent", b:"1px solid white",
		bradius:"4px", w:"120px", mtop:"20px", pad:"5px"}
	  },
	  { selector: ".input_div", 
		settings:{ w:"600px", h:"350px", flt:"left", mleft:"10px"}
	  },
	  { selector:".input_div button",
		settings: {f:"500 1.5em muli, sans-serif", c:"white", bkg:"transparent", b:"1px solid white",
		bradius:"4px", mbottom:"20px", pad:"4px" }
	  },
	  { selector:".input_div label",
		settings:{textalign:"right", mright:"10px", w:"180px", 
		f:"600 1.3em muli, sans-serif", mbottom:"20px", disp:"inline-block"}
	  },
	  { selector:".input_div input[type='text']",
		settings:{textalign:"center", mbottom:"20px", f:"600 1.2em muli, sans-serif"}
	  },
	  { selector:"#create_access_code",
		settings:{f:"1.5em monospace"}
	  },
	  { selector:".input_div h2", 
		settings:{f:"600 1.6em muli, sans-serif", c:"black", opacity:".5",
		mbottom:"30px", textalign:"center", w:"100%"}
	  },
	  { selector:".open_btn",
		settings:{w:"28%", mleft:"36%", mtop:"50px"}
	  },
	]
	this.build = function() {
		var viewName = this.id.substring(1);
		$('body').append('<div id="' + viewName + '" class="view"></div>');

		var html = GenHtml(this.html);
		$(this.id).append(html);

		var css = GenCss(this.css);
		ApplyCss(this.id, css);

		this.events();
	}
	this.display = function() {
		if(DirectOpen) {  // DirectOpen defined in main.js, if true: user specified book in url
			OpenMode = ViewMode;  // set global
			$(this.id + " #view_nb_name").val(LoadBookName);
			this.login();
			return;
		}
		$("body").css("background", Green1);
		$(this.id).show();
	}
	this.events = function() {
		var viewid = this.id;
		var thisView = this;
		$(viewid + " .modebtn").click(function(e) {
			$(viewid + " .modebtn").css("background-color", Green1);
			$(viewid + " .modebtn").css("color", "white");
			$(e.target).css("background-color", "white");
			$(e.target).css("color", "black");		
		});
		$(viewid + " #edit_btn").click(function() {
			$(viewid + " .input_div").hide();
			$(viewid + " #edit_div").show();
			OpenMode = EditMode;
		});
		$(viewid + " #view_btn").click(function() {
			$(viewid + " .input_div").hide();
			$(viewid + " #view_div").show();
			OpenMode = ViewMode;
		});
		$(viewid + " #create_btn").click(function() {
			$(viewid + " .input_div").hide();
			$(viewid + " #create_div").show();
			$(viewid + " #create_access_code").val("");
			OpenMode = CreateMode;
		});
		$(viewid + " #info_btn").click(function() {
			window.open("/info.html", "NoteBurt Info", "width=850, height=500, top=50, left=50");
		});
		$(viewid + " #get_accesscode_btn").click(function() {
			thisView.getAccessCode($(viewid + " #create_access_code"));
		});
		$(viewid + " .open_btn").click(function() {
			if( OpenMode == CreateMode ) {
				Confirm( "Reminder","Please Record Notebook Name & Access Code Before Opening", function() {
					thisView.login();
				});
			} else
				thisView.login();
		}); 
	}
	this.login = function() {
		var viewid = this.id;
		var thisView = this;
		var method, path;
		switch(OpenMode) {
			case EditMode:
				accessCode = $(viewid + " #edit_access_code").val();
				bookName = $(viewid + " #edit_nb_name").val();
				method = "POST";
				path = "open/" + bookName;
				break;
			case CreateMode:
				accessCode = $(viewid + " #create_access_code").val();
				bookName = $(viewid + " #create_nb_name").val();
				method = "POST";
				path = "create/" + bookName;
				break;
			case ViewMode:
				accessCode = "";
				bookName = $(viewid + " #view_nb_name").val();
				method = "GET";
				path = "view/" + bookName;
				break;
			default:
				alert("Program Error: Invalid OpenMode - " + OpenMode);
				return;
		}
		if(OpenMode != ViewMode) {
			if(accessCode.length != 7) {
				Warning("Access Code Not Valid");
				return;
			}
		}
		if(bookName == "") {
			Warning("Book Name Required");
			return;
		}	  
		var requestData = {
			"AccessCode": accessCode,
		}
		$(viewid + " button").prop("disabled", true);
		var settings = {
			url: UrlPrefix + path,
			method: method,
			data: JSON.stringify(requestData),
			dataType: "json",  // returned from server
		}
		$.ajax(settings)
			.done(function(response) {
				thisView.loginSuccess(response)
			})
			.fail(function(xhr) {
				Warning("Open Failed: " + xhr.responseText);
			})
			.always(function() {
				$(viewid + " button").prop("disabled",false);
			})
	  
		trace("login " + settings.url);
		trace(settings.data);
	}
	this.loginSuccess = function(response) {
		trace("loginSuccess");
		DataToken = response.Token;
		DataBookid = response.Bookid;
		DataBookName = response.BookName;
		DataTabs = {};
		var tabNumber, tabName, hidden;
		trace("login response tabs");
		for(tabid in response.Tabs) {
			tabNumber = response.Tabs[tabid].TabNumber;
			tabName = response.Tabs[tabid].TabName;
			hidden = response.Tabs[tabid].Hidden;
			trace(tabNumber + " " + tabName + " " + hidden);
			DataTabs[tabid] = {"tabNumber":tabNumber, "tabName":tabName, "hidden":hidden};
		}
		if(OpenMode == CreateMode)
			OpenMode = EditMode;
		//ShowBroadcast(response.Broadcast);
		Hub.openBook();
	}
	this.getAccessCode = function(domElement) {
		var viewid = this.id;
		var settings = {
			url:	UrlPrefix + "accesscode",
			method: "GET",
			dataType: "text",   // returned from server
		};
		$(viewid + " button").prop("disabled",true);
		$.ajax(settings)
			.done(function(response) {
				domElement.val(response);
				Notice("Please Record Your Access Code & Notebook Name, Thanks !");
			})
			.fail(function(xhr) {
				alert( xhr.responseText );
			})
			.always(function() {
				$(viewid + " button").prop("disabled",false);
			})
	}
}
/*
var sampleLoginResponse = {
  "Token": "abc1defghi",
  "Bookid": "000000000001",
  "BookName": "Misc Ramblings",
  "Broadcast": "broadcast msg",
  "Tabs":{
    "0000000010": {"TabNumber":10, "TabName":"January"},
    "0000000011": {"TabNumber":20, "TabName":"February"},
    "0000000014": {"TabNumber":50, "TabName":"May"},
    "0000000015": {"TabNumber":60, "TabName":"June"},
    "0000000016": {"TabNumber":70, "TabName":"July"},
    "0000000012": {"TabNumber":30, "TabName":"March"},
    "0000000013": {"TabNumber":40, "TabName":"April"},
  }
}

type loginResponse struct {
	Token     string
	Bookid    string
	BookName  string
	Broadcast string
	Tabs      map[string]data.TabInfo
}
*/

