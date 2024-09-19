// myapp/user/AccountService.go
package appSession

import (
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/worker"
	"fmt"

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
	if err := a.db.Where("unique_id = ? AND channel_platform = ? AND channel_source = ? AND tenant_id = ?", uniqueId, channelPlatform, channelSource, tenantId).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (a *SessionService) InitAndCheckSession(uniqueId, channelPlatform, channelSource, tenantId string) (*worker.MetaData, error) {
	fmt.Println("InitAndCheckSession", uniqueId, channelPlatform, channelSource, tenantId)
	session, err := a.FindSession(uniqueId, channelPlatform, channelSource, tenantId)
	if err != nil {
		util.HandleAppError(err, "Whatsapp process", "FindSession", false)
		return nil, nil
	}

	result := worker.MetaData{}
	if session == nil {
		sid, err := util.GenerateId()
		if err != nil {
			util.HandleAppError(err, "Generate SID process", "GenerateId", false)
			return nil, nil
		}

		result.Sid = sid
		result.NewSession = true

	} else {
		result.Sid = session.Sid
		result.NewSession = false
	}

	return &result, nil
}
