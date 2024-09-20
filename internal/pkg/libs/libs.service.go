package libs

import (
	appAccount "bot-middleware/internal/application/account"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/entities"
	"bot-middleware/internal/pkg/util"
	"bot-middleware/internal/webhook"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type LibsService struct {
	db *gorm.DB
}

func NewLibsService(db *gorm.DB) *LibsService {
	return &LibsService{db: db}
}

func (l *LibsService) GetAccount(acc, tenant, platform string) (*appAccount.AccountSetting, error) {
	var account appAccount.AccountSetting
	if err := l.db.Where("account = ?", acc).Where("tenant_id = ?", tenant).Where("account_platform = ?", platform).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, util.HandleAppError(err, "Lib GetAccount", "First", true)
	}
	return &account, nil
}

func (l *LibsService) Text(acc, tenant string, payload interface{}) (interface{}, error) {
	account, err := l.GetAccount(acc, tenant, string(entities.WhatsAppSocio))
	if err != nil {
		return nil, err
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		util.HandleAppError(err, "Lib Text", "JSON Marshal", true)
		return nil, err
	}
	res, StatusCode, errResponse := util.HttpPost(account.BaseURL, payloadData, map[string]string{
		"Content-Type": "application/json",
		"x-key":        account.Token,
	})
	if errResponse != nil {
		return nil, errResponse
	}

	var result interface{}
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		util.HandleAppError(err, "Lib Text", "JSON Unmarshal", true)
		return nil, err
	}

	if StatusCode != 200 {
		return nil, fmt.Errorf("error response: %s", res)
	}

	return result, nil
}

func (l *LibsService) Button(acc, tenant string, payload interface{}) (interface{}, error) {
	account, err := l.GetAccount(acc, tenant, string(entities.WhatsAppSocio))
	if err != nil {
		return nil, err
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, util.HandleAppError(err, "Lib Button", "JSON Marshal", true)
	}

	res, _, errResponse := util.HttpPost(account.BaseURL, payloadData, map[string]string{
		"Content-Type": "application/json",
		"x-key":        account.Token,
		"account_id":   acc,
	})
	if errResponse != nil {
		return nil, errResponse
	}

	var result interface{}
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return nil, util.HandleAppError(err, "Lib Button", "JSON Unmarshal", true)
	}

	return result, nil
}

func (l *LibsService) Carousel(acc, tenant string, payload interface{}) (interface{}, error) {
	account, err := l.GetAccount(acc, tenant, string(entities.WhatsAppSocio))
	if err != nil {
		return nil, err
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, util.HandleAppError(err, "Lib Carousel", "JSON Marshal", true)
	}

	res, _, errResponse := util.HttpPost(account.BaseURL, payloadData, map[string]string{
		"Content-Type": "application/json",
		"x-key":        account.Token,
		"account_id":   acc,
	})
	if errResponse != nil {
		return nil, errResponse
	}

	var result interface{}
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return nil, util.HandleAppError(err, "Lib Carousel", "JSON Unmarshal", true)
	}

	return result, nil
}

func (l *LibsService) FindSessionBySid(sid string) (*appSession.Session, error) {
	var session appSession.Session
	if err := l.db.Where("sid = ?", sid).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, util.HandleAppError(err, "Lib FindSessionBySid", "First", true)
	}
	return &session, nil
}

func (l *LibsService) FindSessionByUniqueId(uniqueId, tenantId string) (*appSession.Session, error) {
	var session appSession.Session
	if err := l.db.Where("unique_id = ? AND tenant_id = ?", uniqueId, tenantId).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, util.HandleAppError(err, "Lib FindSessionByUniqueId", "First", true)
	}
	return &session, nil
}

func (l *LibsService) CreateSession(session *appSession.Session) error {
	if err := l.db.Create(session).Error; err != nil {
		return util.HandleAppError(err, "Lib CreateSession", "Create", false)
	}

	sessionHistory := appSession.SessionHistory{
		Session:         *session,
		BotPlatform:     session.BotPlatform,
		ChannelSource:   session.ChannelSource,
		ChannelPlatform: session.ChannelPlatform,
		ChannelId:       session.ChannelId,
		Omnichannel:     session.Omnichannel,
	}

	if err := l.db.Create(&sessionHistory).Error; err != nil {
		return util.HandleAppError(err, "Lib CreateSession", "Create", false)
	}

	return nil
}

func (l *LibsService) UpdateSession(session *appSession.Session) error {
	if err := l.db.Model(&appSession.Session{}).Where("sid = ?", session.Sid).Updates(session).Error; err != nil {
		return util.HandleAppError(err, "Lib UpdateSession", "Updates", false)
	}

	sessionHistory := appSession.SessionHistory{
		Session:         *session,
		BotPlatform:     session.BotPlatform,
		ChannelSource:   session.ChannelSource,
		ChannelPlatform: session.ChannelPlatform,
		ChannelId:       session.ChannelId,
		Omnichannel:     session.Omnichannel,
	}

	if err := l.db.Create(&sessionHistory).Error; err != nil {
		return util.HandleAppError(err, "Lib UpdateSession", "Create", false)
	}

	return nil
}

func (l *LibsService) DeleteSession(payload *webhook.EndDTO) error {
	if err := l.db.Where("sid = ?", payload.Sid).Delete(&appSession.Session{}).Error; err != nil {
		return util.HandleAppError(err, "Lib DeleteSession", "Delete", true)
	}

	return nil
}

func (l *LibsService) DeleteSessionByUniqueId(payload *webhook.EndDTO) error {
	if err := l.db.Where("unique_id = ?", payload.UniqueId).Delete(&appSession.Session{}).Error; err != nil {
		return util.HandleAppError(err, "Lib DeleteSessionByUniqueId", "Delete", true)
	}

	return nil
}

func (l *LibsService) DeleteSessionOmnix(payload *webhook.EndDTO) error {
	payload.UniqueId = strings.ReplaceAll(payload.UniqueId, "+", "")
	if err := l.db.Where("unique_id = ? AND channel_account = ?", payload.UniqueId, payload.AccountId).Delete(&appSession.Session{}).Error; err != nil {
		return util.HandleAppError(err, "Lib DeleteSessionOmnix", "Delete", true)
	}

	return nil
}
