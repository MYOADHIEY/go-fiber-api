package utils

import (
	"context"
	"errors"
	"kbaa-fiber-api/helper"
	"kbaa-fiber-api/pkg/logruslogger"
	baseUc "kbaa-fiber-api/usecase/base"
	iviewmodels "kbaa-fiber-api/usecase/utils/viewmodels"

	"github.com/rs/xid"
)

type JWTUC struct {
	*baseUc.BaseUc
}

// GenerateToken ...
func (uc JWTUC) GenerateJWTToken(c context.Context, payload map[string]interface{}, res *iviewmodels.JWTVM) (err error) {
	ctx := "JwtUC.GenerateToken"

	deviceID := xid.New().String()
	payload["device_id"] = deviceID

	jwePayload, err := uc.BaseUc.JweCred.Generate(payload)

	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwe", c.Value("requestid"))
		return errors.New(helper.JWT)
	}
	res.Token, res.ExpiredDate, err = uc.BaseUc.JwtCred.GetToken(jwePayload)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwt", c.Value("requestid"))
		return errors.New(helper.JWT)
	}
	res.RefreshToken, res.RefreshExpiredDate, err = uc.BaseUc.JwtCred.GetRefreshToken(jwePayload)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "refresh_jwt", c.Value("requestid"))
		return errors.New(helper.JWT)
	}

	return err
}
