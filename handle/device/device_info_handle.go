package device

type RequestListData struct {
	UserID    uint64 `fomr:"userID" json:"userID" binding:"required"`
	PageNum   int64  `form:"pageNum" json:"pageNum" binding:"required"`   //起始页
	PageSize  int64  `form:"pageSize" json:"pageSize" binding:"required"` //每页大小
	StartTime int64  `form:"startTime" json:"startTime"`
	EndTime   int64  `form:"endTime" json:"startTime"`
}

func (request *RequestListData) GetDeviceList() (interface{}, interface{}) {
	return nil,nil
}

func (request *RequestListData) GetDeviceTransferLogList() (interface{}, interface{}) {
	return nil,nil
}

func (request *RequestListData) GetModuleList() (interface{}, interface{}) {
	return nil,nil
}

func (request *RequestListData) GetModuleConnectLogList() (interface{}, interface{}) {
	return nil,nil
}
