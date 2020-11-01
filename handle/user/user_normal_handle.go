package user

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/extension"
	. "GoServer/model"
	. "GoServer/model/user"
	. "GoServer/utils/respond"
	. "GoServer/utils/security"
	. "GoServer/utils/time"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type PermissionAction struct {
	Action       string `json:"action"`
	DefaultCheck uint8  `json:"defaultCheck"`
	Describe     string `json:"describe"`
}

//权限
type Permission struct {
	PermissionId    string             `json:"permissionId"`
	PermissionName  string             `json:"permissionName"`
	Actions         []PermissionAction `json:"actions"`
	ActionEntitySet []PermissionAction `json:"actionEntitySet"`
	ActionList      []PermissionAction `json:"actionList"`
	DataAccess      []PermissionAction `json:"dataAccess"`
}

// 查询
type UserInfo struct {
	User      UserData     `json:"user"`
	RulesList []Permission `json:"permissions,omitempty"`
}

// 登录
type UserLogin struct {
	UserBase UserBase
	Account  string `form:"account" json:"account" binding:"required"`
	Pwsd     string `form:"pwsd" json:"pwsd" binding:"required"`
}

// 注册
type UserRegister struct {
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

func getLoginType(account string, entity *UserData) uint8 {
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

func analysisRoleList(entity *UserInfo, roleMenus *[]UserRoleMenus) {
	var rootDict = make(map[int16]interface{})
	var menuDict = make(map[int16]interface{})
	entity.RulesList = make([]Permission, 0)

	for _, node := range *roleMenus {
		if node.Type == "page" {
			p := &Permission{
				PermissionId:    node.Key,
				PermissionName:  node.Name,
				ActionEntitySet: make([]PermissionAction, 0),
			}
			rootDict[node.ID] = p
		} else {
			p := &PermissionAction{
				Action:       node.Key,
				DefaultCheck: node.DefaultCheck,
				Describe:     node.Name,
			}
			rootDict[node.ID] = p
		}
	}

	for _, node := range *roleMenus {
		if node.Type == "page" {
			menu := rootDict[node.ID].(*Permission)
			menuDict[node.ID] = menu
			entity.RulesList = append(entity.RulesList, *menu)
		} else {
			parent := menuDict[node.PID].(*Permission)
			action := rootDict[node.ID].(*PermissionAction)
			parent.ActionEntitySet = append(parent.ActionEntitySet, *action)
		}
	}

	for _, node := range menuDict {
		menu := *(node.(*Permission))
		if len(menu.ActionEntitySet) <= 0 {
			menu.ActionEntitySet = nil
		}
	}
}

func fetchUserRules(entity *UserInfo) {
	if entity.User.Rules != "" && len(entity.User.Rules) >= 1 {
		countSplit := strings.Split(entity.User.Rules, ",")
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

func FetchUserInfo(userID uint64,ctx *gin.Context)  (*UserInfo, interface{}) {

	userInfo := &UserInfo{}

	db := ExecSQL().Table("user_base")
	db = db.Select("user_base.uid,user_base.user_name,user_base.user_pwsd,user_base.nick_name,user_base.gender,user_base.birthday,user_base.signature,user_base.face200,user_base.mobile,user_base.email,user_role.rules")
	db = db.Joins("inner join user_role ON user_base.user_role = user_role.id")
	db = db.Where("uid = ?", userID)

	//err := ExecSQL().Table("user_base").Select("user_base.uid,user_base.user_name,user_base.user_pwsd,user_base.nick_name,user_base.gender,user_base.birthday,user_base.signature,user_base.face200,user_base.mobile,user_base.email,user_role.rules").Joins("inner join user_role ON user_base.user_role = user_role.id").Where("uid = ?", userID).Scan(&userInfo.User).Error
	if err := db.Scan(&userInfo.User).Error; err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateErrorMessage(USER_NO_EXIST, nil)
		}
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	//检出菜单
	fetchUserRules(userInfo)

	return userInfo,nil
}

func (M *UserLogin) Run(ctx *gin.Context) (*LoginRespond, interface{}) {
	userInfo := &UserInfo{}
	db := ExecSQL().Table("user_base")
	db = db.Select("user_base.uid,user_base.user_name,user_base.user_pwsd,user_base.nick_name,user_base.gender,user_base.birthday,user_base.signature,user_base.face200,user_base.mobile,user_base.email,user_role.rules")
	db = db.Joins("inner join user_role ON user_base.user_role = user_role.id")
	db = db.Where("email = ? or user_name = ? or mobile = ?", M.Account, M.Account, M.Account)

	if err := db.Scan(&userInfo.User).Error; err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateErrorMessage(USER_NO_EXIST, nil)
		}
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	var loginType uint8 = getLoginType(M.Account, &userInfo.User)

	if chkOk := PasswordVerify(M.Pwsd, userInfo.User.UserPwsd); chkOk != true {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, userInfo.User.UID)
		return nil, CreateErrorMessage(USER_PWSD_ERROR, nil)
	}

	//检出菜单
	fetchUserRules(userInfo)

	tokenData, err := JwtGenerateToken(userInfo.User.UID)
	if err != nil {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, userInfo.User.UID)
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	createLoginLog(ctx, LOGIN_SUCCEED, loginType, userInfo.User.UID)

	respond := LoginRespond{
		UserInfo: userInfo,
		Token:    tokenData,
	}

	return &respond, nil
}

func (M *UserRegister) Build(ctx *gin.Context) interface{} {
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
