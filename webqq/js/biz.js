
//define the object for return word
var oWords = {

	"召唤" : "【机器人】我来也~~",
	"走开" : "【机器人】轻轻的, 我走了...",
	"官网" : "【机器人】http://bakerstreet.club",
	"博客" : "【机器人】http://sherlock.help",
	"作者" : "【机器人】作者很帅",
	"乃贤是一只鸡" : "【机器人】没错，就是！"
};


var oKeyItems = [];
for(var item in oWords){
	oKeyItems.push(item);
}
oWords["help"] = oKeyItems.join("  |  ");


var oCache = {};

// do the Transaction
function fnTrans(oThis){

	//click first
	$(oThis).click();


	setTimeout(function(){

		var sId = oThis.id;

		if(!oCache[sId])
			oCache[sId] = {"status":true};


		var oBuddy = $("#panelBody-5 .chat_content_group.buddy").not(".need_update");
		
		if(oBuddy.length > 0 && oCache[sId]["listenHistory"] && $("#panelBody-5 .chat_content_group.buddy").not(".need_update").last().find(".chat_content").html() == "召唤"){
			oCache[sId]["status"] = true;
		}


		if(oCache[sId]["status"]){

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

				oNewSay.push(sSayWord ? sSayWord : "【机器人】你这样说，我都不知道怎么回你了。。。");

				//go away
				if("走开" == sListenWord){
					oCache[sId]["status"] = false;
				}

				//save to history
				oCache[sId]["listenHistory"].push(oNewListen[item]);
			}

			for(var item in oNewSay){
 
				//chat_textarea
				$("#container #chat_textarea").val(oNewSay[item]);		

				//send_chat_btn
				$("#container #send_chat_btn").click();
				//save to history
				oCache[sId]["sayHistory"].push(oNewSay[item]);	
			}
		}

		$("#panelRightButtonText-5").click();

	}, 500);	 
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
	setInterval(fnRender, 1000);	



	// not business
	//==============================================================
	//background for look up
	setInterval(function(){
		if(0 == $("#current_chat_list").length){
			location = location;
		}
	}, 3600000);//one hour
}