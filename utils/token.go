package utils

import (
	"MyGoApi/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/astaxie/beego/context"
)

const (
	KEY string = "JWT-ARY-STARK"
	DEFAULT_EXPIRE_SECONDS int = 7200 // default 120 minutes

)




// JWT -- json web token
// HEADER PAYLOAD SIGNATURE
// This struct is the PAYLOAD
type MyCustomClaims struct {
	models.Users
	jwt.StandardClaims
}


// update expireAt and return a new token
func RefreshToken(tokenString string)(string, error) {
	// first get previous token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token)(interface{}, error) {
			return []byte(KEY), nil
		})
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return "", err
	}
	mySigningKey := []byte(KEY)
	expireAt  := time.Now().Add(time.Second * time.Duration(DEFAULT_EXPIRE_SECONDS)).Unix()
	newClaims := MyCustomClaims{
		claims.Users,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    claims.Users.Name,
			IssuedAt:  time.Now().Unix(),
		},
	}
	// generate new token with new claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenStr, err := newToken.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("generate new fresh json web token failed !! error :", err)
		return  "" , err
	}
	return tokenStr, err
}


func ValidateToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token)(interface{}, error) {
			return []byte(KEY), nil
		})
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.Users, claims.StandardClaims.ExpiresAt)
		fmt.Println("token will be expired at ", time.Unix(claims.StandardClaims.ExpiresAt, 0))
	} else {
		fmt.Println("validate tokenString failed !!!",err)
		return err
	}
	return nil
}


func GenerateToken(expiredSeconds int, user *models.Users) (tokenString string) {
	if expiredSeconds == 0 {
		expiredSeconds = DEFAULT_EXPIRE_SECONDS
	}
	// Create the Claims
	mySigningKey := []byte(KEY)
	expireAt  := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()
	fmt.Println("token will be expired at ", time.Unix(expireAt, 0) )
	// pass parameter to this func or not
	claims := MyCustomClaims{
		*user,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    user.Name,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("generate json web token failed !! error :", err)
	}
	return tokenStr

}

// return this result to client then all later request should have header "Authorization: Bearer <token> "
func getHeaderTokenValue(tokenString string) string {
	//Authorization: Bearer <token>
	return fmt.Sprintf("Bearer %s", tokenString)
}

var FilterToken = func(ctx *context.Context) {
	logs.Info("current router path is ", ctx.Request.RequestURI)
	if ctx.Request.RequestURI != "/v1/login" && ctx.Input.Header("Authorization") == "" && ctx.Request.RequestURI != "/v1/register"  {
		logs.Error("without token, unauthorized !!")
		ctx.ResponseWriter.WriteHeader(401)
		ctx.ResponseWriter.Write([]byte("no permission"))
	}
	if ctx.Request.RequestURI != "/v1/login" && ctx.Input.Header("Authorization") != "" && ctx.Request.RequestURI != "/v1/register"{
		token := ctx.Input.Header("Authorization")

		err := ValidateToken(token)
		if err != nil {
			ctx.ResponseWriter.WriteHeader(401)
			ctx.ResponseWriter.Write([]byte("token invalid"))
		}
		// invoke ValidateToken in utils/token
		RefreshToken(token)
		// invalid or expired todo res 401
	}
}