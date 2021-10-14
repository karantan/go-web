package auth

import (
	"encoding/gob"
	"net/http"

	"goweb/db"
	"goweb/user"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const UserKey = "user"

func init() {
	// register user.User struct so that we can save it to the session
	gob.Register(&user.User{})
}

// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

func LoginForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)
	if user != nil {
		c.Redirect(http.StatusFound, "/backend/")
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Dunder Mifflin",
	})
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	u, _ := user.FindOneUser(db.DB, &user.User{Username: username})
	err := u.CheckPassword(password)
	// Check for username and password match
	if u.Username == "" || err != nil {
		session.AddFlash("Authentication failed.", "Warn")
		c.HTML(http.StatusForbidden, "login.html", gin.H{
			"MsgWarn": session.Flashes("Warn"),
		})
		return
	}

	// Save the username in the session
	session.Set(UserKey, u) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		session.AddFlash("Failed to save session.", "Warn")
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"MsgWarn": session.Flashes("Warn"),
		})
		return
	}
	c.Redirect(http.StatusFound, "/backend/")
	return
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(UserKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}
