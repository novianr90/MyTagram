package routers

import (
	"final-project-hacktiv8/controllers"
	"final-project-hacktiv8/middlewares"
	"final-project-hacktiv8/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "final-project-hacktiv8/docs"
)

// @title MyGram API
// @version 1.0
// @description This API serve on upload photo and social media
// @termsOfService http://swagger.io/terms
// @contact.name Novian Rachman Iskandar
// @contact.email novianr123@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host mytagram-production.up.railway.app
// @BasePath /
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

		commentService = services.CommentService{
			DB: db,
		}

		socialMediaService = services.SocialMediaService{
			DB: db,
		}

		photoController = controllers.PhotoController{
			PhotoService: &photoService,
			FileService:  &fileService,
		}

		fileController = controllers.FileController{
			FileService: &fileService,
		}

		commentController = controllers.CommentController{
			CommentService: &commentService,
		}

		socialMediaController = controllers.SocialMediaController{
			Service: &socialMediaService,
		}
	)

	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	userRouter := app.Group("/users")
	{
		// Create
		userRouter.POST("/register", userController.Register)
		// Create
		userRouter.POST("/login", userController.Login)
	}

	photoRouter := app.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())

		// Create
		photoRouter.POST("", photoController.CreatePhoto)

		// Read All
		photoRouter.GET("", photoController.GetAll)
		// Read
		photoRouter.GET("/:photoId", middlewares.PhotoAuthorization(&photoService), photoController.GetPhotoById)

		// Update
		photoRouter.PUT("/:photoId", middlewares.PhotoAuthorization(&photoService), photoController.UpdatePhotoById)

		// Delete
		photoRouter.DELETE("/:photoId", middlewares.PhotoAuthorization(&photoService), photoController.DeletePhotoById)

		// Create
		photoRouter.POST("/:photoId/comment", middlewares.PhotoAuthorization(&photoService), commentController.CreateComment)
	}

	filesRouter := app.Group("/files")
	{
		filesRouter.GET("/:image", fileController.GetImages)

	}

	commentRouter := app.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		// Read All
		commentRouter.GET("", commentController.GetAllComments)

		// Read
		commentRouter.GET("/:commentId", middlewares.CommentAuth(&commentService), commentController.GetComment)

		// Update
		commentRouter.PUT("/:commentId", middlewares.CommentAuth(&commentService), commentController.UpdateComment)
		// Delete
		commentRouter.DELETE("/:commentId", middlewares.CommentAuth(&commentService), commentController.DeleteComment)
	}

	socialMediaRouter := app.Group("/accounts")
	{
		socialMediaRouter.Use(middlewares.Authentication())

		// Create
		socialMediaRouter.POST("", socialMediaController.CreateSocialMedia)

		// Read All
		socialMediaRouter.GET("", socialMediaController.GetAllAccounts)

		// Read
		socialMediaRouter.GET("/:accountId", middlewares.AccountAuthorization(&socialMediaService), socialMediaController.GetAccountById)

		// Update
		socialMediaRouter.PUT("/:accountId", middlewares.AccountAuthorization(&socialMediaService), socialMediaController.UpdateAccount)

		// Delete
		socialMediaRouter.DELETE("/:accountId", middlewares.AccountAuthorization(&socialMediaService), socialMediaController.DeleteAccount)
	}

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return app
}
