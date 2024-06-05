package dbstore

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `json:"session_id"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type SessionStore struct {
	db *gorm.DB
}

func NewSessionStore(DB *gorm.DB) *SessionStore {
	return &SessionStore{
		db: DB,
	}
}

func (s *SessionStore) CreateSession(session *Session) (*Session, error) {

	session.SessionID = uuid.New().String()

	result := s.db.Create(session)

	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}

func (s *SessionStore) GetUserFromSession(sessionID string, userID string) (*User, error) {
	var session Session

	err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email")
	}).Where("session_id = ? AND user_id = ?", sessionID, userID).First(&session).Error

	if err != nil {
		return nil, err
	}

	if session.User.ID == 0 {
		return nil, fmt.Errorf("no user associated with the session")
	}

	return &session.User, nil
}
