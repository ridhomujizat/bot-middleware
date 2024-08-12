// myapp/user/AccountService.go
package appAccount

import (
	"bot-middleware/internal/pkg/util"

	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (a *AccountService) GetAccount(acc string, tenant string) (*AccountSetting, error) {
	var account AccountSetting
	if err := a.db.Where("account = ?", acc).Where("tenantId = ?", tenant).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, util.HandleAppError(err, "GetAccount", "First", true)
	}
	return &account, nil
}

func (a *AccountService) SaveAccount(account *AccountSetting) error {
	if err := a.db.Save(account).Error; err != nil {
		return util.HandleAppError(err, "SaveAccount", "Save", true)
	}

	return nil
}
