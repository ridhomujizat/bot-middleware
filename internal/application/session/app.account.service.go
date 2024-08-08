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

func (a *SessionService) FindSession(uniqueId, channelPlatform, channelSource, tenantId string) (*Session, error) {
	var session Session
	if err := a.db.Where("unique_id = ? AND channel_platform = ? AND channel_source = ? AND tenantId = ?", uniqueId, channelPlatform, channelSource, tenantId).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
