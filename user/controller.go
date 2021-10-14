package user

import (
	"goweb/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const UserKey = "user"

func List(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)

	users, err := GetAll(db.DB)
	if err != nil {
		session.AddFlash(err.Error(), "Warn")
	}

	MsgWarn := session.Flashes("Warn")
	MsgInfo := session.Flashes("Info")
	session.Save()
	c.HTML(http.StatusOK, "list.html", gin.H{
		"user":    user,
		"MsgWarn": MsgWarn,
		"MsgInfo": MsgInfo,
		"users":   users,
	})
	return
}

func ShowAdd(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)
	c.HTML(http.StatusOK, "add.html", gin.H{
		"user": user,
	})
	return
}

func Add(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)

	username := c.PostForm("username")
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("password-confirm")
	if password != passwordConfirm {
		session.AddFlash("Passwords are not the same", "Warn")
		c.HTML(http.StatusOK, "add.html", gin.H{
			"user":    user,
			"MsgWarn": session.Flashes("Warn"),
		})
		return
	}
	err := AddUser(db.DB, &User{Username: username}, password)
	if err != nil {
		session.AddFlash(err.Error(), "Warn")
		c.HTML(http.StatusOK, "add.html", gin.H{
			"user":    user,
			"MsgWarn": session.Flashes("Warn"),
		})
		return
	}
	session.AddFlash("User Successfully Created", "Info")
	session.Save()
	c.Redirect(http.StatusFound, "/backend/users")
}

func ShowEdit(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)

	var editUser User
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		session.AddFlash("Invalid user ID.", "Warn")
	} else {
		editUser, err = FindByID(db.DB, userID)
		if err != nil {
			session.AddFlash(err.Error(), "Warn")
		}
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"user":     user,
		"MsgWarn":  session.Flashes("Warn"),
		"editUser": editUser,
	})
	return
}

func Edit(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)

	username := c.PostForm("username")
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("password-confirm")
	if password != passwordConfirm {
		session.AddFlash("Passwords are not the same", "Warn")
		c.HTML(http.StatusOK, "add.html", gin.H{
			"user":    user,
			"MsgWarn": session.Flashes("Warn"),
		})
		return
	}

	var editUser User
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		session.AddFlash("Invalid user ID.", "Warn")
	} else {
		editUser, err = FindByID(db.DB, userID)
		if err != nil {
			session.AddFlash(err.Error(), "Warn")
		}
		editUser.Username = username
		editUser.setPassword(password)
		err = Update(db.DB, &editUser)
		if err != nil {
			session.AddFlash(err.Error(), "Warn")
		}
	}

	session.AddFlash("User Successfully Updated", "Info")

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"user":     user,
		"MsgWarn":  session.Flashes("Warn"),
		"MsgInfo":  session.Flashes("Info"),
		"editUser": editUser,
	})
	return
}

func Remove(c *gin.Context) {
	session := sessions.Default(c)

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		session.AddFlash("Invalid user ID.", "Warn")
	} else {
		err = DeleteByID(db.DB, userID)
		if err != nil {
			session.AddFlash(err.Error(), "Warn")
		}
	}

	session.AddFlash("User Successfully Deleted", "Info")
	session.Save()
	c.Redirect(http.StatusFound, "/backend/users")
}
