package routes

import (
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, UserController controller.UserController, jwtService services.JWTService, BookController controller.BookController) {
	routes := route.Group("/api/user")
	{
		routes.POST("", UserController.RegisterUser)
		routes.GET("", middleware.Authenticate(jwtService), UserController.GetAllUser)
		routes.POST("/login", UserController.LoginUser)
		routes.DELETE("/", middleware.Authenticate(jwtService), UserController.DeleteUser)
		routes.PUT("/", middleware.Authenticate(jwtService), UserController.UpdateUser)
		routes.GET("/me", middleware.Authenticate(jwtService), UserController.MeUser)
		routes.GET("/top", BookController.GetTopBooks)
		routes.GET("/books", BookController.GetAllBooks)
		routes.GET("/books/:book_id", BookController.GetBookPages)
	}

	medias := route.Group("/api/media")
	{
		medias.GET("/get/storage/:path/:dirname/:filename", BookController.GetImage)
	}

}
