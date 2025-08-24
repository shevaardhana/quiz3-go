package routers

import (
	"quiz3-go/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api", controllers.BasicAuthMiddleware)
	{
		//cek login basic auth
		api.GET("/profile", controllers.Profile)

		// category
		category := api.Group("/categories")
		{
			category.GET("", controllers.GetCategories)
			category.POST("", controllers.CreateCategory)
			category.GET("/:id", controllers.GetCategory)
			category.DELETE("/:id", controllers.DeleteCategory)
			category.GET("/:id/books", controllers.GetBooksByCategory)
		}

		book := api.Group("/books")
		{
			book.GET("", controllers.GetBooks)
			book.POST("", controllers.CreateBook)
			book.GET("/:id", controllers.GetBook)
			book.DELETE("/:id", controllers.DeleteBook)
			book.PUT("/:id", controllers.UpdateBook)
		}
	}

	return r
}
