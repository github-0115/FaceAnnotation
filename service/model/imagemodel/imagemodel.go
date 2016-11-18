package imagemodel

import (
	"os"
	"path/filepath"
	"strings"
)

func GetImageList(dirPth string) (files []string, err error) {
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {
			return err
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}

		urls := strings.Split(filename, "\\")
		files = append(files, urls[len(urls)-1])

		return nil
	})
	return files, err
}
