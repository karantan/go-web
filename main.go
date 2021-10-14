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
	router.Use(sessions.Sessions("goweb", sessions.NewCookieStore([]byte("secret"))))

	// Add static assets to binary.
	// See https://github.com/gin-gonic/examples/tree/master/assets-in-binary
	// templ := template.Must(template.New("").ParseFS(staticFS, "templates/index.html"))
	// router.SetHTMLTemplate(templ)
	// router.StaticFS("/public", http.FS(staticFS))

	// Uncoment this for faster html and css development
	router.LoadHTMLGlob("templates/**/*.html")
	router.Static("/assets", "./assets")

	router.GET("/", auth.LoginForm)
	router.POST("/login", auth.Login)
	router.GET("/logout", auth.Logout)

	backend := router.Group("/backend")
	backend.Use(auth.AuthRequired)
	{
		backend.GET("/", dashboard)
		backend.GET("/users", user.List)
		backend.GET("/user/add", user.ShowAdd)
		backend.POST("/user/add", user.Add)
		backend.GET("/user/:id/edit", user.ShowEdit)
		backend.POST("/user/:id/edit", user.Edit)
		backend.GET("/user/:id/delete", user.Remove)
	}

	return router
}

func dashboard(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(auth.UserKey)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"user": user,
	})
	return
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
