package main

import (
	"goweb/db"
	"goweb/logger"
	"goweb/user"
	"goweb/version"

	"github.com/spf13/cobra"
)

var log = logger.New("goweb-admin", true)

var adminCmd = &cobra.Command{
	Use:   "goweb-admin",
	Short: "...",
	Long:  `...`,
}

func init() {
	adminCmd.Run = userManagement
}

func userManagement(cmd *cobra.Command, args []string) {
	log.Infof("Running goweb-admin version %s", version.Version)
	DB, err := db.GetDB("postgres://postgres:secret@localhost:5432/goweb")
	if err != nil {
		panic(err)
	}
	user.AddUser(DB, &user.User{Username: "admin"}, "secret")
}

func main() {
	adminCmd.Execute()
}
