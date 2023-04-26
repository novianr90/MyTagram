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
		userRouter.POST("/register", userController.Register)
		userRouter.POST("/login", userController.Login)
	}

	photoRouter := app.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())

		photoRouter.POST("", photoController.CreatePhoto)

		photoRouter.GET("", photoController.GetAll)
		photoRouter.GET("/:photoId", middlewares.PhotoAuthorization(&photoService), photoController.GetPhotoById)

		photoRouter.PUT("/:photoId", middlewares.PhotoAuthorization(&photoService), photoController.UpdatePhotoById)

		photoRouter.DELETE("/:photoId", middlewares.PhotoAuthorization(&photoService), photoController.DeletePhotoById)

		// Comment
		photoRouter.POST("/:photoId/comment", middlewares.PhotoAuthorization(&photoService), commentController.CreateComment)
	}

	filesRouter := app.Group("/files")
	{
		filesRouter.GET("/:image", fileController.GetImages)

	}

	commentRouter := app.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.GET("", commentController.GetAllComments)

		commentRouter.GET("/:commentId", middlewares.CommentAuth(&commentService), commentController.GetComment)

		commentRouter.PUT("/:commentId", middlewares.CommentAuth(&commentService), commentController.UpdateComment)
		commentRouter.DELETE("/:commentId", middlewares.CommentAuth(&commentService), commentController.DeleteComment)
	}

	socialMediaRouter := app.Group("/accounts")
	{
		socialMediaRouter.Use(middlewares.Authentication())

		socialMediaRouter.POST("/", socialMediaController.CreateSocialMedia)
	}

	return app
}
