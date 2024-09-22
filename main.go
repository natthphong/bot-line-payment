package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/natthphong/bot-line-payment/api"
	"github.com/natthphong/bot-line-payment/config"
	"github.com/natthphong/bot-line-payment/handler/auth"
	"github.com/natthphong/bot-line-payment/handler/user"
	"github.com/natthphong/bot-line-payment/internal/db"
	"github.com/natthphong/bot-line-payment/internal/httputil"
	"github.com/natthphong/bot-line-payment/internal/logz"
	"github.com/natthphong/bot-line-payment/internal/s3util"
	"github.com/natthphong/bot-line-payment/middleware"
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
	httpClient := httputil.InitHttpClient(
		cfg.HTTP.TimeOut,
		cfg.HTTP.MaxIdleConn,
		cfg.HTTP.MaxIdleConnPerHost,
		cfg.HTTP.MaxConnPerHost,
	)
	_ = httpClient
	svc, err := s3util.OpenS3(cfg.AwsS3Config)
	if err != nil {
		panic(err)
	}
	_ = svc
	logger.Info("S3 CONNECT")
	baanBot, err := messaging_api.NewMessagingApiAPI(
		cfg.LineConfig["baan-bot"].ChannelToken,
	)
	if err != nil {
		panic(err)
	}

	_ = baanBot

	logger.Info("line CONNECT")

	//client, e := omise.NewClient(cfg.OmiseConfig.PublicKey, cfg.OmiseConfig.SecretKey)
	//if e != nil {
	//	log.Fatal(e)
	//}
	//
	//_ = client
	//logger.Info("omise CONNECT")
	//// คิดเงิน
	//////TODO insert txn
	//source, err := payment.CreateSource(client)(2000, "IOS",
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
	////fmt.Println("charge:", charge)
	////url := "https://api.omise.co/charges/chrg_test_6161p4jhaa9civutsl5/documents/docu_test_6161p4lyiyf3izc7x9y/downloads/6384C3AE0D9127A4"
	//url := charge.Source.ScannableCode.Image.DownloadURI
	//fmt.Println(url)

	//file, err := httputil.DownloadFile(httpClient, url)
	//if err != nil {
	//	panic(err)
	//}
	//keyFile := fmt.Sprintf("%s/%s", *charge.Source.ScannableCode.Image.Location, charge.Source.ScannableCode.Image.Filename)
	//exp := time.Now().Add(time.Minute * 1)
	//_, err = svc.PutObject(&s3.PutObjectInput{
	//	Expires:     &exp,
	//	Bucket:      aws.String(cfg.AwsS3Config.BucketName),
	//	Key:         aws.String(keyFile),
	//	Body:        bytes.NewReader(file), // Convert byte array to io.Reader
	//	ContentType: aws.String("image/svg+xml"),
	//	ACL:         aws.String("private"), // Set the file access permissions (private, public-read, etc.)
	//})
	//if err != nil {
	//	panic(err)
	//}
	//request, _ := svc.GetObjectRequest(&s3.GetObjectInput{
	//	Bucket: aws.String(cfg.AwsS3Config.BucketName),
	//	Key:    aws.String(keyFile),
	//})
	//
	//urlPresign, err := request.Presign(10 * time.Minute)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(urlPresign)
	////คิดเงิน

	//userTestId := "U5df59f9c7f4fe1f0cb97a1177324be88"
	//
	//_, err = line_message.SendMessage(baanBot)([]messaging_api.MessageInterface{
	//	messaging_api.ImageCarouselTemplate{},
	//	messaging_api.ImageMessage{
	//		OriginalContentUrl: "https://fastly.picsum.photos/id/237/200/300.jpg?hmac=TmmQSbShHz9CdQm0NkEjx1Dyh_Y984R9LpNrpvH2D_U",
	//		PreviewImageUrl:    "https://fastly.picsum.photos/id/237/200/300.jpg?hmac=TmmQSbShHz9CdQm0NkEjx1Dyh_Y984R9LpNrpvH2D_U",
	//	},
	//}, userTestId, uuid.NewString(), true,
	//)
	//if err != nil {
	//	return
	//}

	// TODO end test
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
		// insert finish txn
		var req omise.Charge
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

	group := app.Group("/api/v1")
	group.Post("/login/line", auth.LoginHandler(
		cfg.LineLoginClientId,
		auth.InsertAndMergeUserLineLogin(dbPool),
	))
	groupAuth := group.Group("/auth")
	groupAuth.Get("/me", func(ctx *fiber.Ctx) error {
		body, err := middleware.JwtTokenGetUser(dbPool)(ctx, logger)
		if err != nil {
			return api.JwtError(ctx, err.Error())
		}
		return api.Ok(ctx, body)
	})
	groupAuth.Post("/branch/list", user.BranchListHandler(user.GetListBranches(dbPool)))
	groupAuth.Post("/category/list", user.CategoryListHandler())
	groupAuth.Post("/product/list", user.ProductListHandler())

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
	app.Use(cors.New())
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
