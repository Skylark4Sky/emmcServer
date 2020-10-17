package Mqtt

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const DEVICE_INFO = "device_info"
const TRANSFER = "transfer"
const TRANSFER_RESULT = "resp"
const FIRMWARE_UPDATE = "update_ver"

const (
	COM_NO_WORKING     = 0
	COM_WORKING        = 1
	MIN_PROTO_VERSION0 = 27
	MIN_PROTO_VERSION1 = 32
	MIN_PROTO_VERSION2 = 33
	MIN_PROTO_VERSION3 = 34
	MAX_PROTO_VERSION0 = 261
	MAX_PROTO_VERSION1 = 302
	MAX_PROTO_VERSION2 = 312
	MAX_PROTO_VERSION3 = 322
	MAX_PROTO_VERSION4 = 162
)

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
	Imei     string `json:"imei"`
	Version  string `json:"version"`
	DeviceSn string `json:"device_sn"`
	Sim      struct {
		ICCID string `json:"ICCID"`
		IMEI  string `json:"IMEI"`
	} `json:"sim"`
	CellInfo struct {
		Ci  int    `json:"ci"`
		Lac int    `json:"Lac"`
		Mnc int    `json:"Mnc"`
		Mcc int    `json:"Mcc"`
		Ta  int    `json:"ta"`
		Ext string `json:"Ext"`
	} `json:"cellInfo"`
}

type ComTaskStartTransfer struct {
	ComID          uint8  `json:"id"`          //端口号
	Token          uint32 `json:"token"`       //令牌
	MaxEnergy      uint32 `json:"energy"`      //最大电量
	MaxElectricity uint32 `json:"electricity"` //最大电流
	MaxTime        uint32 `json:"time"`        //最大时间
}

type ComTaskStopTransfer struct {
	ComID     uint8  `json:"id"`        //端口号
	Token     uint32 `json:"token"`     //令牌
	ForceStop uint8  `json:"forceStop"` //是否强制停止 是 停止并清空当前端口数据
}

type ComTaskStatusQueryTransfer struct {
	ComID uint8 `json:"id"` //端口号 0 - 9  大于9 查全部
}

type DeviceSetConfigTransfer struct {
	Time uint8 `json:"time"` //空载时间 秒
}

type DeviceReStartTaskTransfer struct {
}

//某一端口数据
type ComData struct {
	Id             uint8   `json:"id"`             //端口ID
	Token          uint32  `json:"token"`          //令牌
	MaxEnergy      uint32  `json:"maxEnergy"`      //最大电量
	UseEnergy      uint32  `json:"useEnergy"`      //冲电量
	MaxTime        uint32  `json:"maxTime"`        //最大时间
	UseTime        uint32  `json:"useTime"`        //已用时间
	CurElectricity uint32  `json:"curElectricity"` //当前电流
	MaxElectricity uint16  `json:"maxElectricity"` //最大电流
	ChipReset      uint16  `json:"chipReset"`      //芯片复位统计
	ErrCode        uint8   `json:"errCode"`        //错误码
	Enable         uint8   `json:"enable"`         //是否启用
	Behavior       uint8   `json:"behavior"`       //最后状态行为
	CurPower       float64 `json:"c_power"`        //当前端口使用功率
	AveragePower   float64 `json:"a_power"`        //当前端口平均功率
	MaxPower       float64 `json:"m_power"`        //当前端口最高使用功率
}

// 上报数据
type ComList struct {
	Signal         uint8         //设备网络信号
	ComNum         uint8         //端口数量
	ComID          []byte        //端口ID数列
	ComPort        []interface{} //端口数据
	ComBehavior    uint8         //上报行为
	ComProtoVer    uint16        //上报协议版本
	UseEnergy      uint32        //当前使用电量
	UseElectricity uint32        //当前使用电流
	EnableCount    uint8         //充电统计
}

type JosnPacket struct {
	Act      string `json:"act"`
	ID       int    `json:"id"`
	Ctime    int    `json:"ctime,omitempty"`
	Data     string `json:"data"`
	Behavior uint8    `json:"behavior"`
}

type Protocol interface {
	Print() (retString string)
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
		headString := fmt.Sprintf("signal:-%d comNum:%02d comIDs:%s TaskBehavior:%d Version:%d useEnergy:%d useElectricity:%d enableCount:%d\n", comList.Signal, comList.ComNum, hex.EncodeToString(comList.ComID), comList.ComBehavior, comList.ComProtoVer, comList.UseEnergy, comList.UseElectricity, comList.EnableCount)
		buffer.WriteString(headString)
		//buffer.WriteString("=====================================================\n")
		for comID, data := range comList.ComPort {
			comData := data.(ComData)
			var comString string

			switch comList.ComProtoVer {
			case MAX_PROTO_VERSION0:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Uelectricity:%06d errCode:%03d\n", comID, comData.Token, comData.MaxEnergy, comData.UseEnergy,
						comData.MaxTime, comData.UseTime, comData.CurElectricity, comData.ErrCode)
					break
				}
			case MAX_PROTO_VERSION1:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Melectricity:%06d Uelectricity:%06d reset:%06d errCode:%03d \n", comID, comData.Token, comData.MaxEnergy, comData.UseEnergy,
						comData.MaxTime, comData.UseTime, comData.MaxElectricity,
						comData.CurElectricity, comData.ChipReset, comData.ErrCode)
					break
				}
			case MAX_PROTO_VERSION2:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Melectricity:%06d Uelectricity:%06d reset:%06d errCode:%03d enable:%01d \n", comID, comData.Token, comData.MaxEnergy, comData.UseEnergy,
						comData.MaxTime, comData.UseTime, comData.MaxElectricity,
						comData.CurElectricity, comData.ChipReset, comData.ErrCode,
						comData.Enable)
					break
				}
			case MAX_PROTO_VERSION3:
				{
					comString = fmt.Sprintf("comID:%02d token:%012d Menergy:%06d Uenergy:%06d Mtime:%08d Utime:%08d Melectricity:%06d Uelectricity:%06d reset:%06d errCode:%03d enable:%01d behavior:%03d \n", comID, comData.Token, comData.MaxEnergy, comData.UseEnergy,
						comData.MaxTime, comData.UseTime, comData.MaxElectricity,
						comData.CurElectricity, comData.ChipReset, comData.ErrCode,
						comData.Enable, comData.Behavior)
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
		comString := fmt.Sprintf("ComTaskStartTransfer----------> comID:%02d token:%012d M_energy:%06d M_electricity:%06d M_time:%08d \n", taskStart.ComID, taskStart.Token, taskStart.MaxEnergy, taskStart.MaxElectricity, taskStart.MaxTime)
		buffer.WriteString(comString)
	}
	retString = buffer.String()
	return
}

func (taskStop *ComTaskStopTransfer) Print() (retString string) {
	var buffer bytes.Buffer
	if taskStop != nil {
		comString := fmt.Sprintf("ComTaskStopTransfer----------> comID:%02d token:%012d forceStop:%02d\n", taskStop.ComID, taskStop.Token, taskStop.ForceStop)
		buffer.WriteString(comString)
	}
	retString = buffer.String()
	return
}

func (taskStatusQuery *ComTaskStatusQueryTransfer) Print() (retString string) {
	var buffer bytes.Buffer
	var comString string

	if taskStatusQuery != nil {
		if taskStatusQuery.ComID < 10 {
			comString = fmt.Sprintf("ComTaskStatusQueryTransfer----------> comID:%02d\n", taskStatusQuery.ComID)
		} else {
			comString = fmt.Sprintf("ComTaskStatusQueryTransfer----------> Query All Com:%02d\n", taskStatusQuery.ComID)
		}
		buffer.WriteString(comString)
	}
	retString = buffer.String()
	return
}

func (taskSetConfig *DeviceSetConfigTransfer) Print() (retString string) {
	var buffer bytes.Buffer
	if taskSetConfig != nil {
		comString := fmt.Sprintf("DeviceSetConfigTransfer----------> on_load time:%02d\n", taskSetConfig.Time)
		buffer.WriteString(comString)
	}
	retString = buffer.String()
	return
}

func (restart *DeviceReStartTaskTransfer) Print() (retString string) {
	retString = "reStart Device"
	return
}

func (json *JosnPacket) Print() (retString string) {
	var buffer bytes.Buffer
	buffer.WriteString(json.formatData())
	buffer.WriteString(" behavior:")
	buffer.WriteString(strconv.Itoa(int(json.Behavior)))
	retString = buffer.String()
	return
}

func (json *JosnPacket) formatData() (data string) {
	data = ""
	if json != nil {
		if len(json.Data) > 0 {
			data = strings.Replace(json.Data, "\\", "", -1)
		}
	}
	return
}
