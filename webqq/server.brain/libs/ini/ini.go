package ini

import(
	//Built-in
	"os"
	// "io/ioutil"
	// "bytes"
	//"fmt"

	//extends
	. "bakerstreet-club/logs"
	"github.com/bakerstreet-club/ini"
)

const(

	//the config file path
	FILE_PATH = "conf/SEP.DataListener.ini"
)



var(
	cfg = ini.Empty()
)
//init read
func init() {
	//export data change
	//cfg, err := ini.Load([]byte("type = pqsql"), FILE_PATH, ioutil.NopCloser(bytes.NewReader([]byte("type = pqsql"))))
	//find file if another file no find
	//cfg, err := ini.LooseLoad(FILE_PATH, FILE_PATH_404)
	//opinion set
	//cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, FILE_PATH)
	//load key to lower get
	var err error
	cfg, err = ini.InsensitiveLoad(FILE_PATH)
	if nil != err {
		Error(err.Error())
		os.Exit(1)
	}
}
//get config with exactly info of section and key
func GetConfig(sSec, sKey string) (*ini.Key, error) {

	//get section first
	oSec, err := cfg.GetSection(sSec)
	if nil != err {
		Error(err.Error())
		return nil, err
	}
	//get key and check error
	oVal, err := oSec.GetKey(sKey)
	if nil != err {
		Error(err.Error())
		return nil, err
	}
	return oVal, nil
}
//same section kv
func GetSecMap(sSec string, oMapObj interface{}) interface{} {

	//thing the struct first
	if nil == oMapObj {
		Error("sorry, the func named GetSecMap in libs/ini that two param is can not be nil, please check error and try again ")
		return nil
	}
	//get section first
	oSec, err := cfg.GetSection(sSec)
	if nil != err {
		Error(err.Error())
		return nil
	}
	//kv in section map to map
	err = oSec.MapTo(oMapObj)
	if nil != err {
		Error(err.Error())
		return nil
	}
	return oMapObj
}
//get the section then config
func GetTheSection(sSec string) *ini.Section {

	//get the section
	oSec, err := cfg.GetSection(sSec)
	if nil != err {
		Error(err.Error())
		return nil
	}
	return oSec
}
func GetSecConfig(oSec *ini.Section, sKey string) interface{} {
	//think first
		if nil == oSec {
				Error("sorry, the section can not be empty !")
				return nil
		}
		//get the config value
		oVal, err := oSec.GetKey(sKey)
		if nil != err {
			Error(err.Error())
			return nil
		}
		return oVal.Value()
}
