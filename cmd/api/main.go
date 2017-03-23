package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jloup/abotllinaire/app/api"
	"github.com/jloup/abotllinaire/app/db"
	"github.com/jloup/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var log = utils.StandardL().WithField("module", "sfapi")

type Config struct {
	utils.LogConfiguration
	Address   string
	Port      string
	ModelPath string
	TorchPath string
	Dir       string
	NbWorkers int

	api.FacebookConfig
	api.BotConfig
}

func (c *Config) SetToDefaults() {
	*c = Config{Address: "localhost",
		Port:      "8000",
		ModelPath: "/home/jam/lab/abotllinaire/lm_lstm_epoch24.91_1.1464.t7",
		Dir:       "/home/jam/lab/char-rnn",
		TorchPath: "/home/jam/torch/install/bin/th",
		NbWorkers: 2}
	c.LogConfiguration.SetToDefaults()
}

var C Config

func init() {
	utils.MustParseStdConfigFile(&C)
	err := utils.ConfigureStdLogger(C.LogConfiguration)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := api.InitWorkerPool(C.NbWorkers, C.TorchPath, C.Dir, C.ModelPath); err != nil {
		log.Errorf("cannot init worker pool: %v", err)
		os.Exit(-1)
	}

	api.SetFacebookCredentials(C.FacebookConfig)
	api.SetBotParameters(C.BotConfig)

	db.InitDB()

	log.Infof("CONFIG:")
	log.Infof("LISTEN: %s", strings.Join([]string{C.Address, C.Port}, ":"))
	log.Infof("MODEL PATH: '%s'", C.ModelPath)
	log.Infof("TORCH DIR: '%s'", C.Dir)
	log.Infof("TORCH PATH: '%s'", C.TorchPath)
	log.Infof("FB WEBHOOK TOKEN: '%s'", C.FacebookConfig.WebhookToken)
	log.Infof("FB APP SECRET: '%s'", C.FacebookConfig.AppSecret)
	log.Infof("BOT CONFIG: temp=[%v, %v] len=%v", C.BotConfig.TemperatureMin, C.BotConfig.TemperatureMax, C.BotConfig.PoemLen)
}

type EchoRouter interface {
	Use(...echo.MiddlewareFunc)
	GET(string, echo.HandlerFunc, ...echo.MiddlewareFunc)
	POST(string, echo.HandlerFunc, ...echo.MiddlewareFunc)
}

func SetGroup(e *echo.Echo, group api.Group) {
	var r EchoRouter
	if group.Root == "" {
		r = e
	} else {
		r = e.Group(group.Root)
	}

	if group.Middlewares != nil {
		r.Use(group.Middlewares...)
	}

	for _, route := range group.Routes {
		SetApiRoute(r, route)
	}
}

func SetApiRoute(e EchoRouter, route api.Route) {
	switch route.Method {
	case echo.GET:
		e.GET(route.Path, route.Handler, route.Middlewares...)
	case echo.POST:
		e.POST(route.Path, route.Handler, route.Middlewares...)
	default:
		panic(fmt.Errorf("%s method not handled", route.Method))
	}
}

func main() {
	e := echo.New()

	CORSConfig := middleware.CORSConfig{
		AllowOrigins:     []string{"https://jlj.am"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Authorization"},
	}

	e.Use(middleware.CORSWithConfig(CORSConfig))
	//e.Use(middleware.Logger())

	for _, group := range api.Groups {
		SetGroup(e, group)
	}

	e.Start(strings.Join([]string{C.Address, C.Port}, ":"))
}
