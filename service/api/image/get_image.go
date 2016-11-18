package image

import (
	facemodel "FaceAnnotation/service/model/facemodel"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	vars "FaceAnnotation/service/vars"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetImage(c *gin.Context) {

	image_list, err := imagemodel.GetImageList("./origin_images")
	if err != nil {
		log.Error(fmt.Sprintf("get local image list err %S", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrJsonUnmarshal.Code,
			"message": vars.ErrJsonUnmarshal.Msg,
		})
		return
	}

	already_list, err := facemodel.QueryAll()
	if err != nil {
		log.Error(fmt.Sprintf("get already image list err %S", err.Error()))
		if err != facemodel.ErrFaceModelNotFound {
			c.JSON(400, gin.H{
				"code":    vars.ErrFaceCursor.Code,
				"message": vars.ErrFaceCursor.Msg,
			})
			return
		}
	}

	not_list := RemoveDuplicatesAndEmpty(image_list, already_list)

	c.JSON(200, gin.H{
		"code": 0,
		"url":  not_list,
	})
}

func RemoveDuplicatesAndEmpty(a []string, b []*facemodel.FaceUrl) (ret []string) {

	a_len := len(a)
	b_len := len(b)
	for i := 0; i < a_len; i++ {
		for j := 0; j < b_len; j++ {
			if strings.EqualFold(a[i], b[j].Url) {
				continue
			}
			ret = append(ret, a[i])
		}
	}

	return
}
