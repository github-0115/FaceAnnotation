package exportmodel

import (
	imagemodel "FaceAnnotation/service/model/imagemodel"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/inconshreveable/log15"
)

var (
	dirPath         = "./exportfile/"
	ErrDirNotFound  = errors.New("dir ads not found")
	ErrFileNotFound = errors.New("file ads not found")
	ErrCreateDir    = errors.New("create dir err")
	ErrCreateFile   = errors.New("create file err")
	ErrReadFile     = errors.New("read file err")
	ErrWriteFile    = errors.New("write file err")
)

type Rep struct {
	Url string  `json:"url"`
	Res string  `json:"res"`
	P   float64 `json:"p"`
	S   float64 `json:"s"`
	N   float64 `json:"n"`
}

func SaveResFile(filename string, res *imagemodel.ImageModel) error {
	fmt.Println(filename)
	isExist, _ := PathExists(dirPath)

	if !isExist {
		err := os.Mkdir(dirPath, 0777)
		if err != nil {
			log.Error(fmt.Sprintf("create dir err%v", err))
			return ErrCreateDir
		}
	}

	file, err := os.Create(dirPath + filename + ".txt")
	if err != nil {
		log.Error(fmt.Sprintf("create file err%v", err))
		return ErrCreateFile
	}
	defer file.Close()

	resbyte, _ := json.Marshal(res)

	w := bufio.NewWriter(file)
	n4, err := w.Write(resbyte)
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()

	return nil
}

func SaveExportFile(filename string, imageUrl []Rep) error {
	isExist, _ := PathExists(dirPath)

	if !isExist {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			log.Error(fmt.Sprintf("create dir err%v", err))
			return ErrCreateDir
		}
	}

	file, err := os.Create(dirPath + filename)
	if err != nil {
		log.Error(fmt.Sprintf("create file err%v", err))

		return ErrCreateFile
	}

	_, err = file.WriteString(repToString(imageUrl))
	if err != nil {
		log.Error(fmt.Sprintf("write file err%v", err))
		os.Remove(dirPath + filename)
		return ErrWriteFile
	}

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

func repToString(s []Rep) string {

	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {
		out, err := json.Marshal(s[i])
		if err != nil {
			panic(err)
		}
		if i == len(s)-1 {
			buffer.WriteString(string(out))

		} else {
			buffer.WriteString(string(out) + "\n")
		}

	}

	return buffer.String()
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
