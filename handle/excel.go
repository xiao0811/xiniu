package handle

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
)

// ExcelExport 导出 Excel
func ExcelExport(c *gin.Context, head []string, body [][]interface{}, filename string) {
	xlsx := excelize.NewFile()
	_ = xlsx.SetSheetRow("Sheet1", "A1", &head)
	for index, rowData := range body {
		_ = xlsx.SetSheetRow("Sheet1", "A"+strconv.Itoa(index+2), &rowData)
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	_ = xlsx.Write(c.Writer)
}
