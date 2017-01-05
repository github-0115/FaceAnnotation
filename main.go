package main

import (
	//	faceendpoint "FaceAnnotation/service/api/face"
	imageendpoint "FaceAnnotation/service/api/image"
	dataendpoint "FaceAnnotation/service/api/import_data"
	loginendpoint "FaceAnnotation/service/api/login"
	smalltaskendpoint "FaceAnnotation/service/api/small_task"
	taskendpoint "FaceAnnotation/service/api/task"
	userendpoint "FaceAnnotation/service/api/user"
	//	middleware "FaceAnnotation/service/middleware"
	//	thrfacemodelendpoint "FaceAnnotation/service/model/thrfacemodel"
	"flag"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/inconshreveable/log15"
	"github.com/itsjamie/gin-cors"
)

var (
	port  = flag.Int64("p", 80, "port")
	debug = flag.Bool("d", false, "debug model")
)

func ready() {

	flag.Parse()
	if *port == 8060 {
		gin.SetMode(gin.ReleaseMode)
		log.Info(fmt.Sprintf("ReleaseMode"))
	}
}

func main() {
	ready()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "POST",
		RequestHeaders:  "Origin, Authorization, Content-Type, Access-Control-Allow-Headers, LoginToken",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	r.Static("origin_images", "./origin_images")

	//	r.POST("thr", thrfacemodelendpoint.ThrFaceRes)
	//	r.POST("thr_et", thrfacemodelendpoint.EightThreeFace)
	//	r.POST("import_images", imageendpoint.ImportImages)
	r.POST("import_data", dataendpoint.ImportData)
	authorized := r.Group("/user")
	{
		authorized.POST("login", loginendpoint.Login)
	}

	imagegroup := r.Group("/image")
	//	imagegroup.Use(middleware.AuthToken)
	{
		//			imagegroup.GET("get_images", imageendpoint.GetImageList)
		imagegroup.GET("get_one_image", imageendpoint.GetImage)
		imagegroup.POST("save_image", imageendpoint.SaveImageRes)
		imagegroup.POST("import_image", dataendpoint.ImportImage)
	}

	taskgroup := r.Group("/task")
	//	taskgroup.Use(middleware.AuthToken)
	{
		//		taskgroup.GET("task_list", taskendpoint.GetTaskList)
		//		taskgroup.GET("get_one_task", taskendpoint.GetTask)
		taskgroup.POST("create_task", taskendpoint.CreateTask)
		taskgroup.POST("create_small_task", smalltaskendpoint.CreateSmallTask)
		taskgroup.GET("get_small_tasks", smalltaskendpoint.GetSmallTasks)
	}

	//	facegroup := r.Group("/face")
	//	facegroup.Use(middleware.AuthToken)
	//	{
	//		facegroup.POST("upsert", faceendpoint.UpsertFacePoint)
	//	}

	if *debug == false {
		r.POST("add_user", userendpoint.AddUser)
	}

	p := fmt.Sprintf(":%d", *port)
	log.Info("listen port " + p)
	r.Run(p)
}
