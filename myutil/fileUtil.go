package myFileUtil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

// GetFileName :获取url文件名字
func GetFileName(fileUrl string) string {
	u, err := url.Parse(fileUrl)
	if err != nil {
		panic(err)
	}

	filename := path.Base(u.Path)
	return filename
}

// FileExists 判断是否存在，存在为ture，不存在为false
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("文件不存在")
			return false
		} else {
			fmt.Println("发生错误：", err)
		}
	} else {
		fmt.Println("文件存在")
		return true
	}
	return false
}

// DownloadFile 文件下载
// url：文件下载地址
// path：保存地址
func DownloadFile(url string, path string) error {

	fileName := GetFileName(url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP请求错误：", err)
		return errors.New("HTTP请求错误")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP返回码不正确：", resp.StatusCode)
		return errors.New("HTTP返回码不正确")
	}

	file, err := os.Create(path + fileName)
	if err != nil {
		fmt.Println("创建文件失败：", err)
		return errors.New("创建文件失败")
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("拷贝数据失败：", err)
		return errors.New("拷贝数据失败")
	}
	fmt.Println("下载完成：", fileName)
	return nil

}
