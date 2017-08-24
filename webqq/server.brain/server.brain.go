package main 

import(
	//bulit-in
	"net/http"
	"time"
	"encoding/xml"
	"encoding/json"
	"os"
 	"io/ioutil"
 	"fmt"
 	"strings"

	//extends
  	//"server.brain/libs/ini"
  	//"SEP.DataListener/libs/ini"
)

 



//biz here
//============================================================================


//from pgsql
//============================================================================


//============================================================================



//from xml answer
//do xml need here
//============================================================================

type base struct{
	XMLName xml.Name `xml:""`
	Desc string `xml:"desc,attr"`	
	Info string `xml:",innerxml"`
}


type tools struct{
	XMLName xml.Name `xml:""`
	Desc string `xml:"desc,attr"`	
	Info string `xml:",innerxml"`
}


type techs struct{
	XMLName xml.Name `xml:""`
	Desc string `xml:"desc,attr"`	
	Info string `xml:",innerxml"`
}


type plays struct{
	XMLName xml.Name `xml:""`
	Desc string `xml:"desc,attr"`	
	Info string `xml:",innerxml"`
}

type words struct{
	Base base `xml:"base"` 
	Tools tools `xml:"tools"`
	Techs techs `xml:"techs"`
	Plays plays `xml:"plays"`
}

type emotions struct {
	Happy string `xml:"happy"`
	Anguary string `xml:"anguary"`
	Sad string `xml:"sad"`
	Cute string `xml:"cute"`
	Amaze string `xml:"amaze"`
	Helpless string `xml:"helpless"`
}

type robot struct {

	XMLName    xml.Name `xml:"robot"`
	Words words `xml:"words"`
	Emotions emotions `xml:"emotions"` 
}

type Recurlyservers struct {
    //XMLName     xml.Name `xml:"webqq"`
    //Version     string   `xml:"version,attr"`
    Robot         robot `xml:"robot"`//[]
    //Description string   `xml:",innerxml"`
}

func fnThinkWord(sQuestion string, oMapWords words) string{



	var oAllWords = map[string]string{
		oMapWords.Base.Desc : oMapWords.Base.Info,
		oMapWords.Tools.Desc : oMapWords.Tools.Info,
		oMapWords.Techs.Desc : oMapWords.Techs.Info,
		oMapWords.Plays.Desc : oMapWords.Plays.Info,
	}


	var oHelp []string
	oSecHelp := make(map[string]string)

	for k, v := range oAllWords {


		oHelp = append(oHelp, k)

		var oItemMap map[string]interface{}
		if err := json.Unmarshal([]byte(v), &oItemMap); err != nil {

	        return ""
	    }


	    var oSecHelpItem []string
	    for iK, iV := range oItemMap {

	    	oSecHelpItem = append(oSecHelpItem, iK)

	    	if strings.Index(sQuestion, iK) >= 0 {

	    		return iV.(string)
	    	}
	    }

	    oSecHelp[k] = strings.Join(oSecHelpItem, " | ")

    }

    if "help" == sQuestion {
    	return strings.Join(oHelp, " | ")
    }

    for k, v := range oSecHelp {

    	if strings.Index(sQuestion, k) >= 0 {

    		return v
    	}
    }


    return "【机器人】┌( ಠ_ಠ)┘ 不懂怎么回答你哦~"
}

func fnBiz(sQuestion string) string{
 		
 	//sQuestion

 	oXMLFile, err := os.Open("./conf/server.brain.xml") // For read access.     

    if err != nil {

        fmt.Printf("error: %v", err)
        return ""
    }
	    
    defer oXMLFile.Close()

    oData, err := ioutil.ReadAll(oXMLFile)

    if err != nil {
        fmt.Printf("error: %v", err)
        return ""
    }


    v := Recurlyservers{}
    err = xml.Unmarshal(oData, &v)
    if err != nil {
        fmt.Printf("error: %v", err)
        return ""
    }


    //fmt.Println(v.Robot.Words)


//    var oMapWords words //map[string]interface{}

    // if err := json.Unmarshal([]byte(v.Robot.Words), &oMapWords); err != nil {

    //     return ""
    // }

    //fmt.Println(oMapWords)

    //sQuestion 
    return fnThinkWord(sQuestion, v.Robot.Words) 
    //fmt.Println(v)
}

//=====================================================================================================================


func init(){
 	

}


func main() {

	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request){

		//parse querydate
		r.ParseForm()

		sId := r.Form["sid"][0]
		sName := r.Form["n"][0]
		sQuestion := r.Form["q"][0]


		//read from ini file
		//oIni = ini.GetSecMap("webqq.robot.words", &Robot_Words{}).(*Robot_Words)

		w.Write([]byte(fnBiz(sQuestion)))


	})

	oHttpServer := &http.Server{
		Addr : ":8521",
		Handler : nil,
		ReadTimeout: time.Duration(20) * time.Second,
		WriteTimeout: time.Duration(20) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//start listen and server
	oHttpServer.ListenAndServe()
}