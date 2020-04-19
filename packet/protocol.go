package Packect

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const GISUNLINK_RECV_TOPIC = "/device"
const GISUNLINK_SEND_TOPIC = "/device_state"

const DEVICE_INFO = "device_info"
const TRANSFER = "transfer"
const TRANSFER_RESULT = "transfer_result"
const FIRMWARE_UPDATE = "update_state"

const (
	//下发
	GISUNLINK_CHARGE_TASK      = 0x10 //下发
	GISUNLINK_DEVIDE_STATUS    = 0x11 //查询
	GISUNLINK_EXIT_CHARGE_TASK = 0x12 //终止
	GISUNLINK_SET_CONFIG       = 0x13 //配置
	GISUNLINK_RESTART          = 0x1E //重启
	//上报
	GISUNLINK_START_CHARGE     = 0x14 //开始
	GISUNLINK_CHARGEING        = 0x15 //执行中
	GISUNLINK_CHARGE_FINISH    = 0x16 //完成
	GISUNLINK_CHARGE_LEISURE   = 0x17 //空闲中
	GISUNLINK_CHARGE_BREAKDOWN = 0x18 //故障
	GISUNLINK_CHARGE_NO_LOAD   = 0x19 //noload
	GISUNLINK_UPDATE_FIRMWARE  = 0x1A //升级
	GISUNLINK_COM_UPDATE       = 0x1B //参数刷新
	GISUNLINK_STOP_CHARGE      = 0x1C //停止
	GISUNLINK_COM_NO_UPDATE    = 0x1D //参数没有刷新
)

type Protocol interface {
	Print() (retString string)
}

//更新状态
type UpdateState struct {
	Msg string `json:"msg"`
}

//传输回复
type TransferResult struct {
	ReqID   int64  `json:"req_id"`
	Success string `json:"success"`
	Msg     string `json:"msg"`
}

//信息截取
type DeviceInfo struct {
	Imei     string `json:"imei "`
	Version  string `json:"version"`
	DeviceSn string `json:"device_sn"`
	Sim      struct {
		ICCID string `json:"ICCID"`
		IMEI  string `json:"IMEI"`
	} `json:"sim"`
	CellInfo struct {
		Ci  int    `json:"ci "`
		Lac int    `json:"Lac"`
		Mnc int    `json:"Mnc"`
		Mcc int    `json:"Mcc"`
		Ta  int    `json:"ta"`
		Ext string `json:"Ext"`
	} `json:"cellInfo"`
}

type ComTaskStartTransfer struct {
	comID          uint8
	token          uint32
	maxEnergy      uint32
	maxElectricity uint32
	maxTime        uint32
}

type ComTaskStopTransfer struct {
	comID     uint8
	token     uint32
	forceStop uint8
}

type ComData struct {
	Id             uint8
	token          uint32
	maxEnergy      uint32
	useEnergy      uint32
	maxTime        uint32
	useTime        uint32
	curElectricity uint32
	chipReset      uint16
	maxElectricity uint16
	errCode        uint8
	enable         uint8
	behavior       uint8
}

type ComList struct {
	signal         uint8
	comNum         uint8
	comID          []byte
	comPort        []interface{}
	comBehavior    uint8
	comProtoVer    uint16
	useEnergy      uint32
	useElectricity uint32
	enableCount    uint8
}

type JosnPacket struct {
	Act      string `json:"act"`
	ID       int    `json:"id"`
	Ctime    int    `json:"ctime"`
	Data     string `json:"data"`
	Behavior int    `json:"behavior"`
}

func (update *UpdateState) Print() (retString string) {
	retString = "update_state"
	return

}

func (result *TransferResult) Print() (retString string) {
	retString = "transfer_result"
	return
}

func (device *DeviceInfo) Print() (retString string) {
	retString = "device_info"
	return
}

func (comList *ComList) Print() (retString string) {
	var buffer bytes.Buffer
	if comList != nil {
		headString := fmt.Sprintf("signal:-%d comNum:%02d comIDs:%s TaskBehavior:%d Version:%d useEnergy:%d useElectricity:%d enableCount:%d\n", comList.signal, comList.comNum, hex.EncodeToString(comList.comID), comList.comBehavior, comList.comProtoVer, comList.useEnergy, comList.useElectricity, comList.enableCount)
		buffer.WriteString(headString)
		//buffer.WriteString("=====================================================\n")
		for comID, data := range comList.comPort {
			comData := data.(ComData)
			var comString string

			switch comList.comProtoVer {
			case MAX_PROTO_VERSION0:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Uelectricity:%06d errCode:%03d\n", comID, comData.token, comData.maxEnergy, comData.useEnergy,
						comData.maxTime, comData.useTime, comData.curElectricity, comData.errCode)
					break
				}
			case MAX_PROTO_VERSION1:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Melectricity:%06d Uelectricity:%06d reset:%06d errCode:%03d \n", comID, comData.token, comData.maxEnergy, comData.useEnergy,
						comData.maxTime, comData.useTime, comData.maxElectricity,
						comData.curElectricity, comData.chipReset, comData.errCode)
					break
				}
			case MAX_PROTO_VERSION2:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Melectricity:%06d Uelectricity:%06d reset:%06d errCode:%03d enable:%01d \n", comID, comData.token, comData.maxEnergy, comData.useEnergy,
						comData.maxTime, comData.useTime, comData.maxElectricity,
						comData.curElectricity, comData.chipReset, comData.errCode,
						comData.enable)
					break
				}
			case MAX_PROTO_VERSION3:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Melectricity:%06d Uelectricity:%06d reset:%06d errCode:%03d enable:%01d behavior:%03d \n", comID, comData.token, comData.maxEnergy, comData.useEnergy,
						comData.maxTime, comData.useTime, comData.maxElectricity,
						comData.curElectricity, comData.chipReset, comData.errCode,
						comData.enable, comData.behavior)
					break
				}
			}
			buffer.WriteString(comString)
		}
	}

	retString = buffer.String()
	return
}

func (taskStart *ComTaskStartTransfer) Print() (retString string) {
	var buffer bytes.Buffer
	if taskStart != nil {
		comString := fmt.Sprintf("ComTaskStartTransfer----------> comID:%02d token:%012d M_energy:%06d M_electricity:%06d M_time:%08d \n", taskStart.comID, taskStart.token, taskStart.maxEnergy, taskStart.maxElectricity, taskStart.maxTime)
		buffer.WriteString(comString)
	}
	retString = buffer.String()
	return
}

func (taskStop *ComTaskStopTransfer) Print() (retString string) {
	var buffer bytes.Buffer
	if taskStop != nil {
		comString := fmt.Sprintf("ComTaskStopTransfer----------> comID:%02d token:%012d forceStop:%02d\n", taskStop.comID, taskStop.token, taskStop.forceStop)
		buffer.WriteString(comString)
	}
	retString = buffer.String()
	return
}

func (json *JosnPacket) Print() (retString string) {
	var buffer bytes.Buffer
	buffer.WriteString(json.FormatData())
	buffer.WriteString(" behavior:")
	buffer.WriteString(strconv.Itoa(json.Behavior))
	retString = buffer.String()
	return
}

func (json *JosnPacket) FormatData() (data string) {
	data = ""
	if json != nil {
		if len(json.Data) > 0 {
			data = strings.Replace(json.Data, "\\", "", -1)
		}
	}
	return
}
