package middlewares

import (
	"encoding/base64"
	"errors"
	"fmt"
	"kbaa-fiber-api/helper"
	"kbaa-fiber-api/pkg/functioncaller"
	"kbaa-fiber-api/pkg/logruslogger"
	"kbaa-fiber-api/usecase/base"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type JWTMiddleware struct {
	*base.BaseUc
}

// ferify basic
func (jwtmw JWTMiddleware) VerifyBasic(ctx *fiber.Ctx) (err error) {
	basic := base64.StdEncoding.EncodeToString([]byte(jwtmw.EnvConfig["BASIC_USERNAME"] + ":" + jwtmw.EnvConfig["BASIC_PASSWORD"]))
	fmt.Println("Baseic nya", basic)
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
