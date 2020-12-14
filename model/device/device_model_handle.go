package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/utils/log"
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

func findComChargeTaskRecord(entity *DeviceCom) (bool, error) {
	err := ExecSQL().Where("device_id = ? AND charge_id = ? AND com_id = ?", entity.DeviceID, entity.ChargeID, entity.ComID).Order("create_time desc").First(&entity).Error
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

func DeviceComChargeTaskOps(entity *DeviceCom, isAck bool) error {
	taskRecord := &DeviceCom{
		DeviceID: entity.DeviceID,
		ChargeID: entity.ChargeID,
		ComID:    entity.ComID,
	}

	hasRecord, err := findComChargeTaskRecord(taskRecord)
	if err != nil {
		return err
	}

	//存在记录
	if hasRecord {
		if isAck {
			entity.State = (taskRecord.State | COM_CHARGE_START_ACK_BIT)
		}
		updateParam := map[string]interface{}{"max_energy": entity.MaxEnergy, "max_time": entity.MaxTime, "max_electricity": entity.MaxElectricity, "state": entity.State}
		if err := ExecSQL().Model(taskRecord).Updates(updateParam).Error; err != nil {
			SystemLog("update Data Error:", zap.Any("SQL", taskRecord), zap.Error(err))
		}
	} else { //不存在记录
		if isAck {
			entity.State |= COM_CHARGE_START_ACK_BIT
		}
		if err := ExecSQL().Create(entity).Error; err != nil {
			structTpey := reflect.Indirect(reflect.ValueOf(entity)).Type()
			SystemLog("Create ", structTpey, " Error ", zap.Any("SQL", entity), zap.Error(err))
		}
	}
	return nil
}