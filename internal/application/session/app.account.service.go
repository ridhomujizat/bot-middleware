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

func (a *SessionService) FindSession(
	unique_id string,
	channel_platform string,
	channel_source string,
	tenantId string,
) (*Session, error) {
	var session Session
	if err := a.db.Where("unique_id = ? AND channel_platform = ? AND channel_source = ? AND tenant_id = ?", unique_id, channel_platform, channel_source, tenantId).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
