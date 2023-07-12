package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading the env file")
	}

	//gin set up
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")
	r.MaxMultipartMemory = 32 << 20

	//setup s3
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Printf("error:%v", err)
		return
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "upload images",
		})
	})

	r.POST("/", func(c *gin.Context) {

		file, err := c.FormFile("image")
		if err != nil {
			c.HTML(http.StatusBadRequest, "index.html", gin.H{
				"error": "failed to upload the image" + err.Error(),
			})
			return
		}

		f, openErr := file.Open()
		if openErr != nil {
			c.HTML(http.StatusBadRequest, "index.html", gin.H{
				"error": "failed to open the image" + openErr.Error(),
			})
			return
		}

		result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("jerseyhub"),
			Key:    aws.String(file.Filename),
			Body:   f,
			ACL:    "public-read",
		})
		if uploadErr != nil {
			c.HTML(http.StatusBadRequest, "index.html", gin.H{
				"error": "failed to upload the image" + uploadErr.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"image": result.Location,
		})
	})

	r.Run()
}
