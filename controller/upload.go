package controller

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/handle"
)

// UploadImages 图片上传
func UploadImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	var images []string
	for _, file := range files {
		s := strings.Split(file.Filename, ".")
		name := time.Now().Format("20060102150405") + strconv.Itoa(handle.RandInt(1000, 9999)) + s[len(s)-1]
		// 上传文件到指定的路径
		if err := c.SaveUploadedFile(file, "./upload/images/"+name); err != nil {
			log.Println(err)
			continue
		}
		images = append(images, "./upload/images/"+name)
	}
	handle.ReturnSuccess("ok", images, c)
}

// ShowImage 获取图片
func ShowImage(c *gin.Context) {
	imageName := "./upload/images/" + c.Param("images_name")
	c.File(imageName)
}
