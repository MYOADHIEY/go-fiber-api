package main

import (
	"fmt"
	configuration "kbaa-fiber-api/config"
	"kbaa-fiber-api/pkg/str"
	"kbaa-fiber-api/server/middlewares"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	baseuc "kbaa-fiber-api/usecase/base"

	"kbaa-fiber-api/server/bootsrap"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverMiddleware "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/xid"
)

var (
	validatorDriver *validator.Validate
	uni             *ut.UniversalTranslator
	translator      ut.Translator
	logFormat       = `{"host":"${host}","pid":"${pid}","time":"${time}","req_id":"${locals:requestid}","status":"${status}","method":"${method}","latency":"${latency}","path":"${path}",` +
		`"user_agent":"${ua}","in":"${bytesReceived}",  "out":"${bytesSent}","res_body":"${resBody}", "ip":"${ip}", "query":"${query}"  }`
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")
	loc, err := time.LoadLocation("Asia/Jakarta")
	// handle err
	time.Local = loc
	configurations, err := configuration.LoadConfigurations()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer configurations.DB.Close()
	// KS := 0
	// roe := configurations.DB.QueryRow("select count(*) from _user")
	// err = roe.Scan(&KS)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// init validation driver
	validatorInit(&configurations)

	appVer := configurations.EnvConfig["APP_VERSION"]
	fmt.Println("versinya ", appVer)
	app := fiber.New(fiber.Config{
		BodyLimit:    str.StringToInt(configurations.EnvConfig["FILE_MAX_UPLOAD_SIZE"]),
		ErrorHandler: middlewares.InternalServer,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	BaseUC := baseuc.BaseUc{
		ReqID:     xid.New().String(),
		EnvConfig: configurations.EnvConfig,
		DB:        configurations.DB,
		Validate:  validatorDriver,
		JweCred:   configurations.JweCred,
		JwtCred:   configurations.JwtCred,
	}

	appBoot := bootsrap.AppBoot{
		App:        app,
		BaseUC:     BaseUC,
		Validator:  validatorDriver,
		Translator: translator,
	}

	appBoot.App.Use(requestid.New())
	appBoot.App.Use(logger.New(logger.Config{
		Format:     logFormat + "\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Jakarta",
	}))

	appBoot.App.Use(recoverMiddleware.New())

	app.Use(func(c *fiber.Ctx) error {
		method := c.Method()
		path := c.Path()

		defer func() {
			if r := recover(); r != nil {
				// Ambil hanya stack trace dari project Anda
				stack := string(debug.Stack())
				lines := strings.Split(stack, "\n")

				// Filter: ambil hanya baris yang mengandung package project Anda
				var relevant []string
				for _, line := range lines {
					if strings.Contains(line, "stellar-4.0-backend") {
						relevant = append(relevant, strings.TrimSpace(line))
					}
				}

				log.Printf("🚨 PANIC %s %s : %v\n📍 Relevant stack:\n%s",
					method, path, r, strings.Join(relevant, "\n"))

				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal Server Error",
				})
			}
		}()

		return c.Next()
	})

	appBoot.App.Use(cors.New(cors.Config{
		AllowOrigins:     configurations.EnvConfig["APP_CORS_DOMAIN"],
		AllowMethods:     http.MethodHead + "," + http.MethodGet + "," + http.MethodPost + "," + http.MethodPut + "," + http.MethodPatch + "," + http.MethodDelete,
		AllowHeaders:     "*",
		AllowCredentials: false,
	}))

	appBoot.App.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	// app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	appBoot.RegisterRouters()
	log.Fatal(appBoot.App.Listen(configurations.EnvConfig["APP_HOST"]))

}

func validatorInit(configs *configuration.Configurations) {
	en := en.New()
	id := id.New()
	uni = ut.New(en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	validatorDriver = validator.New()

	enTranslations.RegisterDefaultTranslations(validatorDriver, transEN)
	idTranslations.RegisterDefaultTranslations(validatorDriver, transID)

	switch configs.EnvConfig["APP_LOCALE"] {
	case "en":
		translator = transEN
	case "id":
		translator = transID
	}
}
