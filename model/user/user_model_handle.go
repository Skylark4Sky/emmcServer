package user

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/utils/log"
	"go.uber.org/zap"
)

func CreateNewUser(entity *CreateUserInfo, hasAuth bool) error {
	var id []uint64
	tx := ExecSQL().Begin()
	if err := tx.Create(&entity.Base).Error; err != nil {
		SystemLog("add UserBase Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
		SystemLog("get LastID Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	var userID uint64 = id[0]

	if hasAuth == true {
		entity.Auth.UID = userID
		if err := tx.Create(&entity.Auth).Error; err != nil {
			SystemLog("add UserAuth Error", zap.Error(err))
			tx.Rollback()
			return err
		}
	}

	entity.Log.UID = userID
	if err := tx.Create(&entity.Log).Error; err != nil {
		SystemLog("add UserRegisterLog Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	entity.Extra.UID = userID
	if err := tx.Create(&entity.Extra).Error; err != nil {
		SystemLog("add UserExtra Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	entity.Location.UID = userID
	if err := tx.Create(&entity.Location).Error; err != nil {
		SystemLog("add UserLocation Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

