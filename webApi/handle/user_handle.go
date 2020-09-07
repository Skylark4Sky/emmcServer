package handle

import (
	. "GoServer/webApi/middleWare"
	. "GoServer/webApi/model"
	. "GoServer/webApi/utils"
	"errors"
	"fmt"
)

// 登录成功返回
type UserLoginRespond struct {
	UserID    int64
	UserName  string
	NickName  string
	Gender    int8
	Birthday  int64
	Signature string
	Mobile    string
	Email     string
	Face      string
	Face200   string
	Srcface   string
}

// 登录绑定
type UserLogin struct {
	UserBase UserBase
	Account  string `json:"account"`
	Pwsd     string `json:"pwsd"`
}

func createLoginRespond(entity *UserBase) *UserLoginRespond {
	return &UserLoginRespond{
		UserID:    entity.UId,
		UserName:  entity.UserName,
		NickName:  entity.NickName,
		Gender:    entity.Gender,
		Birthday:  entity.Birthday,
		Signature: entity.Signature,
		Mobile:    entity.Mobile,
		Email:     entity.Email,
		Face:      entity.Face,
		Face200:   entity.Face200,
		Srcface:   entity.Srcface,
	}
}

func (M *UserLogin) Login(ip string) (*JwtObj, error) {
	if M.Pwsd == "" {
		return nil, errors.New("password is required")
	}
	entity := &M.UserBase
	cond := fmt.Sprintf("email = '%s' or user_name = '%s' or mobile = '%s'", M.Account, M.Account, M.Account)
	err := DBInstance.Debug().Where(cond).First(&entity).Error
	if err != nil {
		if IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if chkOk := PasswordVerify(M.Pwsd, entity.UserPwsd); chkOk != true {
		return nil, err
	}

	return JwtGenerateToken(createLoginRespond(entity), entity.UId)
}
