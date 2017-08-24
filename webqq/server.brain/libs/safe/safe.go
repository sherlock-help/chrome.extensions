package safe

import (
  //bulit-in package
  //"fmt"
  "strings"
  "strconv"
  "regexp"

  //extends
  . "bakerstreet-club/logs"
  "SEP.DataListener/libs/ini"
)


//init load
func init() {

}

var (
    oSafeSec = ini.GetTheSection("safe")
)

//filter input and out put handle
func Filter(oFT interface{}) interface{} {

    //declare for return
    var oReFilter interface{}
    switch oFT.(type) {
        case string: oReFilter = handleSkinString(oFT.(string))
        default: oReFilter = nil
    }
    return oReFilter
}


//handle func
//-----------------------------------------
func handleSkinString(sQUeryString string) string  {

  //think first
  sPre := ini.GetSecConfig(oSafeSec, "safe_webio_key_pre").(string)
  sRegexLen := ini.GetSecConfig(oSafeSec, sPre + "len").(string)
  iRegexLen, err := strconv.Atoi(sRegexLen)
  if nil != err {
      Error(err.Error())
      return ""
  }
  sRegexRpLen := ini.GetSecConfig(oSafeSec, sPre + "rplen").(string)
  iRegexRpLen, err := strconv.Atoi(sRegexRpLen)
  if nil != err {
      Error(err.Error())
      return ""
  }
  for i:=1; i <= iRegexLen; i++ {
      sRegexCkey := sPre + strconv.Itoa(i)
      sRegex := ini.GetSecConfig(oSafeSec, sRegexCkey).(string)
      bMat, err := regexp.MatchString(sRegex, sQUeryString)
      if nil != err {
          Error(err.Error())
          continue
      }
      //success
      if bMat {
        for j := 1; j <= iRegexRpLen; j++ {
          sRegex_no := ini.GetSecConfig(oSafeSec, sRegexCkey + "no" + strconv.Itoa(j)).(string)
          sRegex_ok := ini.GetSecConfig(oSafeSec, sRegexCkey + "ok" + strconv.Itoa(j)).(string)
          sQUeryString = strings.Replace(sQUeryString, sRegex_no, sRegex_ok, -1)
        }
      //faild
      }
  }
  return sQUeryString
}



//-----------------------------------------
