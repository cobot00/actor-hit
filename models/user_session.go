package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserSession struct {
	ID                   uint
	Turn                 uint
	Hit                  uint
	ExpiredAt            time.Time
	Answer               string
	AnswerCharacterNames string
	ResultSummary        string
}

func FindSession(db *gorm.DB, sessionId int) *UserSession {
	var session UserSession
	recordNotFound := db.Table("user_session").
		Where("id = ?", sessionId).Where("expired_at >= ?", time.Now()).
		First(&session).RecordNotFound()
	if recordNotFound {
		session = UserSession{
			Turn:      0,
			Hit:       0,
			ExpiredAt: time.Now().Add(time.Duration(1) * time.Hour),
		}
		db.Table("user_session").Create(&session)
	}
	return &session
}

func UpdateSession(db *gorm.DB, session *UserSession) {
	session.ExpiredAt = time.Now().Add(time.Duration(1) * time.Hour)
	db.Table("user_session").Save(session)
}
