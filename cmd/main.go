package main

import (
	"github.com/vvigg0/l3/url-shortener/cmd/app"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.Init()
	if err := app.Run(); err != nil {
		zlog.Logger.Fatal().Msgf("ошибка запуска: %v", err)
	}
}
