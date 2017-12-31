package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo"
)

// 接入微信接口服务
func api(c echo.Context) error {
	input := c.Param("url")
	if input == "" {
		PrintErrorHandler(c.Response().Writer, c.Request())

	}
	uDec, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		PrintErrorHandler(c.Response().Writer, c.Request())
	} else {
		PrintHandler(string(uDec), c.Response().Writer, c.Request())
	}
	var err2 error
	return err2
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SaveImg 保存图片到本地
func SaveImg(imageURL, saveName string) (n int64, err error) {
	out, err := os.Create(saveName)
	defer out.Close()
	if err != nil {
		return
	}
	resp, err := http.Get(imageURL)

	if err != nil {
		return
	}
	pix, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	if err != nil {
		return
	}
	n, err = io.Copy(out, bytes.NewReader(pix))

	if err != nil {
		return
	}
	// todo 获取图片类型
	// fmt.Println(resp.Header.Get("Content-Type"))
	return
}

func PrintErrorHandler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "images/404.png")
}

func PrintHandler(u string, w http.ResponseWriter, r *http.Request) {

	imgname := GetMd5String(u)

	imgpath := fmt.Sprintf("file/%v.jpg", imgname)

	// 如果本地服务器不存在缓存，再去拿
	_, err := os.Stat(imgpath)
	if os.IsNotExist(err) {
		_, err2 := SaveImg(u, imgpath)
		if err2 != nil {
			imgpath = "images/404.png"
		} else {
			src, err := imaging.Open(imgpath)
			if err != nil {
				fmt.Println("Open failed: %v", err)
				imgpath = "images/404.png"
			} else {
				// src = imaging.Resize(src, 256, 0, imaging.Lanczos)
				src = imaging.Resize(src, 400, 0, imaging.Lanczos)
				src = imaging.CropAnchor(src, 400, 300, imaging.Center)
				err = imaging.Save(src, imgpath)
				if err != nil {
					fmt.Println("Save failed: %v", err)
					imgpath = "images/404.png"
				}
			}
		}
	}
	http.ServeFile(w, r, imgpath)
}

func main() {
	e := echo.New()

	// Handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "pic crop save server")
	})

	e.File("/favicon.ico", "images/favicon.ico")

	e.GET("/:url", api)

	// Handler
	// e.GET("/:url/:param", func(c echo.Context) error {
	// 	input := c.Param("url")

	// 	// input = "http://mmbiz.qpic.cn/mmbiz_jpg/Z8SUoc8pJqdBfxCtd51ibGNr7IOXNI4DuUVbpToIqdhZUibOYDmW0S8nCGchoExiaMIPJ8oaMsXB7KSyKNcsVjibBg/0?wx_fmt=jpeg"
	// 	// uEnc := base64.URLEncoding.EncodeToString([]byte(input))
	// 	// aHR0cDovL21tYml6LnFwaWMuY24vbW1iaXpfanBnL1o4U1VvYzhwSnFkQmZ4Q3RkNTFpYkdOcjdJT1hOSTREdVVWYnBUb0lxZGhaVWliT1lEbVcwUzhuQ0djaG9FeGlhTUlQSjhvYU1zWEI3S1N5S05jc1ZqaWJCZy8wP3d4X2ZtdD1qcGVn

	// 	// fmt.Println(string(uEnc))

	// 	uDec, err := base64.URLEncoding.DecodeString(input)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// 	fmt.Println(string(uDec))
	// 	// fmt.Println(string(uEnc))
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	// Start server
	e.Logger.Fatal(e.Start(":8003"))
}
