package main

import (
	"context"
	"encoding/json"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/natthphong/bot-line-payment/config"
	"github.com/natthphong/bot-line-payment/internal/db"
	"github.com/natthphong/bot-line-payment/internal/logz"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
)

func main() {
	ctx := context.Background()

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

	botCat, err := messaging_api.NewMessagingApiAPI(
		cfg.LineConfig["bot-tar-cat"].ChannelToken,
	)
	if err != nil {
		panic(err)
	}
	botDog, err := messaging_api.NewMessagingApiAPI(
		cfg.LineConfig["bot-tar-cat"].ChannelToken,
	)
	if err != nil {
		panic(err)
	}
	_ = botCat
	_ = botDog
	req := &messaging_api.BroadcastRequest{
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: "Hello, world",
			},
		},
		NotificationDisabled: true,
	}
	_, err = botCat.Broadcast(req, "")
	if err != nil {
		return
	}

	logger.Info("line CONNECT")

}
