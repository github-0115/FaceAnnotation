package exportmodel

import (
	cfg "FaceAnnotation/config"
	imagemodel "FaceAnnotation/service/model/imagemodel"
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/inconshreveable/log15"
)

var (
	dirPath         = "./exportData/"
	ErrDirNotFound  = errors.New("dir ads not found")
	ErrFileNotFound = errors.New("file ads not found")
	ErrCreateDir    = errors.New("create dir err")
	ErrCreateFile   = errors.New("create file err")
	ErrReadFile     = errors.New("read file err")
	ErrWriteFile    = errors.New("write file err")
	domain          = cfg.Cfg.Domian
)

type Rep struct {
	Url string  `json:"url"`
	Res string  `json:"res"`
	P   float64 `json:"p"`
	S   float64 `json:"s"`
	N   float64 `json:"n"`
}

type Res struct {
	Name   string              `json:"name"`
	Points []*imagemodel.Point `json:"url"`
}

var (
	imagesDomain = "http://faceannotation.oss-cn-hangzhou.aliyuncs.com/"
)

func ImageDataZip(images []*imagemodel.ImageModel, ress []*Res, resName string) (string, error) {
	// 创建一个缓冲区用来保存压缩文件内容
	buf := new(bytes.Buffer)

	// 创建一个压缩文档
	w := zip.NewWriter(buf)

	// 将文件加入压缩文档
	for _, image := range images {
		file, err := w.Create("img/" + image.Url)
		if err != nil {
			log.Error(fmt.Sprintf("create file err%v", err))
			return "", ErrCreateFile
		}

		resbyte, _ := getUrlImg(imagesDomain + image.Md5)
		_, err = file.Write(resbyte)
		if err != nil {
			log.Error(fmt.Sprintf("Write file err%v", err))
		}
	}

	file, err := w.Create("points.txt")
	if err != nil {
		log.Error(fmt.Sprintf("create file err%v", err))
		return "", ErrCreateFile
	}
	for _, res := range ress {
		resbyte, _ := json.Marshal(res)

		_, err = file.Write([]byte(res.Name + "\t" + string(resbyte) + "\n"))
		if err != nil {
			log.Error(fmt.Sprintf("Write file err%v", err))
		}
	}

	// 关闭压缩文档
	err = w.Close()
	if err != nil {
		log.Error(fmt.Sprintf("close file err%v", err))
	}

	// 将压缩文档内容写入文件
	f, err := os.OpenFile(dirPath+resName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(fmt.Sprintf("open file err%v", err))
	}
	buf.WriteTo(f)

	return domain + "get_export_data/" + resName, nil
}

func SaveImage(md5 string, filename string) error {

	isExist, _ := PathExists(dirPath + "images/")

	if !isExist {
		err := os.Mkdir(dirPath+"images/", 0777)
		if err != nil {
			log.Error(fmt.Sprintf("create dir err%v", err))
			return ErrCreateDir
		}
	}

	file, err := os.Create(dirPath + "images/" + filename)
	if err != nil {
		log.Error(fmt.Sprintf("create file err%v", err))
		return ErrCreateFile
	}
	defer file.Close()

	resbyte, _ := getUrlImg(imagesDomain + md5)

	w := bufio.NewWriter(file)
	n4, err := w.Write(resbyte)
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()

	return nil
}

func SaveImageRes(filename string, res []*imagemodel.Point) error {

	isExist, _ := PathExists(dirPath)

	if !isExist {
		err := os.Mkdir(dirPath, 0777)
		if err != nil {
			log.Error(fmt.Sprintf("create dir err%v", err))
			return ErrCreateDir
		}
	}
	var seekn int64 = 0
	file, err := os.OpenFile(dirPath+"points.txt", os.O_WRONLY, 0644)
	if err != nil {
		file, err = os.Create(dirPath + "points.txt")
		if err != nil {
			log.Error(fmt.Sprintf("create file err%v", err))
			return ErrCreateFile
		}

	} else {
		// 查找文件末尾的偏移量
		seekn, _ = file.Seek(0, os.SEEK_END)
		fmt.Println(seekn)
	}
	defer file.Close()

	res1B, _ := json.Marshal(res)
	// 从末尾的偏移量开始写入内容
	_, err = file.WriteAt([]byte(filename+"\t"+string(res1B)+"\n"), seekn)

	return nil
}

func GetAllExportFile() ([]string, error) {
	isExist, _ := PathExists(dirPath)

	if !isExist {
		log.Error(fmt.Sprintf(" dir not Exists err"))
		return nil, ErrDirNotFound
	}

	files, err := getFilelist(dirPath)
	if err != nil {
		log.Error(fmt.Sprintf("get filelist err%v", err))

		return nil, ErrFileNotFound
	}

	return files, nil
}

func ReadExportFile(filename string) (string, error) {

	isExist, _ := PathExists(dirPath + filename)

	if !isExist {
		log.Error(fmt.Sprintf(" file not Exists err"))
		return "", ErrFileNotFound
	}

	data, err := ioutil.ReadFile(dirPath + filename)
	if err != nil {
		log.Error(fmt.Sprintf("read file err%v", err))

		return "", ErrFileNotFound
	}

	return string(data), nil
}

func RemoveExportFile(filename string) error {

	isExist, _ := PathExists(dirPath + filename)

	if !isExist {
		log.Error(fmt.Sprintf(" file not Exists err"))
		return ErrFileNotFound
	}
	err := os.Remove(dirPath + filename)
	if err != nil {
		log.Error(fmt.Sprintf("Remove file err%v", err))

		return ErrFileNotFound
	}

	return nil
}

func ReadLocalFile(filePath string) (string, error) {

	isExist, _ := PathExists(filePath)

	if !isExist {
		log.Error(fmt.Sprintf(" file not Exists err"))
		return "", ErrFileNotFound
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(fmt.Sprintf("read file err%v", err))

		return "", ErrFileNotFound
	}

	return string(data), nil
}

func getFilelist(path string) ([]string, error) {
	fileList := make([]string, 0, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)

		files := strings.Split(path, "/")
		var (
			name string
		)
		if len(files) > 1 {
			name = files[len(files)-1]
		}
		fileList = append(fileList, name)

		return nil
	})

	if err != nil {
		log.Error(fmt.Sprintf("get file list %v\n", err))
		return nil, ErrFileNotFound
	}

	return fileList, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func repToString(filename string, res []*imagemodel.Point) string {

	var buffer bytes.Buffer
	for i := 0; i < len(res); i++ {
		out, err := json.Marshal(res[i])
		if err != nil {
			panic(err)
		}
		if i == len(res)-1 {
			buffer.WriteString(string(out) + "\n")
		}

	}

	return filename + ":" + buffer.String()
}

func toString(s []string) string {

	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {

		if i == len(s)-1 {
			buffer.WriteString(s[i])
		} else {
			buffer.WriteString(s[i] + "\n")
		}

	}

	return buffer.String()
}

func getUrlImg(url string) (pix []byte, err error) {
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	fmt.Println(name)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err = ioutil.ReadAll(resp.Body)

	return

}
