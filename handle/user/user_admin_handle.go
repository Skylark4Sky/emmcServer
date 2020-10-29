package user

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/extension"
	. "GoServer/model"
	. "GoServer/model/user"
	. "GoServer/utils/log"
	. "GoServer/utils/respond"
	. "GoServer/utils/security"
	. "GoServer/utils/time"
	//	"bytes"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type PermissionAction struct {
	action       string `json:"action"`
	defaultCheck uint8  `json:"defaultCheck"`
	describe     string `json:"describe"`
}

//权限
type Permission struct {
	permissionId    string             `json:"permissionId"`
	permissionName  string             `json:"permissionName"`
	actions         []PermissionAction `josn:"actions"`
	actionEntitySet []PermissionAction `json:"actionEntitySet"`
	actionList      []PermissionAction `json:"actionList"`
	dataAccess      []PermissionAction `json:"dataAccess"`
}

// 查询
type AdminUser struct {
	UID       uint64        `json:"uid"`
	UserName  string        `json:"username"`
	UserPwsd  string        `json:"-"`
	NickName  string        `json:"nickname"`
	Gender    uint8         `json:"gender"`
	Birthday  int64         `json:"birthday"`
	Signature string        `json:"signature"`
	Face200   string        `json:"face200"`
	Mobile    string        `json:"mobile"`
	Email     string        `json:"email"`
	Rules     string        `json:"-"`
	RulesList []interface{} `json:"menus"`
}

// 登录
type AdminLogin struct {
	UserBase UserBase
	Account  string `form:"account" json:"account" binding:"required"`
	Pwsd     string `form:"pwsd" json:"pwsd" binding:"required"`
}

// 注册
type AdminRegister struct {
	Source    uint8  `form:"source" json:"source" binding:"required"`
	Name      string `form:"userName" json:"userName"`
	Pwsd      string `form:"userPwsd" json:"userPwsd" binding:"required"`
	NickName  string `form:"nickName" json:"nickName"`
	Gender    uint8  `form:"gender" json:"gender"`
	Birthday  int64  `form:"birthDay" json:"birthDay"`
	Signature string `form:"signature" json:"signature"`
	Mobile    string `form:"mobile" json:"mobile"`
	Email     string `form:"email" json:"email"`
}

func getLoginType(account string, entity *AdminUser) uint8 {
	loginType := UNKNOWN

	switch account {
	case entity.Email:
		loginType = EMAIL
	case entity.UserName:
		loginType = USERNAME
	case entity.Mobile:
		loginType = MOBILE
	}
	return loginType
}

func (entity *Permission) setAction(action PermissionAction) {
	entity.actionEntitySet = append(entity.actionEntitySet, action)
}

func analysisRoleList(entity *AdminUser, roleMenus *[]UserRoleMenus) {

	var rootDict = make(map[int16]interface{})
	var menuList = make([]Permission, 0)

	for _, node := range *roleMenus {
		if node.Type == "page" {
			var actionList = make([]PermissionAction, 0)
			var p Permission = Permission{
				permissionId:    node.Key,
				permissionName:  node.Name,
				actionEntitySet: actionList,
			}
			rootDict[node.ID] = p
		} else {
			var p PermissionAction = PermissionAction{
				action:       node.Key,
				defaultCheck: node.DefaultCheck,
				describe:     node.Name,
			}
			rootDict[node.ID] = p
		}
	}

	for _, node := range *roleMenus {
		if node.Type == "page" {
			menu := rootDict[node.ID].(Permission)
			menuList = append(menuList, menu)
		} else {
			parent := rootDict[node.PID].(Permission)
			action := rootDict[node.ID].(PermissionAction)
			parent.setAction(action)
		}
	}

	SystemLog("menuList:-->", menuList)

}

func fetchUserRules(entity *AdminUser) {
	if entity.Rules != "" && len(entity.Rules) >= 1 {
		countSplit := strings.Split(entity.Rules, ",")
		ruleids := make([]int16, len(countSplit))
		for idx, ids := range countSplit {
			if v, err := strconv.Atoi(ids); err == nil {
				ruleids[idx] = int16(v)
			}
		}

		if len(ruleids) > 0 {
			var roleMenus []UserRoleMenus
			ExecSQL().Where("id IN (?)", ruleids).Find(&roleMenus)
			analysisRoleList(entity, &roleMenus)
		}
	}
}

func (M *AdminLogin) Login(ctx *gin.Context) (*JwtObj, interface{}) {
	adminResults := &AdminUser{}
	err := ExecSQL().Debug().Table("user_base").Select("user_base.uid,user_base.user_name,user_base.user_pwsd,user_base.nick_name,user_base.gender,user_base.birthday,user_base.signature,user_base.face200,user_base.mobile,user_base.email,user_role.rules").Joins("inner join user_role ON user_base.user_role = user_role.id").Where("email = ? or user_name = ? or mobile = ?", M.Account, M.Account, M.Account).Scan(&adminResults).Error

	if err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateErrorMessage(USER_NO_EXIST, nil)
		}
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	var loginType uint8 = getLoginType(M.Account, adminResults)

	if chkOk := PasswordVerify(M.Pwsd, adminResults.UserPwsd); chkOk != true {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, adminResults.UID)
		return nil, CreateErrorMessage(USER_PWSD_ERROR, nil)
	}

	//检出菜单
	fetchUserRules(adminResults)

	JwtData, err := JwtGenerateToken(adminResults, adminResults.UID)
	if err != nil {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, adminResults.UID)
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	createLoginLog(ctx, LOGIN_SUCCEED, loginType, adminResults.UID)
	return JwtData, nil
}

func (M *AdminRegister) Register(ctx *gin.Context) interface{} {

	var user CreateUserInfo
	user.Base = UserBase{
		RegisterSource: M.Source,
		UserRole:       ADMIN_USER,
		UserName:       M.Name,
		UserPwsd:       M.Pwsd,
		NickName:       M.NickName,
		Gender:         M.Gender,
		Birthday:       M.Birthday,
		Signature:      M.Signature,
		Mobile:         M.Mobile,
		Email:          M.Email,
		CreateTime:     GetTimestamp(),
	}

	user.Log = UserRegisterLog{
		RegisterMethod: M.Source,
		RegisterTime:   user.Base.CreateTime,
		RegisterIP:     ctx.ClientIP(),
	}

	user.Extra = UserExtra{
		CreateTime: user.Base.CreateTime,
	}

	user.Location = UserLocation{}

	CreateAsyncSQLTask(ASYNC_CREATE_NORMAL_USER, user)

	return CreateMessage(SUCCESS, nil)
}
