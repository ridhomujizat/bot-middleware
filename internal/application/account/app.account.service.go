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

func (a *AccountService) GetUserByName(name string) (*AccountSetting, error) {
	var account AccountSetting
	if err := a.db.Where("account = ?", name).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
