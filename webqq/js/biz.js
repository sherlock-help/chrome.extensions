
//define the object for return word
//default here
var oWords = {

	"召唤" : "【机器人】我来也~~",
	"走开" : "【机器人】轻轻的, 我走了...",
	"官网" : "【机器人】http://bakerstreet.club",
	"博客" : "【机器人】http://sherlock.help",
	"作者" : "【机器人】作者很帅"
};


var oKeyItems = [];
for(var item in oWords){
	oKeyItems.push(item);
}
oWords["help"] = oKeyItems.join("  |  ");

//ask from golang
var sReAnswer = "";


var oCache = {};

// do the Transaction
function fnTrans(oThis){

	//click first
	$(oThis).click();


	setTimeout(function(){

		var sId = oThis.id;

		if(!oCache[sId])
			oCache[sId] = {"status":false};


		var oBuddy = $("#panelBody-5 .chat_content_group.buddy").not(".need_update");
		
		if(oBuddy.length > 0 && oCache[sId]["listenHistory"] && $("#panelBody-5 .chat_content_group.buddy").not(".need_update").last().find(".chat_content").html() == "召唤"){
			oCache[sId]["status"] = true;
		}

		//oCache[sId]["status"]
		if(true){

			//init think
			if(!oCache[sId]["listenHistory"]){
				oCache[sId]["listenHistory"] = [];	
			}
			if(!oCache[sId]["sayHistory"]){
				oCache[sId]["sayHistory"] = [];
			}

			//friend can answer not same		
			oCache[sId]["isFriend"] = sId.indexOf("recent-item-friend") > -1;
			oCache[sId]["isGroup"] = sId.indexOf("recent-item-group") > -1;


			//panelBody-5
			//chat_content_group buddy
			//chat_content_group self
			//chat_content 
			var oNewListen = [];
			var oNewSay = [];
			 
			$("#panelBody-5 .chat_content_group.buddy").not(".need_update").each(function(iIndex){

				//say group have need_update class 
				//if(this.getAttribute("class").indexOf("need_update") > -1) return;


				var sListen = $(this).find(".chat_content").html();

				if(iIndex >= oCache[sId]["listenHistory"].length && sListen.indexOf("【机器人】") == -1){

					oNewListen.push(sListen);
				}

			});
			// $("#panelBody-5 .chat_content_group.self").each(function(iIndex){

			// 	if(iIndex >= oCache[sId]["sayHistory"].length){

			// 		oNewSay.push($(this).find(".chat_content").html());
			// 	}
				
			// });

			//chat_content_group buddy need_update 


			//save to history
			// and answer
			var bIfGoAway = oCache[sId]["status"];
			for(var item in oNewListen){

				var sListenWord = oNewListen[item];

				var sSayWord = "";

				if(oWords[sListenWord]){
					sSayWord = oWords[sListenWord];
				}else{

					//模糊匹配
					for(var key in oWords){					
						if(sListenWord.indexOf(key) > -1){
							sSayWord = oWords[key];
						}
					} 
				}


				//query back
				if(!sSayWord){
					sendAjaxRequest("http://127.0.0.1:8521/ask?q="+sListenWord);
					sSayWord = sReAnswer;
				}


				oNewSay.push(sSayWord ? sSayWord : "【机器人】你这样说，我都不知道怎么回你了。。。");

				//go away
				if("走开" == sListenWord)
					oCache[sId]["status"] = false;

				//save to history
				oCache[sId]["listenHistory"].push(oNewListen[item]);
			}

			for(var item in oNewSay){
 	
 				if(bIfGoAway){
					//chat_textarea
					$("#container #chat_textarea").val(oNewSay[item]);		

					//send_chat_btn
					$("#container #send_chat_btn").click();
 				}

				//save to history
				oCache[sId]["sayHistory"].push(oNewSay[item]);	
			}

		}

		$("#panelRightButtonText-5").click();

	}, 200);	 
}



function fnRender(){
	$("#current_chat_list .notify").each(function(){
		//recent-item-friend
		//recent-item-group
		fnTrans(this);
	});
}

 
if(location.href.indexOf("http://w.qq.com") > -1){

	//will 1s find info
	setInterval(fnRender, 800);	



	// not business
	//==============================================================
	//background for look up
	setInterval(function(){
		//2 hours after reflash
		location = location;
	}, 7200000);//one hour
}



//invoke back
//==========================================================
var XMLHttpReq;  
function createXMLHttpRequest() {  
    try {  
        XMLHttpReq = new ActiveXObject("Msxml2.XMLHTTP");//IE高版本创建XMLHTTP  
    }  
    catch(E) {  
        try {  
            XMLHttpReq = new ActiveXObject("Microsoft.XMLHTTP");//IE低版本创建XMLHTTP  
        }  
        catch(E) {  
            XMLHttpReq = new XMLHttpRequest();//兼容非IE浏览器，直接创建XMLHTTP对象  
        }  
    }  
  
}  

function sendAjaxRequest(url) {  
    createXMLHttpRequest();                                //创建XMLHttpRequest对象  
    XMLHttpReq.open("post", url, false);  //false 同步
    XMLHttpReq.onreadystatechange = processResponse; //指定响应函数  
    XMLHttpReq.send(null);
}  
//回调函数  
function processResponse() {  
    if (XMLHttpReq.readyState == 4) {  
        if (XMLHttpReq.status == 200) {  
            var sResponseText = XMLHttpReq.responseText;  
            /** 
             *实现回调 
             */  
            sReAnswer = window.decodeURI(sResponseText);
        }  
    }  
  
}  