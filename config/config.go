package config

import (
	"database/sql"
	"log"
	"time"

	"kbaa-fiber-api/pkg/jwe"
	"kbaa-fiber-api/pkg/jwt"
	postgresqlPkg "kbaa-fiber-api/pkg/postgresql"
	"kbaa-fiber-api/pkg/str"

	"github.com/joho/godotenv"
)

type Configurations struct {
	EnvConfig map[string]string
	DB        *sql.DB
	JweCred   jwe.Credential
	JwtCred   jwt.Credential
}

var (
	envConfigs, _ = godotenv.Read("../.env")
)

func LoadConfigurations() (res Configurations, err error) {
	res.EnvConfig, err = godotenv.Read("../.env")

	dbConn := postgresqlPkg.Connection{
		Host:                    res.EnvConfig["DATABASE_HOST"],
		DbName:                  res.EnvConfig["DATABASE_DB"],
		User:                    res.EnvConfig["DATABASE_USER"],
		Password:                res.EnvConfig["DATABASE_PASSWORD"],
		Port:                    str.StringToInt(res.EnvConfig["DATABASE_PORT"]),
		SslMode:                 res.EnvConfig["DATABASE_SSL_MODE"],
		DBMaxConnection:         str.StringToInt(res.EnvConfig["DATABASE_MAX_CONNECTION"]),
		DBMAxIdleConnection:     str.StringToInt(res.EnvConfig["DATABASE_MAX_IDLE_CONNECTION"]),
		DBMaxLifeTimeConnection: str.StringToInt(res.EnvConfig["DATABASE_MAX_LIFETIME_CONNECTION"]),
	}
	res.DB, err = dbConn.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	res.DB.SetMaxOpenConns(dbConn.DBMaxConnection)
	res.DB.SetMaxIdleConns(dbConn.DBMAxIdleConnection)
	res.DB.SetConnMaxLifetime(time.Duration(dbConn.DBMaxLifeTimeConnection) * time.Second)

	// jwe
	res.JweCred = jwe.Credential{
		KeyLocation: res.EnvConfig["APP_PRIVATE_KEY_LOCATION"],
		Passphrase:  res.EnvConfig["APP_PRIVATE_KEY_PASSPHRASE"],
	}

	// jwt
	res.JwtCred = jwt.Credential{
		Secret:           res.EnvConfig["TOKEN_SECRET"],
		ExpSecret:        str.StringToInt(res.EnvConfig["TOKEN_EXP_SECRET"]),
		RefreshSecret:    res.EnvConfig["TOKEN_REFRESH_SECRET"],
		RefreshExpSecret: str.StringToInt(res.EnvConfig["TOKEN_EXP_REFRESH_SECRET"]),
	}

	return
}
