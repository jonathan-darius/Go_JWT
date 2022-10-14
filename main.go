package main

import (
	"Latihan_Mongo/controller"
	docs "Latihan_Mongo/docs"
	"Latihan_Mongo/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

// @BasePath /v1
var (
	server       *gin.Engine
	us           service.UserService
	uc           controller.UserController
	ctx          context.Context
	userc        *mongo.Collection
	mongoclient  *mongo.Client
	err          error
	currentImage *imageupload.Image
)

func init() {
	ctx = context.TODO()

	//mongoconn := options.Client().ApplyURI("mongodb://admin:jaringan123@192.168.29.86:27017/")
	mongoconn := options.Client().ApplyURI("mongodb://admin:jaringan123@192.168.29.86:27017/")
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	fmt.Println("mongo connection established")

	userc = mongoclient.Database("userdb").Collection("users")
	us = service.NewUserService(userc, ctx)
	uc = controller.New(us)
	server = gin.Default()
}

func main() {
	defer mongoclient.Disconnect(ctx)
	docs.SwaggerInfo.BasePath = "/v1"
	basepath := server.Group("/v1")
	uc.RegisterUserRoutes(basepath)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r := *server

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.GET("/image", func(c *gin.Context) {
		if currentImage == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		currentImage.Write(c.Writer)
	})

	r.GET("/thumbnail", func(c *gin.Context) {
		if currentImage == nil {
			c.AbortWithStatus(http.StatusNotFound)
		}

		t, err := imageupload.ThumbnailJPEG(currentImage, 300, 300, 80)

		if err != nil {
			panic(err)
		}

		t.Write(c.Writer)
	})

	r.POST("/upload", func(c *gin.Context) {
		img, err := imageupload.Process(c.Request, "file")
		if err != nil {
			panic(err)
		}

		currentImage = img

		c.Redirect(http.StatusMovedPermanently, "/")
	})

	log.Fatal(server.Run(":8080"))

}
