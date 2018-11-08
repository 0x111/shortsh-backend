package main

import (
	"fmt"
	"github.com/0x111/shortsh-backend/models"
	"github.com/0x111/shortsh-backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime"
)

var engine *xorm.Engine

func main() {
	// maximize CPU usage for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error

	// read config
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath("./_config") // optionally look for config in the working directory
	viper.AddConfigPath(".")         // optionally look for config in the working directory
	err = viper.ReadInConfig()       // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	engine, err = xorm.NewEngine("mysql", viper.GetString("mysql_dsn"))

	if err != nil {
		log.Fatalf("We could not connect to the database %v\n", err)
	}

	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	// Sync models
	engine.Sync2(new(models.Url), new(models.ShortDomains), new(models.Visitors))

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{viper.GetString("allow_origins")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, "X-Requested-With", "Cache-Control"},
	}))

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "short.sh engine")
	})

	e.POST("/shorten", func(c echo.Context) error {
		urlMeta := new(models.Url)

		if err := c.Bind(urlMeta); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "msg": "There was an error while processing your request!"})
		}

		parsedURL, err := url.Parse(urlMeta.Url)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "msg": "This does not seem to be a valid URL!"})
		}

		host, _, _ := net.SplitHostPort(parsedURL.Host)

		if len(host) == 0 {
			host = parsedURL.Host
		}

		if host == "a2.to" || host == "short.sh" || len(host) == 0 {
			return c.JSON(http.StatusForbidden, echo.Map{"success": false, "msg": "This does not seem to be a valid URL!"})
		}

		urlMeta.SetRandomID(3)
		urlMeta.SetDomain()

		data, exists := utils.UrlExists(engine, urlMeta.Url)
		shortDomain := utils.GetShortDomain(engine, urlMeta)

		if exists && shortDomain != "" {
			return c.JSON(http.StatusOK, echo.Map{"success": true, "url": "https://" + shortDomain + "/" + data.ShortId})
		}

		_, err = engine.Insert(urlMeta)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "msg": "There was an error while communicating with the database!"})
		}

		return c.JSON(http.StatusOK, echo.Map{"success": true, "url": "https://" + shortDomain + "/" + urlMeta.ShortId})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
