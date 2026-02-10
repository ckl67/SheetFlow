package controllers

import (
	"net/http"
	//"path/filepath"

	//"backend/"
	"backend/api/middlewares"

	"github.com/gin-gonic/gin"
)

func (server *Server) SetupRouter() {
	r := gin.New()
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Version info
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": server.Version})
	})

	// API root
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	api := r.Group("/api")

	// Public routes
	api.POST("/login", server.Login)
	api.POST("/request_password_reset", server.RequestPasswordReset)

	// Secure routes
	secure := api.Group("")
	secure.Use(middlewares.AuthMiddleware())

	// Users
	api.POST("/reset_password", server.ResetPassword)
	secure.GET("/users", server.GetUsers)
	secure.GET("/users/:id", server.GetUser)
	secure.POST("/users", server.CreateUser)
	secure.PUT("/users/:id", server.UpdateUser)
	secure.DELETE("/users/:id", server.DeleteUser)

	// Sheets
	secure.GET("/sheets", server.GetSheetsPage)
	secure.POST("/sheets", server.GetSheetsPage)
	secure.GET("/sheet/:sheetName", server.GetSheet)
	secure.PUT("/sheet/:sheetName", server.UpdateSheet)
	secure.DELETE("/sheet/:sheetName", server.DeleteSheet)
	secure.POST("/upload", server.UploadFile)
	secure.PUT("/sheet/:sheetName/info", server.UpdateSheetInformationText)
	secure.POST("/sheet/:sheetName/info", server.UpdateSheetInformationText)

	// Thumbnails & PDFs
	api.GET("/sheet/thumbnail/:name", server.GetThumbnail)
	secure.GET("/sheet/pdf/:composer/:sheetName", server.GetPDF)

	// Search
	secure.GET("/search/:searchValue", server.SearchSheets)
	secure.GET("/search/composers/:searchValue", server.SearchComposers)

	// Tags
	secure.POST("/tag/sheet/:sheetName", server.AppendTag)
	secure.DELETE("/tag/sheet/:sheetName", server.DeleteTag)
	secure.GET("/tag", server.FindSheetsByTag)
	secure.POST("/tag", server.FindSheetsByTag)

	// Composers
	secure.GET("/composers", server.GetComposersPage)
	secure.POST("/composers", server.GetComposersPage)
	secure.PUT("/composer/:composerName", server.UpdateComposer)
	secure.DELETE("/composer/:composerName", server.DeleteComposer)
	api.GET("/composer/portrait/:composerName", server.ServePortraits)

	server.Router = r
}
