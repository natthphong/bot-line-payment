package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/natthphong/bot-line-payment/config"
	"github.com/natthphong/bot-line-payment/internal/db"
	"github.com/natthphong/bot-line-payment/internal/logz"
	"github.com/natthphong/bot-line-payment/model"
	"github.com/omise/omise-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {

	//Test

	// end test
	ctx := context.Background()
	app := initFiber()
	config.InitTimeZone()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to initial config."))
	}

	logz.Init(cfg.Log.Level, cfg.Server.Name)
	//defer logz.Drop()

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	logger := zap.L()

	jsonCfg, err := json.Marshal(cfg)
	_ = jsonCfg
	logger.Debug("after cfg : " + string(jsonCfg))
	dbPool, err := db.Open(ctx, cfg.DBConfig)
	if err != nil {
		logger.Fatal("server connect to db", zap.Error(err))
	}
	defer dbPool.Close()
	logger.Info("DB CONNECT")

	// tune perfomance with appCache

	baanBot, err := messaging_api.NewMessagingApiAPI(
		cfg.LineConfig["baan-bot"].ChannelToken,
	)
	if err != nil {
		panic(err)
	}

	_ = baanBot

	//userTestId := "U5df59f9c7f4fe1f0cb97a1177324be88"
	//response, err := line_message.SendMessage(baanBot)([]messaging_api.MessageInterface{
	//	messaging_api.TextMessage{
	//		Text: "Hello, world" + time.Now().Format("2006-01-02 15:04:05"),
	//	},
	//}, userTestId, uuid.NewString(), true,
	//)
	//if err != nil {
	//	return
	//}

	logger.Info("line CONNECT")

	client, e := omise.NewClient(cfg.OmiseConfig.PublicKey, cfg.OmiseConfig.SecretKey)
	if e != nil {
		log.Fatal(e)
	}

	_ = client
	logger.Info("omise CONNECT")
	// คิดเงิน
	//source, err := payment.CreateSource(client)(2000, "",
	//	"promptpay",
	//	"thb",
	//	uuid.NewString(),
	//	"Ufa2dce3870921d4d4604bf39024bab7f")
	//if err != nil {
	//	return
	//}
	//fmt.Println(source.ID)
	//charge, err := payment.CreateCharge(client,
	//	[]string{})(2000,
	//	"https://invalid-magic-stations-internal.trycloudflare.com/omise/event/hook",
	//	"thb", time.Now().Add(time.Minute*60), source)
	//if err != nil {
	//	panic(err)
	//	return
	//}
	//fmt.Println("charge:", charge)
	// คิดเงิน

	//app.Post("/line/event/hook", func(c *fiber.Ctx) error {
	//	cb, err := webhook.ParseRequest(os.Getenv("LINE_CHANNEL_SECRET"))
	//	if err != nil {
	//		// Handle any errors that occur.
	//	}
	//	fmt.Println(string(c.BodyRaw()))
	//	return c.Status(http.StatusOK).JSON(fiber.Map{
	//		"status": "ok",
	//	})
	//},
	//)

	app.Post("/omise/event/hook", func(c *fiber.Ctx) error {
		var req model.Event
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		fmt.Println(ToJson(req))
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	},
	)

	if err = app.Listen(fmt.Sprintf(":%v", cfg.Server.Port)); err != nil {
		logger.Fatal(err.Error())
	}

}

func ToJson(v interface{}) string {
	jsonStr, _ := json.Marshal(v)
	return string(jsonStr)
}
func initFiber() *fiber.App {
	app := fiber.New(
		fiber.Config{
			ReadTimeout:           5 * time.Second,
			WriteTimeout:          5 * time.Second,
			IdleTimeout:           30 * time.Second,
			DisableStartupMessage: true,
			CaseSensitive:         true,
			StrictRouting:         true,
		},
	)
	//app.Use(fiber.co)
	app.Use(SetHeaderID())
	return app
}

func SetHeaderID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		randomTrace := uuid.New().String()
		traceId := c.Get("traceId")
		//refId := c.Get("RequestRef")
		if traceId == "" {
			traceId = randomTrace
		}

		c.Accepts(fiber.MIMEApplicationJSON)
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		c.Request().Header.Set("traceId", traceId)
		return c.Next()
	}
}
