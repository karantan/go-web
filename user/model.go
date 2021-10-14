package user

import (
	"errors"
	"goweb/db"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User database model - `users` table in postgres
type User struct {
	gorm.Model          // include fields ID, CreatedAt, UpdatedAt, DeletedAt
	Username     string `gorm:"column:username;unique"`
	PasswordHash string `gorm:"column:password;not null"`
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

// FindByID retrievs a user by ID
func FindByID(db db.Database, ID int) (User, error) {
	var u User
	result := db.Find(&u, ID)
	return u, result.Error
}

func GetAll(db db.Database) ([]User, error) {
	var users []User
	result := db.Find(&users)
	return users, result.Error
}

func Update(db db.Database, user *User) error {
	return db.Save(&user).Error
}

// DeleteByID soft deletes the user `user`.
func DeleteByID(db db.Database, ID int) error {
	var u User
	return db.Delete(&u, ID).Error
}
