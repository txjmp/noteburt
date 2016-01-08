// tag literals used in view defs, for convenience so tag names don't have to be enclosed in quotes
const ztag = 	'tag';
const zlbl = 	'label';
const zinp = 	'input';
const zh1 = 	'h1';
const zh2 = 	'h2';
const zh3 = 	'h3';
const zdiv = 	'div';
const zp = 		'p';
const zimg = 	'img';
const zbtn = 	'button';
const zul = 	'ul';
const zol = 	'ol';
const zli = 	'li';
const zend = 	'end';
const zditto = 	'ditto';  // means to copy input string as is
const zbr = 	'br';
const zselect = 'select';
const ztextarea = 'textarea';

var htmlCodes = {
	id: 		"id",
	typ: 		"type",
	type: 		"type",
	val: 		"value",
	cls: 		"class",
	clas: 		"class",
	class: 		"class",
	name: 		"name",
	chk: 		"checked",
	checked: 	"checked",
	disabled: 	"disabled",
	readonly: 	"readonly",
	rows: 		"rows",
	size: 		"size",
	target: 	"target",
	href: 		"href",
	src: 		"src",
}
var cssCodes = {
	b: 		"border",
	bright: "border-right",
	bleft:	"border-left",
	btop: 	"border-top",
	bbottom: "border-bottom",
	bradius: "border-radius",
	bkg:	 "background",
	bkgimage: "background-image",
	bkgrepeat: "background-repeat",
	c: 		"color",
	f: 		"font",
	ff: 	"font-family",
	fs: 	"font-size",
	fw: 	"font-weight",
	m: 		"margin",
	mleft: 	"margin-left",
	mright: "margin-right",
	mtop: 	"margin-top",
	mbottom: "margin-bottom",
	p:		"padding",
	pad: 	"padding",
	pleft: "padding-left",
	pright: "padding-right",
	ptop:	"padding-top",
	pbottom: "padding-bottom",
	w: 		"width",
	h: 		"height",
	flt: 	"float",
	pos: 	"position",
	position: "position",
	textalign: "text-align",
	boxshadow: "box-shadow",
	textshadow: "text-shadow",
	cursor:	"cursor",
	top:	"top",
	left:	"left",
	bottom:	"bottom",
	right:	"right",
	disp:	"display",
	opacity: "opacity",
	wordwrap: "word-wrap",
	maxwidth: "max-width",
	minheight: "min-height",
	maxheight: "max-height",
	minwidth: "min-width",
	clear:	"clear",
	overflow: "overflow",
	readonly: "readonly",
	tab:	"tab-size",
}

function GenHtml(screenDef) {
	var result = '';
	var tag, attrs, text, style, htmlAttrs, lineEnd;
	for(var i=0; i<screenDef.length; i+=2) {
		tag = screenDef[i];
		if(tag == 'end') {
			result += '</' + screenDef[i+1] + '>\n';
			continue;
		}
		if(tag == 'ditto') {
			result += screenDef[i+1] + '\n';
			continue;
		}
		if(tag == "br" ) {
			result += "<br>\n";
			i = i - 1;  // no associated value for this tag
			continue;
		}
		attrs = screenDef[i+1];
		text = '';
		style = {};		  // css attributes (margin, border, etc.)	
		htmlAttrs = {};   // html attributes (id, name, class, etc.)
		lineEnd = '';
		for(attr in attrs) {
			attrVal = attrs[attr];
			if(attr == "text") {
				text = attrVal;
				continue;
			}
			if( htmlCodes[attr] ) {
				htmlAttrs[htmlCodes[attr]] = attrVal;
				continue;
			}
			if( cssCodes[attr] ) {
				style[cssCodes[attr]] = attrVal;
				continue;
			}
			if( attr == "end" ) {
				lineEnd = attrVal;
				continue;
			}
			alert('invalid screenDef attribute: ' + attr);
		}
		result += '<' + tag;
		if(Object.keys(htmlAttrs).length > 0) {
			for(var a in htmlAttrs) {
				result += ' ' + a + '=' + '"' + htmlAttrs[a] + '"';
			}
		}
		if(Object.keys(style).length > 0) {
			result += ' style="';
			var sp = '';
			for(var a in style) {
				result += sp + a + ':' + style[a] + ';';
				sp = ' ';
			}
			result += '"';
		}
		result += '>' + text;

		if(lineEnd == 'tag') {
			result += '</' + tag + '>';
		}
		result += '\n';
	}
	return result;
}


// var shortCss = [
//   { selector: "button", settings:{ ff:"sans-serif" }}
// ]
// returns [{selector:"button",settings:{font-family:"sans-serif"}]
function GenCss(shortCss) {
	var regCss, regAttr;
	var response = [];
	shortCss.forEach(function(line) {
		regCss = {};
		for( attr in line.settings ) {
			if( !cssCodes[attr] ) {
				alert("invalid attr inCss - " + attr);
				return;
			}
			regAttr = cssCodes[attr];
			regCss[regAttr] = line.settings[attr];
		}
		response.push({selector:line.selector, settings:regCss});
	});
	return response;
}

// css is result of GenCss function
function ApplyCss(parentSelector, css) {
  //console.log("ApplyCss");
  css.forEach(function(line) {
    //console.log(parentSelector + " " + line.selector);
    //console.log(line.settings);
    $(parentSelector + " " + line.selector).css(line.settings);
  });
}

function GetCookie(cname) {
    var name = cname + "=";
    var cookies = document.cookie.split(';');
    for(var i=0; i<cookies.length; i++) {
        var cookie = cookies[i];
        while (cookie.charAt(0) == ' ')
          cookie = cookie.substring(1);  // remove leading space
        if (cookie.indexOf(name) != -1)  // compare each cookie found to one desired, if match then return value
          return cookie.substring(name.length,cookie.length);
    }
    return "";
}

function Notice(text,duration) {
  if(duration == undefined)
    var duration = 2500;
  $("#notice_div").css("background-color", "rgba(0,0,200,.9)");
  $("#notice_text").css("color", "white");
  $("#notice_text").text(text);
  $("#notice_div").show(); 
  setTimeout(function() {
    $("#notice_div").hide();
  }, duration);
}
function Warning(text,duration) {
  if(duration == undefined)
    var duration = 3000;
  $("#notice_div").css("background-color", "rgba(250,225,20,.9)");
  $("#notice_text").css("color", "black");
  $("#notice_text").text(text);
  $("#notice_div").show(); 
  setTimeout(function() { $("#notice_div").hide(); }, duration);
}

var confirm_okFunc;
var confirm_funcs_set = false;

function Confirm(title, text, okFunc) {
	confirm_okFunc = okFunc;
	if(!confirm_funcs_set) {  // run 1st time Confirm called, had problems defining click funcs other ways
		confirm_funcs_set = true;
		$("#confirm_cancel").click(function() {
			$("#confirm_div").hide();      
		});
		$("#confirm_ok").click(function() {
			$("#confirm_div").hide();
			if(confirm_okFunc != undefined)
				confirm_okFunc();
		});
	}
	$("#confirm_title").text(title);
	$("#confirm_text").text(text);
	$("#confirm_div").show();
}

function Err(msg) {
  console.log("*** ERROR ***");
  console.log(msg);
  alert("Error\n" + msg);
}

// inp = [2015, 12, 25, 9, 55, 55]  (dec 25, 2015 9:55:55)
function ArrayToDate(inp) {
	var yr = inp[0];
	var mo = inp[1] - 1;  // js month is 0 - 11
	var da = inp[2];
	var hr = inp[3];
	var min = inp[4];
	var sec = inp[5];
	return new Date(yr, mo, da, hr, min, sec, 0);
}

var fmtDateMonths = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
// convert date to Jan 1, 2016 13:44
function FmtDate(inp) {
	var yr = inp.getFullYear();
	var mo = inp.getMonth();
	var da = inp.getDate();
	var hr = inp.getHours();
	var min = inp.getMinutes();
	var mth = fmtDateMonths[mo];
	if(min < 10)
		min = "0" + min;
	return mth + " " + da + "," + " " + yr + " " + hr + ":" + min; 
}