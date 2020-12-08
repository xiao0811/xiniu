package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/handle"
)

// UploadImages 图片上传
func UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "图片格式不正确", c)
		return
	}
	files := form.File["upload[]"]
	var images []string
	for _, file := range files {
		s := strings.Split(file.Filename, ".")
		name := time.Now().Format("20060102150405") + strconv.Itoa(handle.RandInt(1000, 9999)) + "." + s[len(s)-1]
		// 上传文件到指定的路径
		if err := c.SaveUploadedFile(file, "/var/www/html/"+time.Now().Format("20060102")+"/"+name); err != nil {
			log.Println(err)
			continue
		}
		images = append(images, time.Now().Format("20060102")+"/"+name)
	}
	handle.ReturnSuccess("ok", images, c)
}

// UploadImage 上传单个图片
func UploadImage(c *gin.Context) {
	// 单文件
	file, err := c.FormFile("file")
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "图片格式不正确", c)
		return
	}
	s := strings.Split(file.Filename, ".")
	name := time.Now().Format("20060102150405") + strconv.Itoa(handle.RandInt(1000, 9999)) + "." + s[len(s)-1]
	// 上传文件到指定的路径
	if err := c.SaveUploadedFile(file, "/var/www/html/"+time.Now().Format("20060102")+"/"+name); err != nil {
		handle.ReturnError(http.StatusBadRequest, "图片上传失败", c)
		return
	}
	image := time.Now().Format("20060102") + name
	handle.ReturnSuccess("ok", image, c)
}

// ShowImage 获取图片
func ShowImage(c *gin.Context) {
	imageName := "/var/www/html" + c.Param("images_name")
	c.File(imageName)
}
