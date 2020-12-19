package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/utils/log"
	. "GoServer/utils/string"
	. "GoServer/utils/time"
	"go.uber.org/zap"
	"reflect"
)

func CreateDeviceAndModuleInfo(entity *CreateDeviceInfo) error {
	var err error
	tx := ExecSQL().Begin()
	switch entity.Type {
	case NO_DEVICE_WITH_MODULE:	//建立 模组 和 设备信息
		entity.Device.ID, err = TXCreateSQLAndRetLastID(tx,&entity.Device)
		if err != nil {
			SystemLog("add DeviceInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		loopInsertComList := StringJoin([]interface{}{"call InsertDeviceComList(",entity.Device.ID,",9)"})
		if err = tx.Exec(loopInsertComList).Error; err != nil {
			SystemLog(loopInsertComList," Error")
			tx.Rollback()
			return err
		}
		entity.Module.DeviceID = entity.Device.ID
		entity.Module.ID, err = TXCreateSQLAndRetLastID(tx,&entity.Module)
		if err != nil {
			SystemLog("add ModuleInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}
	case DEVICE_BUILD_BIT: //需单独创建Module
		entity.Module.DeviceID = entity.Device.ID
		entity.Module.UID = entity.Device.UID
		entity.Module.ID, err = TXCreateSQLAndRetLastID(tx,&entity.Module)
		if err != nil {
			SystemLog("add ModuleInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}
	case MODULE_BUILD_BIT: //需单独创建Device
		entity.Device.ID, err = TXCreateSQLAndRetLastID(tx,&entity.Device)
		if err != nil {
			SystemLog("add DeviceInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		loopInsertComList := StringJoin([]interface{}{"call InsertDeviceComList(",entity.Device.ID,",9)"})
		if err = tx.Exec(loopInsertComList).Error; err != nil {
			SystemLog(loopInsertComList," Error")
			tx.Rollback()
			return err
		}
		entity.Module.DeviceID = entity.Device.ID
		updateModule := StringJoin([]interface{}{"UPDATE `module_info` SET `uid`=0, `device_id`=",entity.Device.ID," WHERE `id`=",entity.Module.ID})
		if err = tx.Exec(updateModule).Error; err != nil {
			SystemLog(loopInsertComList," Error")
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	entity.Type = HAS_DEVICE_WITH_MODULE
	return nil
}

func findComChargeTaskRecord(entity *DeviceCharge) (bool, error) {
	err := ExecSQL().Where("device_id = ? AND token = ? AND com_id = ?", entity.DeviceID, entity.Token, entity.ComID).Order("create_time desc").First(&entity).Error
	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			SystemLog("createComChargeTaskRecord select Error", zap.Error(err))
			return true, err
		}
	}
	return hasRecord, nil
}

func isChaegeEnding(state uint32) (ret bool) {
	ret = true
	switch state {
	case COM_CHARGE_START_BIT, COM_CHARGE_START_ACK_BIT:
		ret = false
	}
	return ret
}

func DeviceComChargeTaskOps(entity *DeviceCharge, state uint32) error {
	taskRecord := &DeviceCharge{
		DeviceID: entity.DeviceID,
		Token:    entity.Token,
		ComID:    entity.ComID,
	}

	hasRecord, err := findComChargeTaskRecord(taskRecord)
	if err != nil {
		return err
	}

	//存在记录
	if hasRecord {
		if (taskRecord.State&state) == state && state != COM_CHARGE_RUNING_BIT {
			return nil
		}

		entity.State = (taskRecord.State | state)
		updateParam := make(map[string]interface{})

		if state == COM_CHARGE_STOP_BIT {
			updateParam["state"] = entity.State
			updateParam["end_time"] = GetTimestampMs()
		} else {
			updateParam["max_energy"] = entity.MaxEnergy
			updateParam["max_time"] = entity.MaxTime
			updateParam["max_electricity"] = entity.MaxElectricity
			updateParam["state"] = entity.State
			updateParam["use_energy"] = entity.UseEnergy
			updateParam["use_time"] = entity.UseTime
			updateParam["max_charge_electricity"] = entity.MaxChargeElectricity
			updateParam["average_power"] = entity.AveragePower
			updateParam["max_power"] = entity.MaxPower
			if state == COM_CHARGE_RUNING_BIT {
				updateParam["update_time"] = GetTimestampMs()
			} else if isChaegeEnding(state) {
				updateParam["end_time"] = GetTimestampMs()
			}
		}

		if err := ExecSQL().Model(taskRecord).Updates(updateParam).Error; err != nil {
			SystemLog("update Data Error:", zap.Any("SQL", taskRecord), zap.Error(err))
		}
	} else { //不存在记录
		entity.State = state
		if state == COM_CHARGE_RUNING_BIT {
			entity.UpdateTime = GetTimestampMs()
		} else if isChaegeEnding(state) {
			entity.EndTime = GetTimestampMs()
		}
		if err := ExecSQL().Create(entity).Error; err != nil {
			structTpey := reflect.Indirect(reflect.ValueOf(entity)).Type()
			SystemLog("Create ", structTpey, " Error ", zap.Any("SQL", entity), zap.Error(err))
		}
	}
	return nil
}
