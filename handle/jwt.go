package handle

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	// ErrorReasonServerBusy .
	ErrorReasonServerBusy = "服务器繁忙"
	// ErrorReasonReLogin .
	ErrorReasonReLogin = "请重新登陆"
)

func sayHello(c *gin.Context) {
	strToken := c.Param("token")
	claim, err := VerifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, "hello,", claim.Phone)
}

// JWTClaims token里面添加用户信息，验证token后可能会用到用户信息
type JWTClaims struct {
	jwt.StandardClaims
	UserID      uint     `json:"user_id"`
	Password    string   `json:"password"`
	Phone       string   `json:"username"`
	FullName    string   `json:"full_name"`
	Permissions []string `json:"permissions"`
}

var (
	// Secret 加盐
	Secret = "dong_tech"
	// ExpireTime token有效期
	ExpireTime = 3600
)

func login(c *gin.Context) {
	username := c.Param("username")
	password := c.Param("password")
	claims := &JWTClaims{
		UserID:      1,
		Phone:       username,
		Password:    password,
		FullName:    username,
		Permissions: []string{},
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	signedToken, err := GetToken(claims)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, signedToken)
}

func verify(c *gin.Context) {
	// strToken := c.Param("token")
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 2003,
			"msg":  "请求头中auth为空",
		})
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusOK, gin.H{
			"code": 2004,
			"msg":  "请求头中auth格式有误",
		})
		c.Abort()
		return
	}
	strToken := parts[1]

	claim, err := VerifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, "verify,", claim.Phone)
}

func refresh(c *gin.Context) {
	strToken := c.Param("token")
	claims, err := VerifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := GetToken(claims)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, signedToken)
}

// VerifyAction 验证token
func VerifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, errors.New("token错误")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReasonReLogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorReasonReLogin)
	}
	return claims, nil
}

// GetToken 创建一个新的token
func GetToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", errors.New(ErrorReasonServerBusy)
	}
	return signedToken, nil
}
