package routers

import (
	"final-project-hacktiv8/controllers"
	"final-project-hacktiv8/middlewares"
	"final-project-hacktiv8/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartServer(db *gorm.DB) *gin.Engine {
	var (
		app = gin.Default()

		userService = services.UserService{
			DB: db,
		}

		userController = controllers.UserController{
			UserService: &userService,
		}

		photoService = services.PhotoService{
			DB: db,
		}

		fileService = services.FileService{
			DB: db,
		}

		photoController = controllers.PhotoController{
			PhotoService: &photoService,
			FileService:  &fileService,
		}

		fileController = controllers.FileController{
			FileService: &fileService,
		}
	)

	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	userRouter := app.Group("/users")
	{
		userRouter.POST("/register", userController.Register)
		userRouter.POST("/login", userController.Login)
	}

	photoRouter := app.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())

		photoRouter.POST("", photoController.CreatePhoto)

		photoRouter.GET("", photoController.GetAll)
		photoRouter.GET("/:photoId", middlewares.PhotoAuthorization(&photoService), photoController.GetPhotoById)

	}

	filesRouter := app.Group("/files")
	{
		filesRouter.GET("/:image", fileController.GetImages)

	}

	return app
}
