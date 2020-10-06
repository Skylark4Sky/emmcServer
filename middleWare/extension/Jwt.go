package middleWare

import (
	. "GoServer/utils/config"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Jwt, _ = GetJwt()

var (
	AppSecret  = Jwt.AppSecret
	AppIss     = Jwt.AppIss
	ExpireTime = time.Hour * time.Duration(Jwt.ExpireTime)
)

const (
	JwtCtxUidKey = "authedUserId"
	bearerLength = len("Bearer ")
)

type JwtObj struct {
	User     interface{}
	Token    string    `json:"token"`
	Expire   time.Time `json:"expire"`
	ExpireTs int64     `json:"expire_ts"`
}

func jwtTokenVerify(tokenString string) (int64, error) {
	if tokenString == "" {
		return 0, errors.New("no token is found in Authorization Bearer")
	}
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppSecret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return 0, errors.New("token is expired")
	}
	if !claims.VerifyIssuer(AppIss, true) {
		return 0, errors.New("token's issuer is wrong")
	}
	uid, err := strconv.ParseUint(claims.Id, 10, 64)
	return int64(uid), err
}

func JwtGenerateToken(obj interface{}, userID int64) (*JwtObj, error) {
	expireTime := time.Now().Add(ExpireTime)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", userID),
		Issuer:    AppIss,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		logrus.WithError(err).Error("config is wrong, can not generate jwt")
	}
	data := &JwtObj{User: obj, Token: tokenString, Expire: expireTime, ExpireTs: expireTime.Unix()}
	return data, err
}

func JwtIntercept(context *gin.Context) {
	token, ok := context.GetQuery("_t")

	if !ok {
		//		fmt.Println("!ok")
		hToken := context.GetHeader("Authorization")

		//		fmt.Println("hToken:", hToken)
		if len(hToken) < bearerLength {
			context.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
			return
		}
		token = strings.TrimSpace(hToken[bearerLength:])
	}

	userId, err := jwtTokenVerify(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": err.Error()})
		return
	}

	context.Set(JwtCtxUidKey, userId)
	context.Next()
}
