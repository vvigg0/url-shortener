package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vvigg0/l3/url-shortener/internal/handler"
	"github.com/vvigg0/l3/url-shortener/internal/myredis"
	repocleaner "github.com/vvigg0/l3/url-shortener/internal/repoCleaner"
	"github.com/vvigg0/l3/url-shortener/internal/repository"
	"github.com/vvigg0/l3/url-shortener/internal/service"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func Run() error {
	dsn := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"))

	URLsTTL, err := time.ParseDuration(os.Getenv("URLS_TTL"))
	if err != nil {
		return err
	}

	cachedURLsTTL, err := time.ParseDuration(os.Getenv("CACHED_URLS_TTL"))
	if err != nil {
		return err
	}

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	client := redis.New("shortener_redis:6379", "", 0)
	rredis := myredis.New(client, cachedURLsTTL)
	defer rredis.Client.Close()
	if err := rredis.Client.Ping(context.Background()); err != nil {
		return err
	}
	db, err := dbpg.New(dsn, nil, &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5})
	if err != nil {
		return fmt.Errorf("ошибка создания объекта БД: %w", err)
	}
	defer db.Master.Close()
	if err := db.Master.Ping(); err != nil {
		return err
	}

	cleaner := repocleaner.New(db)

	h := handler.New(service.New("localhost:8080", repository.New(db, URLsTTL), rredis))
	router := ginext.New("")
	registerRoutes(router, h)

	router.Static("/static", "./web")
	router.LoadHTMLFiles("./web/index.html")

	srvr := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigs)

	errCh := make(chan error, 1)
	go func() {
		zlog.Logger.Info().Msg("сервер запущен")
		if err := srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	go func() {
		if err := cleaner.Start(appCtx); err != nil {
			errCh <- err
			return
		}
		zlog.Logger.Info().Msg("cleaner успешно завершился")
	}()
	select {
	case sig := <-sigs:
		zlog.Logger.Info().Msgf("получен сигнал завершения %v, завершение работы...", sig)
		appCancel()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srvr.Shutdown(ctx); err != nil {
			return err
		}
		zlog.Logger.Info().Msg("server успешно завершился")
	case e := <-errCh:
		return e
	}
	return nil
}

func registerRoutes(r *ginext.Engine, h *handler.Handler) {
	r.GET("/", func(ctx *ginext.Context) {
		ctx.HTML(200, "index.html", nil)
	})
	r.GET("/shortened", h.GetAllLinks)
	r.GET("/s/:short_code", h.ShortLinkRedirect)
	r.GET("/analytics/:short_code", h.GetAnalytics)

	r.POST("/shorten", h.CreateShortLink)
}
