package dao

import(
	//bulit-int package
	"database/sql"
	"strings"
	"reflect"
	//"fmt"
	"strconv"
	"encoding/json"
	"time"
	//"io/ioutil"
	"os"
	"path/filepath"

	//extends
	_ "github.com/lib/pq"
	//_ "github.com/bakerstreet-club/pgsql"
	//_ "github.com/bakerstreet-club/gopgsql"
 	. "bakerstreet-club/logs"


	. "SEP.DataListener/libs/safe"
 	"SEP.DataListener/libs/ini"
)

//when use this page need the obj
func init() {
	//if need can init do here
	CreateConnect()
}

var (
	pgdb *sql.DB
	oIni *pgsqlIni
)

//the func for create sql connect
func CreateConnect()  {

	//the func for create connect
	oIni = ini.GetSecMap("db.pgsql", &pgsqlIni{}).(*pgsqlIni)
	var err error
	pgdb, err = sql.Open("postgres", "user="+ oIni.Usr +" password="+ oIni.Pwd +" dbname="+ oIni.Dbn +" sslmode="+ oIni.Sslmode)
	if nil != err {
		Error(err.Error())
	}
}

//the db obj
func GetPgObj() *PGSql {

	//load from config
	err := pgdb.Ping()
	//fmt.Println(err)
	if nil != err {
			Error(err.Error())
			//reCreate
			CreateConnect()
	}
	// defer func(){
	// 	pgdb.Close()
	// }()
	//sql.Open("postgres", "postgres://username:password@localhost/db_name?sslmode=disable")
	oPGdb := new(PGSql)
	oPGdb.obj = pgdb
	oPGdb.ini = *oIni

	return oPGdb
}

//the object from config
type pgsqlIni struct {
	Host string
	Port string
	Dbn string
	Usr string
	Pwd string
	Sslmode string
}
type PGSql struct {
	obj *sql.DB
	ini pgsqlIni
}

// select column_name from information_schema.columns where table_name = ''
//get all Columns
func (oPgsql *PGSql) GetColumns(sTableName string) []string {

	//query obj
	oMapQuery := map[string] PGQueryObj {
		"information_schema.columns" : PGQueryObj{
			Fields : []string{"column_name"},
			ObjWhere : []PGQueryWhere {
				*&PGQueryWhere{
					Field : "table_name",
					Op : "equal",
					Val : sTableName,
				},
			},
		},
	}
	oMapResult := oPgsql.QueryData(oMapQuery)
	var oColumns []string
	for _, v := range oMapResult["information_schema.columns"] {
		oColumns = append(oColumns, v["column_name"])
	}
	return oColumns
}


//return operate status
func (oPgSql *PGSql) Insert(sTBName string,oMapKV map[string]string) bool{

	//check first
	if "" == sTBName || 0 == len(oMapKV) {
		Error("sorry, the func insert in model with param db or sTBName or oMapKV is not satisfy condition ! ")
		return false
	}

	//obj
	oFields := []string{}
	oValues := []string{}
	iFieldIndex := 1
	var sValueOcc string
	for k, v := range oMapKV {
		oFields = append(oFields, Filter(k).(string))
		oValues = append(oValues, Filter(v).(string))

		sValueOcc += "$" + strconv.Itoa(iFieldIndex) + ","
		iFieldIndex++
	}
	//just for safe
	sTBName = Filter(sTBName).(string)
	if len(oFields) == 0 {
			Error("sorry, the fields length for pgsql insert is zero !")
			return false;
	}
	sInsSql := `INSERT INTO ` + sTBName + `("` + strings.Join(oFields, `","`) + `") VALUES(` + sValueOcc[0: len(sValueOcc) - 1] + `)`
 	//log.Println(sInsSql)
	//Info(sInsSql)
	stmt, err := oPgSql.obj.Prepare(sInsSql)
	defer stmt.Close()

	if err != nil {
		Error(err.Error())
		return false
	}
	//	stmt.Exec("test", "test")
	oInValues := make([]interface{}, len(oValues))
	for i, s  := range oValues {
		//Info(s)
		oInValues[i] = s
	}
	_, errDo := stmt.Exec(oInValues...)
	if errDo != nil {
		Error(errDo.Error())
		return false
	}
	//LastInsertId
	return true
}

//return operate status
func (oPgSql *PGSql) Update(sTBName string, oMapKV map[string]string, oKeyKV map[string]string) bool{

	if "" == sTBName || 0 == len(oKeyKV) || 0 == len(oMapKV) {
		Error("sorry, can not find the sTBName or len(oKeyKV) == 0 or len(oMapKV) == 0 for update the table in pqsql database ! ")
		return false
	}

	// oMainKey := []string{}
	// oMainValue := []string{}
	oStrKeyKV := []string{" 1=1 "}
	for k, v := range oKeyKV {
		// oMainKey = append(oMainKey, k)
		// oMainValue = append(oMainValue, v)
			oStrKeyKV = append(oStrKeyKV, `"` + Filter(k).(string) + `" = '` + Filter(v).(string) + `' `)
	}

	oFields := []string{}
	oValues := []string{}
	oFieldValue := []string{}
	iParamIndex := 1
	for k, v := range oMapKV {
		sFK := Filter(k).(string)
		sFV := Filter(v).(string)
		oFields = append(oFields, sFK)
		oValues = append(oValues, sFV)

		oFieldValue = append(oFieldValue, ` "`+ sFK + `" = $`+ strconv.Itoa(iParamIndex) +` `)

		iParamIndex++
	}
	//just for safe
	sTBName = Filter(sTBName).(string)

	sUpSql := ` UPDATE ` + sTBName + ` SET ` + strings.Join(oFieldValue, ` , `) + `  where ` + strings.Join(oStrKeyKV, ` and `)
	//Info(sUpSql)

	stmt, err := oPgSql.obj.Prepare(sUpSql)
	defer stmt.Close()

	if err != nil {
		Error(err.Error())
		return false
	}

	//	stmt.Exec("test", "test")
	oInValues := make([]interface{}, len(oValues))
	for i, s  := range oValues {
		oInValues[i] =  s
	}
	//fmt.Println(oInValues)
	_, errDo := stmt.Exec(oInValues...)
	if errDo != nil {
		Error(errDo.Error())
		return false
	}
	//LastInsertId
	return true
}

const(
	PREDEL_DATA_BAK_Folder = "$owns/pgsql-bak"
)
//return operate status
func (oPgSql *PGSql) Delete(sTBName string, oKeyKV []PGQueryWhere) bool{

	if "" == sTBName {
		Error("sorry, can not find the db object for delete the table in pqsql database or the sTBName is not empty ! ")
		return false
	}
	sTBName = Filter(sTBName).(string)

	oTBColumns := oPgSql.GetColumns(sTBName)
	if len(oTBColumns) == 0 {
			return false
	}
	//when delete just bak data
	oMapQuery := map[string] PGQueryObj {
		sTBName : PGQueryObj{
				Fields : oTBColumns,
				ObjWhere : oKeyKV,
		},
	}
	oPreDelDatas := oPgSql.QueryData(oMapQuery)
	if len(oPreDelDatas) == 0 {
		Error("sorry, can not find the data from table named " + sTBName + " with the condition you want ! ")
		return false
	}
	//write to file in folder $owns/pgsql-bak
	oPreDelJD, err := json.Marshal(oPreDelDatas)
	if nil != err {
		Error(err.Error())
		return false
	}
	//write data to file from memory
	//sSaveFilePath := strings.Trim(PREDEL_DATA_BAK_Folder, "/") + "/" + sTBName + "/" + sTBName + ".bak"
	_, err = os.Stat(PREDEL_DATA_BAK_Folder)
	if nil != err {
		Error(err.Error())
		err = os.MkdirAll(PREDEL_DATA_BAK_Folder, 0777)
		if nil != err {
				Error(err.Error())
				return false
		}
	}
	sSaveFolderPath := filepath.Join(PREDEL_DATA_BAK_Folder, sTBName)
	_, err = os.Stat(sSaveFolderPath)
	if nil != err {
		Error(err.Error())
		err = os.MkdirAll(sSaveFolderPath, 0777)
		if nil != err {
			Error(err.Error())
			return false
		}
	}
	//save file path
	sSaveFilePath := filepath.Join(sSaveFolderPath, time.Now().Format("2006-01-02") + ".bak")
	//firlst open it
	oFile, err := os.OpenFile(sSaveFilePath, os.O_RDWR|os.O_APPEND, 0644)
	defer oFile.Close()
	if nil != err {
		Error(err.Error() + "   => the file is not exist pre and create now ... ")
		oFile, err = os.Create(sSaveFilePath)
		defer oFile.Close()
		if nil != err {
			Error(err.Error())
			return false
		}
	}
	_, err = oFile.WriteString(string(oPreDelJD) + "\n")
	if nil != err {
		Error(err.Error())
		return false
	}
	// oMainKey := []string{}
	// oMainValue := []string{}
	sDeSql := ` DELETE FROM ` + sTBName + fnPGWhereHandle(oKeyKV)

	//now do delete
	stmt, err := oPgSql.obj.Prepare(sDeSql)
	defer stmt.Close()

	if err != nil {
		Error(err.Error())
		return false
	}
	//	stmt.Exec("test", "test")
	_, errDo := stmt.Exec()
	//_.RowsAffected()
	if errDo != nil {
		Error(errDo.Error())
		return false
	}
	//LastInsertId
	//if delete success save the info
	//delete is import things
	Info("the data is true delete from " + sTBName + " but bak in $owns/pgsql-bak")
	return true
}


//query data
//the where struct
type PGQueryWhere struct {
	Field string
	Op string
	Val string
}
// fields is a array param that have you will show field
// orders is a array param that you will order
// the objwhere is a struct array for get where condition
// PageInfo is a struct for page info
type PGQueryObj struct {
	Fields []string
	Orders []string
	ObjWhere []PGQueryWhere
	PageInfo struct {
			star int
			size int
	}
}

//think function of where condition
func fnPGWhereHandle(oWhereObj []PGQueryWhere) string {
	//return the sql clolumn
	sWhere := ""
	if len(oWhereObj) > 0 {
		oWhereAnds := []string{}
		bWhereJoin := false
		for _, oWKV := range oWhereObj {
			sWKV := ""
			sWK := Filter(oWKV.Field).(string)
			sWV := Filter(oWKV.Val).(string)

			switch oWKV.Op  {
				case "equal" : sWKV = ` "`+ sWK +`" = '`+ sWV +`'`
				case "lowerEqual" : sWKV = ` lower("`+ sWK +`") = '`+ sWV +`'`
				case "contains" : sWKV = ` "`+ sWK +`" like '%`+ sWV +`%'`
				case "containsTop" : sWKV = ` "`+ sWK +`" like '`+ sWV +`%'`
				case "containsLast" : sWKV = ` "`+ sWK +`" like '%`+ sWV +`'`
			}
			if !reflect.DeepEqual(oWKV, PGQueryWhere{}) {
				if !bWhereJoin {
						bWhereJoin = true
				}
				// /Info(sWKV)
				oWhereAnds = append(oWhereAnds, sWKV)
			}
		}

		if len(oWhereAnds) > 0 && bWhereJoin {
			sWhereAnd := strings.Join(oWhereAnds, " and ")
			if "" != sWhereAnd{
				sWhere = " where " + sWhereAnd
			}
		}
	}
	//final need return condition of the obj
	return sWhere
}

//look up all field with the format array
func (oPgSql *PGSql) QueryData(oTBQueryObj map[string]PGQueryObj) map[string][]map[string]string {

	//think first about param
	if nil == oTBQueryObj {
		Error("sorry, the query object is can not be empty !")
		return nil
	}

	//this obj will return
	oReObj := make(map[string][]map[string]string)

	for k, v := range oTBQueryObj {

		if 0 == len(v.Fields) {
			continue
		}
		//new a oFields to save the value

		//think where
		sWhere := fnPGWhereHandle(v.ObjWhere)
		//think order
		sOrder := ""
		// bIsOrder := true
		// for _, bv := range v.Fields {
		// 	for fk, fv := range v.Orders {
		// 		if fv != bv && bIsOrder{
		// 				bIsOrder = false
		// 		}
		// 		v.Orders[fk] = Filter(fv).(string)
		// 	}
		// }
		//&& bIsOrder
		if len(v.Orders) > 0 {
			sOrder = ` order by "` + strings.Join(v.Orders, `","`) + `" `
		}
		//Info(`select "`+ strings.Join(v.Fields, `","`) +`" from ` + k + sWhere + sOrder)
 		rows, err := oPgSql.obj.Query(`select "`+ strings.Join(v.Fields, `","`) +`" from ` + k + sWhere + sOrder)
		//Info(`select "`+ strings.Join(v.Fields, `","`) +`" from ` + k + sWhere + sOrder)
 		if nil != err {
 			Error(err.Error())
 			return nil
 		}

		oInFields := make([]interface{}, len(v.Fields))
		for i, v  := range v.Fields {
			oNewV := v
			oInFields[i] = &oNewV
		}
 		var oGet []map[string]string
 		//rows.Columns()
 		//fmt.Println(oInFields)
	 	for rows.Next() {
 		//fmt.Println(rows)
 			errScan := rows.Scan(oInFields...)
 			if nil != errScan {
 				Error(errScan.Error())
 				return nil
 			}
 			oRowMap := make(map[string]string)
 			for im, is := range v.Fields {
 				// if nil == oInFields[im] {
 				// 	oRowMap[is] = ""
 				// 	continue
 				// }
 				oRowMap[is] = *(oInFields[im].(*string))
 			}
 			oGet = append(oGet, oRowMap)
 		}
		oReObj[k] = oGet
	}
	return oReObj
}


//the fn to exsql
func (oPgSql *PGSql) ExPGSql(sSql string) bool {

	//think first
	if "" == sSql {
		Error("sorry, the param of sSql you send can not be empty ! ")
		return false
	}
	stmt, err := oPgSql.obj.Prepare(Filter(sSql).(string))
	defer stmt.Close()
	if nil != err {
		Error(err.Error())
		return false
	}
	_, errDo := stmt.Exec()
	if nil != errDo {
		Error(errDo.Error())
		return false
	}
	return true
}
