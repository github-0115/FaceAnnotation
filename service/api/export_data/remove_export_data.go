package export_data

import (
	exportmodel "FaceAnnotation/service/model/exportmodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func RemoveExportData(c *gin.Context) {
	fileName := c.PostForm("filename")
	if strings.EqualFold(fileName, "") {
		log.Error(fmt.Sprintf("fileName = nil"))
		c.JSON(400, gin.H{
			"code":    vars.ErrExportParmars.Code,
			"message": vars.ErrExportParmars.Msg,
		})
		return
	}

	err := exportmodel.RemoveExportFile(fileName)
	if err != nil {
		log.Error(fmt.Sprintf("remove export data err %s", err))
		c.JSON(400, gin.H{
			"code":    vars.ErrExportNotFound.Code,
			"message": vars.ErrExportNotFound.Msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "expoet data remove success",
	})
}
