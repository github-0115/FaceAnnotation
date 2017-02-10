package thr_result

import (
	imageend "FaceAnnotation/service/api/image"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	thrfacemodel "FaceAnnotation/service/model/thrfacemodel"
	usermodel "FaceAnnotation/service/model/usermodel"
	vars "FaceAnnotation/service/vars"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
)

func GetFaceResult(c *gin.Context) {
	name, _ := c.Get("username")
	username := name.(string)
	url := c.Query("url")
	pointType, err := strconv.Atoi(c.Query("point_type"))

	if url == "" || err != nil {
		log.Error(fmt.Sprintf("parmar nil err%v"))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	_, err = usermodel.QueryUser(username)
	if err != nil {
		log.Error(fmt.Sprintf("find user error:%s", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrLoginParams.Code,
			"message": vars.ErrLoginParams.Msg,
		})
		return
	}

	urlStrs := strings.Split(url, "/")
	md5 := urlStrs[len(urlStrs)-1]
	imageColl, err := imagemodel.QueryImage(md5)
	if err != nil {
		log.Error(fmt.Sprintf("query image err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	imageBytes, err := getImg(url)
	if err != nil {
		log.Error(fmt.Sprintf("query image err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	imageColl, err = faceRes(md5, imageBytes, imageColl)
	if err != nil {
		log.Error(fmt.Sprintf("query image err", err.Error()))
		c.JSON(400, gin.H{
			"code":    vars.ErrImageModelNotFound.Code,
			"message": vars.ErrImageModelNotFound.Msg,
		})
		return
	}

	pointRep := imageend.SwitchPoint(int64(pointType), imageColl)

	if pointRep == nil {
		pointRep = &imageend.PointsRep{}
	}
	log.Error(fmt.Sprintf("pointRep:=%s", pointRep))
	c.JSON(200, gin.H{
		"code":   0,
		"points": pointRep,
	})
}

func faceRes(photoName string, fileByte []byte, imageColl *imagemodel.ImageModel) (*imagemodel.ImageModel, error) {
	//face++ res
	//	var thrRes *thrfacemodel.FaceModelV3
	thrRes, err := thrfacemodel.ThrFaceFileResV3(photoName, fileByte)
	if err != nil {
		log.Error(fmt.Sprintf("get face++ res fail err:%s", err))
	}
	res1B, _ := json.Marshal(thrRes)

	var result interface{}
	if err := json.Unmarshal(res1B, &result); err != nil {
		fmt.Println("json unmarshal err=%s", err)
	}
	fmt.Println("-----result%s-----", result)
	if imageColl.ThrFaces == nil {
		imageColl.ThrFaces["face++"] = make(map[string]interface{})
	}
	imageColl.ThrFaces["face++"] = make(map[string]interface{})
	imageColl.ThrFaces["face++"]["83"] = result

	_, err = imagemodel.UpsertImageModel(imageColl)
	if err != nil {
		log.Error(fmt.Sprintf("image update err", err.Error()))

		return nil, err
	}
	fmt.Println("-----imageColl%s-----", imageColl)
	return imageColl, nil
}

func getImg(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Error(fmt.Sprintf("get url=%s pic err=%s", url, err))
		return nil, err
	}
	defer resp.Body.Close()

	pix, err := ioutil.ReadAll(resp.Body)
	return pix, err
}
