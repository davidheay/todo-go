package dbstore

import (
	"todo-go/internal/util/hash"

	"gorm.io/gorm"
)

// User model info
// @Description User  information
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
	Todos    []Todo `json:"-"`
	Roles    []Role `gorm:"many2many:user_roles;"`
}
type UserStore struct {
	db           *gorm.DB
	passwordhash hash.PasswordHash
}

func NewUserStore(DB *gorm.DB, PasswordHash hash.PasswordHash) *UserStore {
	return &UserStore{
		DB,
		PasswordHash,
	}
}

func (s *UserStore) CreateUser(email string, password string) error {

	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	return s.db.Create(&User{
		Email:    email,
		Password: hashedPassword,
	}).Error
}

func (s *UserStore) GetUser(email string) (*User, error) {

	var user User
	err := s.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}
func (s *UserStore) GetUserById(userId string) (*User, error) {
	var user User
	err := s.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}
