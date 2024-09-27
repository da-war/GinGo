package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("templates/**/**")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "views/index.html", gin.H{
			"title": "Main website",
		})
	})

	log.Println("Server started on port 8080")
	r.Run(":8080")
}
