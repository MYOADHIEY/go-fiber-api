package middlewares

import (
	"encoding/base64"
	"errors"
	"fmt"
	"kbaa-fiber-api/helper"
	"kbaa-fiber-api/pkg/functioncaller"
	"kbaa-fiber-api/pkg/interfacepkg"
	"kbaa-fiber-api/pkg/logruslogger"
	ihanldler "kbaa-fiber-api/server/handlers/base"
	"kbaa-fiber-api/usecase/base"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type AuthorizationMiddleware struct {
	*base.BaseUc
}

// ferify basic
func (authmiddleware AuthorizationMiddleware) VerifyBasic(ctx *fiber.Ctx) (err error) {
	basic := base64.StdEncoding.EncodeToString([]byte(authmiddleware.EnvConfig["BASIC_USERNAME"] + ":" + authmiddleware.EnvConfig["BASIC_PASSWORD"]))

	header := ctx.Get("Authorization")
	if !strings.Contains(header, "Basic") {
		logruslogger.Log(logruslogger.WarnLevel, helper.HeaderNotPresent, functioncaller.PrintFuncName(), "middleware-jwt-header")
		return errors.New(helper.HeaderNotPresent)
	}
	token := strings.Replace(header, "Basic ", "", -1)
	if token != basic {
		logruslogger.Log(logruslogger.WarnLevel, basic, functioncaller.PrintFuncName(), "invalid-token")
		return errors.New(helper.UnexpectedClaims)
	}

	return ctx.Next()
}

func (authmiddleware AuthorizationMiddleware) verifyBearer(ctx *fiber.Ctx, role string) (res map[string]interface{}, err error) {
	claims := &jwt.StandardClaims{}

	header := ctx.Get("Authorization")
	if !strings.Contains(header, "Bearer") {
		logruslogger.Log(logruslogger.WarnLevel, helper.HeaderNotPresent, functioncaller.PrintFuncName(), "middleware-jwt-header")
		return res, errors.New(helper.HeaderNotPresent)
	}

	//check claims and signing method
	token := strings.Replace(header, "Bearer ", "", -1)
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			logruslogger.Log(logruslogger.WarnLevel, helper.UnexpectedSigningMethod, functioncaller.PrintFuncName(), "middleware-jwt-checkSigningMethod")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := authmiddleware.EnvConfig["TOKEN_SECRET"]
		return []byte(secret), nil
	})

	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "middleware-jwt-checkClaims")
		return res, errors.New(helper.UnexpectedClaims)
	}

	//check token live time
	if claims.ExpiresAt < time.Now().Unix() {
		logruslogger.Log(logruslogger.WarnLevel, helper.ExpiredToken, functioncaller.PrintFuncName(), "middleware-jwt-checkTokenLiveTime")
		return res, errors.New(helper.ExpiredToken)
	}

	//jwe roll back encrypted id
	res, err = authmiddleware.JweCred.Rollback(claims.Id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, helper.Unauthorized, functioncaller.PrintFuncName(), "pkg-jwe-rollback")
		return res, errors.New(helper.Unauthorized)
	}
	if res == nil {
		logruslogger.Log(logruslogger.WarnLevel, helper.Unauthorized, functioncaller.PrintFuncName(), "pkg-jwe-resultNil")
		return res, errors.New(helper.Unauthorized)
	}

	if role != "" && fmt.Sprintf("%v", res["role"]) != role {
		logruslogger.Log(logruslogger.WarnLevel, helper.InvalidRole, functioncaller.PrintFuncName(), "pkg-jwe-resultNil")
		return res, errors.New(helper.InvalidRole)
	}

	logruslogger.Log(logruslogger.InfoLevel, interfacepkg.Marshal(res), functioncaller.PrintFuncName(), "user", ctx.Locals("requestid"))

	return res, nil
}

// VerifyUser jwt middleware
func (authmiddleware AuthorizationMiddleware) VerifyJWT(ctx *fiber.Ctx) (err error) {
	handler := ihanldler.BaseHandler{BaseUC: authmiddleware.BaseUc}

	jweRes, err := authmiddleware.verifyBearer(ctx, "")
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "verify")
		return handler.SendResponse(ctx, nil, nil, err.Error(), http.StatusUnauthorized)
	}
	// set id to uce case contract
	ctx.Locals("user_id", jweRes["user_id"].(string))

	return ctx.Next()
}
