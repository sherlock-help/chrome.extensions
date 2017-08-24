package libs

import(
    //bulit-in
    "net/http"
    "io/ioutil"
    "net/url"
    "strings"

    //extends
    . "bakerstreet-club/logs"
    GoQuery "github.com/bakerstreet-club/goquery"
//    iconv "github.com/sherlock-help/iconv-go"
)

func init(){
  oDocBuffer = make(map[string] interface{})
}
//buffer here
var oDocBuffer map[string]interface{}
//get doc
func GoQueryDoc(sURL string) interface{} {

      if _, ok := oDocBuffer[sURL]; ok{
          Info("get doc from buffer")
          return oDocBuffer[sURL]
      }
      //check first
      if "" == sURL {
          Error("sorry, the param named sURL is can not be empty!")
          return ""
      }

      //go query begin
      doc, err := GoQuery.NewDocument(sURL)
      if nil != err {
          Error(err.Error())
      }
      oDocBuffer[sURL] = doc
      return doc
}

//go query
func GoQueryByURLAndSelect(sURL, sSelector string) []string {
    doc := GoQueryDoc(sURL).(*GoQuery.Document)
    var oReText []string
    doc.Find(sSelector).Each(func(i int, s *GoQuery.Selection){
        //for Each
        oReText = append(oReText, s.Text())
    })
    return oReText
}

//param:
//sURL
//oURLBring the first is param; and second is header
func QueryByURL(sURL string, oURLBring ...map[string]string) interface{}{

     return DoPostURL(sURL, oURLBring)
}

func GetPostResponse(sURL string, oURLBring ...map[string]string) *http.Response {

    return getPostResponse(sURL, oURLBring)
}

func getPostResponse(sURL string, oURLBring []map[string]string) *http.Response {

      //check first
      if sURL == "" {
          Error("sorry, the param named sURL is can not be emtpy! ")
          return nil
      }

      //query here

      oParam := url.Values{}
      if len(oURLBring) > 0 {
        //set param
         for k, v := range oURLBring[0] {
            oParam.Set(k, v)
         }
      }
      bodyPost := ioutil.NopCloser(strings.NewReader(oParam.Encode()))
      clientPost := &http.Client{}

      rq, _ := http.NewRequest("POST",
        sURL,
        bodyPost)
      //set header
      rq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
      //rq.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
      rq.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
      rq.Header.Set("Cache-Control", "max-age=0")
      rq.Header.Set("Connection", "keep-alive")
      if len(oURLBring) > 1 {
        //set Header
        for k, v := range oURLBring[1] {
            rq.Header.Set(k, v)
        }
      }

      respPost, _ := clientPost.Do(rq)

      return respPost
}

func DoPostURL(sURL string, oURLBring []map[string]string) interface{} {

   respPost := getPostResponse(sURL, oURLBring)
   defer respPost.Body.Close()

   dataPost, _ := ioutil.ReadAll(respPost.Body)

   //change encode
   //out := make([]byte,len(dataPost))
  // iconv.Convert(dataPost, out, "gb2312", "utf-8")

   return string(dataPost)
}

func GetPostDocSelection(oDoc interface{}, sSelector string) *GoQuery.Selection {

    if nil == oDoc {
      Error("sorry, the param named oDoc is can not be nil")
      return nil
    }

    //doc goquery here
    doc, err := GoQuery.NewDocumentFromResponse(oDoc.(*http.Response))
    if nil != err {
        Error(err.Error())
        return nil
    }

    //var oReSelection []*GoQuery.Selection
    return doc.Find(sSelector)

    // .Each(func(i int, s *GoQuery.Selection){
    //     //for Each
    //     oReSelection = append(oReSelection, s)
    // })
    // return oReSelection
}
