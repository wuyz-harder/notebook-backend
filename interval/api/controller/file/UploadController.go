package file

import (
	"GetHotWord/common"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-tika/tika"
	"github.com/google/uuid"
	"github.com/yanyiwu/gojieba"
)

// 这个服务需要用docker生成：docker run -d -p 9998:9998 apache/tika:latest
var client = tika.NewClient(nil, "http://127.0.0.1:9998")

// 接口，只要每种类型实现了里面的接口，那么该类型的每个对象都可以调用里面的方法
type MyInterface interface {
	getName()
	getAge()
}
type User struct {
	name string
	age  int
}

type get interface {
}

func (u User) getName() {
	fmt.Println(u.name)
}

func (u User) getAge() {
	fmt.Println(u.age)
}

// 去除字符串中的html标签
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func UploadFile(context *gin.Context) {

	var (
		jieba = gojieba.NewJieba()
	)

	uploadfile, err := context.FormFile("file")

	if err != nil {
		context.Error(common.NewError(400, 200400, "file参数不能为空"))
	}
	uid := uuid.NewString()
	saveErr := context.SaveUploadedFile(uploadfile, "../file"+string(os.PathSeparator)+uid+"-"+uploadfile.Filename)
	//判断是否保存失败
	if saveErr != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{
			"mes": saveErr.Error(),
		})
		return
	}
	ufile, _ := uploadfile.Open()
	ufile.Close()
	file, _ := os.OpenFile("../file"+string(os.PathSeparator)+uid+"-"+uploadfile.Filename, os.O_RDONLY, 0755)
	defer file.Close()
	var res []string
	//pdf文件用tiga
	if strings.Contains(uploadfile.Filename, ".pdf") {
		content, err := ReadPdf("file" + string(os.PathSeparator) + uid + "-" + uploadfile.Filename) // Read local pdf file
		content = TrimHtml(content)
		fmt.Println(content)
		if err != nil {
			panic(err)
		}
		res = jieba.Extract(content, 5)

	} else {
		all, err := io.ReadAll(file)
		if err != nil {
			return
		}
		res = jieba.Extract(string(all), 5)
	}
	//tiga
	context.JSON(
		http.StatusOK, gin.H{
			"msg":      "上传成功",
			"key_word": res,
			"url":      "file" + string(os.PathSeparator) + uid + "-" + uploadfile.Filename,
		})

}

func ReadPdf(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return "", err
	}
	//Background就是纯文本
	return client.Parse(context.Background(), f)
}

// 上传头像

func UploadAvatar(context *gin.Context) {

	uploadfile, err := context.FormFile("file")

	if err != nil {
		context.Error(common.NewError(400, 200400, "file参数不能为空"))
	}
	uid := uuid.NewString()
	saveErr := context.SaveUploadedFile(uploadfile, "../file/avatar"+string(os.PathSeparator)+uid+"-"+uploadfile.Filename)
	//判断是否保存失败
	if saveErr != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{
			"mes": saveErr.Error(),
		})
		return
	}
	//tiga
	context.JSON(
		http.StatusOK, gin.H{
			"msg": "上传成功",
			"url": "file/avatar" + string(os.PathSeparator) + uid + "-" + uploadfile.Filename,
		})

}
