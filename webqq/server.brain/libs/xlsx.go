package libs

import(

  //extends
  xlsx "github.com/bakerstreet-club/xlsx-go"
  . "bakerstreet-club/logs"

)

func SaveXlsx(sXlsxName, sSheetName string, oTitle []string, oData [][]string) {

      //check first
      if "" == sXlsxName {
          Error("sorry, the xlsx name is can not be empty")
          return
      }
      if "" == sSheetName {
          Error("sorry, the sheet name is can not be empty")
          return
      }
      if 0 == len(oTitle){
          Error("sorry, the title array of length is can not equal zero")
          return
      }

      var oFile *xlsx.File
      var oSheet *xlsx.Sheet
      var oRow *xlsx.Row
      var oCell *xlsx.Cell
      var oErr error

      oFile = xlsx.NewFile()
      oSheet, _ = oFile.AddSheet(sSheetName)

      oRow = oSheet.AddRow()
      for _, v := range oTitle {
        oCell = oRow.AddCell()
        oCell.Value = v
      }

      if(0 == len(oData)){
        oRow = oSheet.AddRow()
        oCell = oRow.AddCell()
        oCell.Value = "（暂无数据！）"
        oCell.Merge(len(oTitle), 1)
      }else{
        for _, v := range oData {
          oRow = oSheet.AddRow()
          for _, vv := range v {
            oCell = oRow.AddCell()
            oCell.Value = vv
          }
        }
      }

      //save xlsx
      oErr = oFile.Save(sXlsxName)
      if oErr != nil {
      	Error(oErr.Error())
      }
}
