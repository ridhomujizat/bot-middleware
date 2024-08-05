// myapp/user/AccountService.go
package appSession

import (
	"gorm.io/gorm"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (a *SessionService) GetUserByName(name string) (*Session, error) {
	var session Session
	if err := a.db.Where("account = ?", name).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
