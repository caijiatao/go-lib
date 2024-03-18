package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golib/libs/gin_helper"
	"golib/libs/logger"
	"golib/libs/redis"
	"net/http"
	"time"
)

const AuthSecret = "SECRET"
const CtxUserID = "CTX-USERID"
const CtxUserPhone = "CTX-USER-PHONE"

const TokenKey = "TOKEN-KEY"
const TokenExpire = 24 * time.Hour

type User struct {
	UserID      int64
	PhoneNumber string
}

func AuthenticationVarious(fn func(ctx *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ParseToken(c)
		if err != nil {
			OutPutErr(c, err)
			return
		}
		err = fn(c)
		if err != nil {
			OutPutErr(c, err)
			return
		}
	}
}

func OutPutErr(ctx *gin.Context, err error) {
	myErr := err.(*gin_helper.MyError)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":       myErr.Code,
		"msg":        myErr.Msg,
		"serverTime": myErr.ServerTime,
		"data":       myErr.Data,
	})
	return
}

// ParseToken 解析验证jwt
func ParseToken(ctx *gin.Context) error {
	authToken := ctx.GetHeader("Token")
	if authToken == "" {
		logger.Errorf("refuse unauthorized request")
		return gin_helper.UNAUTHORIZED
	}
	token, err := jwt.Parse(authToken, func(*jwt.Token) (interface{}, error) {
		return []byte(AuthSecret), nil
	})
	if err != nil {
		logger.Errorf("token parse error: %s", err)
		return gin_helper.UNAUTHORIZED
	}

	// 将token的claims放入到上下文中
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Errorf("token turn to claims")
		return gin_helper.INNER_ERROR
	}
	// 计算token是否超时
	exp := int64(claims["exp"].(float64))
	if exp < time.Now().Unix() {
		userPhone, _ := claims["phone_number"].(string)
		err = DoLogOut(ctx, authToken, userPhone)
		if err != nil {
			return gin_helper.INNER_ERROR
		}
		return gin_helper.TOKEN_TIMEOUT
	}

	// 缓存查询
	userPhone := redis.ClusterClient().Get(ctx, TokenKey+":token:"+authToken).Val()
	if userPhone == "" {
		return gin_helper.UNAUTHORIZED
	}

	ctx.Set("token-claims", claims)
	ctx.Set(CtxUserID, claims["user_id"])
	ctx.Set(CtxUserPhone, claims["phone_number"])
	return nil
}

// GenerateToken 生成jwt
func GenerateToken(userID int64, phoneNumber string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// 过期时间2h
	claims["exp"] = time.Now().Add(TokenExpire).Unix()
	claims["iat"] = time.Now().Unix()
	claims["issuer"] = "sics"
	claims["user_id"] = userID
	claims["phone_number"] = phoneNumber
	token.Claims = claims

	return token.SignedString([]byte(AuthSecret))
}

func GetUser(ctx *gin.Context) *User {
	id, exist := ctx.Get(CtxUserID)
	if !exist {
		return nil
	}
	uid, _ := id.(float64)

	phone, exist := ctx.Get(CtxUserPhone)
	if !exist {
		return nil
	}
	userPhone, _ := phone.(string)

	return &User{
		UserID:      int64(uid),
		PhoneNumber: userPhone,
	}
}

func DoLogOut(ctx context.Context, token, userPhone string) error {
	_, err := redis.ClusterClient().Del(ctx, TokenKey+":user:"+userPhone).Result()
	if err != nil {
		return err
	}
	_, err = redis.ClusterClient().Del(ctx, TokenKey+":token:"+token).Result()
	if err != nil {
		return err
	}
	return nil
}

func DoLogin(ctx context.Context, token, userPhone string) error {
	lastToken := redis.ClusterClient().Get(ctx, TokenKey+":user:"+userPhone).Val()
	if lastToken != "" {
		_, err := redis.ClusterClient().Del(ctx, TokenKey+":token:"+lastToken).Result()
		if err != nil {
			return err
		}
	}
	redis.ClusterClient().SetEx(ctx, TokenKey+":token:"+token, userPhone, TokenExpire)
	redis.ClusterClient().SetEx(ctx, TokenKey+":user:"+userPhone, token, TokenExpire)
	return nil
}
