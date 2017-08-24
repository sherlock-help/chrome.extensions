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

type robot struct {

	XMLName    xml.Name `xml:"robot"`
	Words string   `xml:"words"`
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


    var oMapWords map[string]interface{}

    if err := json.Unmarshal([]byte(v.Robot.Words), &oMapWords); err != nil {

        return ""
    }

    //fmt.Println(oMapWords)

    //sQuestion
    for k, v := range oMapWords {
    	if strings.Index(sQuestion, k) >= 0 {

    		return v.(string)
    	}
    }

    return "【机器人】┌( ಠ_ಠ)┘ 不懂怎么回答你哦~"
    //fmt.Println(v)
}


//do xml need here
//============================================================================
type Recurlyservers struct {
    //XMLName     xml.Name `xml:"webqq"`
    //Version     string   `xml:"version,attr"`
    Robot         robot `xml:"robot"`//[]
    //Description string   `xml:",innerxml"`
}
//============================================================================



func init(){
 	

}


func main() {

	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request){

		//parse querydate
		r.ParseForm()

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