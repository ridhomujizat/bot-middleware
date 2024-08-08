// myapp/user/AccountService.go
package appAccount

import (
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
		return nil, err
	}
	return &account, nil
}

func (a *AccountService) SaveAccount(account *AccountSetting) error {
	return a.db.Save(account).Error
}
