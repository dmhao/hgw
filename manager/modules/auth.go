package modules

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"time"
)

const Issuer = "hgw_admin"
const TokenExpire = 24 * time.Hour


type AdminUser struct {
	UserId			string
	UserName		string
	Password		string
	Salt			string
}

func Index(_ *gin.Context) {

}

func AuthHandler(c *gin.Context) {
	jwtStr,_ := c.Cookie("jwt")
	userId,_ := c.Cookie("userId")
	if jwtStr == "" || userId == "" {
		mContext{c}.ErrorOP(SystemErrorNotLogin)
		return
	}
	userRsp, err := adminUser(userId)
	if err != nil || userRsp.Count == 0 {
		mContext{c}.ErrorOP(SystemErrorNotLogin)
		return
	}

	user := new(AdminUser)
	err = json.Unmarshal(userRsp.Kvs[0].Value, user)
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	token, err := checkToken(userId, user.Salt, jwtStr)
	if token != nil && err == nil {
		if token.Valid {
			claims := token.Claims.(*jwt.StandardClaims)
			if userId == claims.Subject {
				c.Next()
				return
			}
		}
	}
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
}

func AuthInit(c *gin.Context) {
	rsp, err := authInitData()
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	if rsp.Count > 0 {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	salt := uuid.Must(uuid.NewV4()).String()
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}
	userId := createUserNameMd5(username)
	pwMd5 := createPasswordMd5(password, salt)

	adminUser := AdminUser{userId,username, pwMd5, salt}
	userBytes, err := json.Marshal(adminUser)
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}

	result := putAdminUser(userId, string(userBytes))
	if !result {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	mContext{c}.SuccessOP(nil)
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.SetCookie("userId", "", -1, "/", "", false, true)
	mContext{c}.SuccessOP(nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		mContext{c}.ErrorOP(DataParseError)
		return
	}
	userId := createUserNameMd5(username)
	userRsp, err := adminUser(userId)
	if err != nil || userRsp.Count == 0 {
		mContext{c}.ErrorOP(LoginParamsError)
		return
	}

	user := new(AdminUser)
	err = json.Unmarshal(userRsp.Kvs[0].Value, user)
	if err != nil {
		mContext{c}.ErrorOP(SystemError)
		return
	}

	if user.Password != createPasswordMd5(password, user.Salt) {
		mContext{c}.ErrorOP(LoginParamsError)
		return
	}

	token := getToken(userId, user.Salt)
	if token == "" {
		mContext{c}.ErrorOP(SystemError)
		return
	}
	c.SetCookie("jwt", token, int(TokenExpire.Seconds()), "/", "", false, false)
	c.SetCookie("userId", userId, int(TokenExpire.Seconds()), "/", "", false, false)
	mContext{c}.SuccessOP(nil)
}

func createUserNameMd5(username string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(username))
	userMd5 := md5Ctx.Sum(nil)
	return hex.EncodeToString(userMd5)
}

func createPasswordMd5(password, salt string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(salt))
	md5Ctx.Write([]byte(password))
	md5Ctx.Write([]byte(salt))
	passwordMd5 := md5Ctx.Sum(nil)
	return hex.EncodeToString(passwordMd5)
}

func getToken(userId string, secretKey string) string {
	signKey := []byte(secretKey)
	claims := &jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: time.Now().Add(TokenExpire).Unix(),
		Issuer:    Issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signKey)
	if err != nil {
		return ""
	}
	return ss
}

func checkToken(userId string, secretKey string, tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return token, err
}