package main

import (
	"embed"
	"goweb/auth"
	"goweb/db"
	"goweb/logger"
	"goweb/user"
	"goweb/version"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var log = logger.New("main", false)

func Migrate() {
	user.AutoMigrate(db.DB)
}

//go:embed assets/*
var staticFS embed.FS

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))

	// Add static assets to binary.
	// See https://github.com/gin-gonic/examples/tree/master/assets-in-binary
	// templ := template.Must(template.New("").ParseFS(staticFS, "templates/index.html"))
	// router.SetHTMLTemplate(templ)
	// router.StaticFS("/public", http.FS(staticFS))

	// Uncoment this for faster html and css development
	router.LoadHTMLGlob("templates/**.html")
	router.Static("/assets", "./assets")

	router.GET("/", index)
	router.GET("/ping", ping)
	router.POST("/login", auth.Login)
	router.GET("/logout", auth.Logout)

	private := router.Group("/private")
	private.Use(auth.AuthRequired)
	{
		private.GET("/me", auth.Me)
		private.GET("/status", auth.Status)
	}

	return router
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Dunder Mifflin",
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"foo": "bar"})
}

func main() {
	_, err := db.GetDB("postgres://postgres:secret@localhost:5432/goweb")
	if err != nil {
		panic(err)
	}
	Migrate()

	log.Infof("Starting the app version %s", version.Version)
	r := setupRouter()
	r.Run()
}
