package user

import (
	"errors"
	"goweb/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User database model - `users` table in postgres
type User struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password;not null"`
	CreatedAt    time.Time
}

// Migrate the schema of database if needed
func AutoMigrate(db db.Database) error {
	return db.AutoMigrate(&User{})
}

func (u *User) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func AddUser(db db.Database, user *User, password string) error {
	user.setPassword(password)
	user.CreatedAt = time.Now()
	return db.Save(user).Error
}

func FindOneUser(db db.Database, user *User) (User, error) {
	var u User
	result := db.Where(user).First(&u)
	return u, result.Error
}

// // You could update properties of an UserModel to database returning with error info.
// //  err := db.Model(userModel).Update(UserModel{Username: "wangzitian0"}).Error
// func (u *User) Update(db db.Database, user *User) error {
// 	db := common.GetDB()
// 	err := db.Model(u).Update(user).Error
// 	return err
// }
