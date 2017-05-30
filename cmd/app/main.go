package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jloup/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Config struct {
	utils.LogConfiguration
	Address   string
	Port      string
	PublicDir string
}

func (c *Config) SetToDefaults() {
	*c = Config{Address: "localhost", Port: "8000", PublicDir: "./public"}
	c.LogConfiguration.SetToDefaults()
}

var C Config
var log = utils.StandardL().WithField("module", "abotllinaire")

func init() {
	utils.MustParseStdConfigFile(&C)
	err := utils.ConfigureStdLogger(C.LogConfiguration)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	log.Infof("CONFIG:\n")
	log.Infof("PUBLIC_DIR: %s\n", C.PublicDir)
	log.Infof("LISTEN: %s\n", strings.Join([]string{C.Address, C.Port}, ":"))
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	e.Static("/public", C.PublicDir)
	e.File("/", path.Join(C.PublicDir, "/index.html"))

	e.Start(strings.Join([]string{C.Address, C.Port}, ":"))
}
