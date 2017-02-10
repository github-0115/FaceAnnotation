package main

import (
	exportendpoint "FaceAnnotation/service/api/export_data"
	imageendpoint "FaceAnnotation/service/api/image"
	imagetaskendpoint "FaceAnnotation/service/api/image_task"
	dataendpoint "FaceAnnotation/service/api/import_data"
	loginendpoint "FaceAnnotation/service/api/login"
	smalltaskendpoint "FaceAnnotation/service/api/small_task"
	taskendpoint "FaceAnnotation/service/api/task"
	thrResultendpoint "FaceAnnotation/service/api/thr_result"
	userendpoint "FaceAnnotation/service/api/user"
	middleware "FaceAnnotation/service/middleware"
	"flag"
	"fmt"
	"time"

	"github.com/gin-gonic/contrib/gzip"
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
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,POST,PUT,DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Access-Control-Allow-Headers, LoginToken,X-Requested-With, X-CSRF-Token",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	r.Static("origin_images", "./origin_images")
	r.Static("get_export_data", "./exportData")
	r.POST("import_data", dataendpoint.ImportData)
	authorized := r.Group("/user")
	{
		authorized.POST("login", loginendpoint.Login)
	}

	imagegroup := r.Group("/image")
	imagegroup.Use(middleware.AuthToken)
	{
		imagegroup.GET("get_one_image", imageendpoint.GetImage)
		imagegroup.POST("save_image", imageendpoint.SaveImageRes)
		imagegroup.POST("import_image", dataendpoint.ImportImage)
		imagegroup.POST("import_images", dataendpoint.ImportImages)
		imagegroup.DELETE("delete_image", imageendpoint.DeleteImage)
		imagegroup.POST("import_res", dataendpoint.ImportResult)
		imagegroup.POST("export_data", exportendpoint.ExportData)
		imagegroup.PUT("remove_data", exportendpoint.RemoveExportData)
	}

	taskgroup := r.Group("/task")
	taskgroup.Use(middleware.AuthToken)
	{
		taskgroup.GET("create_image_task", imagetaskendpoint.CreateImageTask)
		taskgroup.POST("image_create_task", taskendpoint.ImageCreateTask)
		taskgroup.GET("all_images", imagetaskendpoint.GetAllImages)
		taskgroup.PUT("start_task", taskendpoint.StartTask)
		taskgroup.PUT("stop_task", taskendpoint.StopTask)
		taskgroup.GET("image_task_list", imagetaskendpoint.ImageTaskList)
		taskgroup.GET("get_small_tasks", smalltaskendpoint.GetSmallTasks)
		taskgroup.GET("task_all_images", taskendpoint.GetTaskAllImages)
		taskgroup.GET("small_task_all_images", smalltaskendpoint.GetSmallTaskAllImages)
		taskgroup.PUT("remove_image_task", imagetaskendpoint.RemoveImageTask)
		taskgroup.PUT("remove_task", taskendpoint.RemoveTask)
		taskgroup.POST("create_small_task", smalltaskendpoint.CreateSmallTask)
		taskgroup.GET("small_task_list", smalltaskendpoint.SmallTaskList)
		taskgroup.GET("small_image_list", smalltaskendpoint.GetSmallTaskImages)
		taskgroup.POST("create_task", taskendpoint.CreateTask)
		taskgroup.GET("task_list", taskendpoint.TaskList)
	}

	thrgroup := r.Group("/thr")
	thrgroup.Use(middleware.AuthToken)
	{
		thrgroup.GET("face_res", thrResultendpoint.GetFaceResult)
	}

	if *debug == false {
		r.POST("add_user", userendpoint.AddUser)
	}

	p := fmt.Sprintf(":%d", *port)
	log.Info("listen port " + p)
	r.Run(p)
}
