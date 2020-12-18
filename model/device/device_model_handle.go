package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/utils/log"
	. "GoServer/utils/time"
	"go.uber.org/zap"
	"reflect"
)

func CreateDevInfo(entity *CreateDeviceInfo) error {
	device := entity.Device
	module := entity.Module
	log := entity.Log

	err := ExecSQL().Where("device_sn = ?", device.DeviceSn).First(&device).Error
	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			SystemLog("transactionCreateDevInfo select Error", zap.Error(err))
			return err
		}
	}

	if hasRecord {
		//设备已存在，单独建立模组信息
		tx := ExecSQL().Begin()

		module.DeviceID = device.ID
		if err := tx.Create(&module).Error; err != nil {
			SystemLog("add ModuleInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		var id []uint64
		if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
			SystemLog("get LastID Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		var ModuleID uint64 = id[0]
		log.ModuleID = ModuleID
		if err := tx.Create(&log).Error; err != nil {
			SystemLog("add module connect log Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		//事务建立 模组 和 设备信息
		var id []uint64
		tx := ExecSQL().Begin()
		if err := tx.Create(&device).Error; err != nil {
			SystemLog("add DeviceInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
			SystemLog("get LastID Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		var DeviceID uint64 = id[0]
		device.ID = DeviceID
		module.DeviceID = DeviceID
		if err := tx.Create(&module).Error; err != nil {
			SystemLog("add ModuleInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
			SystemLog("get LastID Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		var ModuleID uint64 = id[0]
		log.ModuleID = ModuleID
		if err := tx.Create(&log).Error; err != nil {
			SystemLog("add module connect log Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	Redis().UpdateDeviceIDToRedisByDeviceSN(device.DeviceSn, device.ID)
	return nil
}

func findComChargeTaskRecord(entity *DeviceCharge) (bool, error) {
	err := ExecSQL().Debug().Where("device_id = ? AND token = ? AND com_id = ?", entity.DeviceID, entity.Token, entity.ComID).Order("create_time desc").First(&entity).Error
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
		if (taskRecord.State & state) == state {
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
			if isChaegeEnding(state) {
				updateParam["end_time"] = GetTimestampMs()
			}
		}
		if err := ExecSQL().Model(taskRecord).Updates(updateParam).Error; err != nil {
			SystemLog("update Data Error:", zap.Any("SQL", taskRecord), zap.Error(err))
		}
	} else { //不存在记录
		entity.State = state
		if isChaegeEnding(state) {
			entity.EndTime = GetTimestampMs()
		}
		if err := ExecSQL().Create(entity).Error; err != nil {
			structTpey := reflect.Indirect(reflect.ValueOf(entity)).Type()
			SystemLog("Create ", structTpey, " Error ", zap.Any("SQL", entity), zap.Error(err))
		}
	}
	return nil
}
