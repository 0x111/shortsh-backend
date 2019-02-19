package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/short-sh/shortsh-backend/models"
	"github.com/short-sh/shortsh-backend/utils"
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
		var protocol string // protocol of the short domain

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

		if shortDomain == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "msg": "There was an error while communicating with the database!"})
		}

		if shortDomain.Secure {
			protocol = "https"
		} else {
			protocol = "http"
		}

		if exists {
			return c.JSON(http.StatusOK, echo.Map{"success": true, "url": protocol + "://" + shortDomain.ShortDomain + "/" + data.ShortId})
		}

		_, err = engine.Insert(urlMeta)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "msg": "There was an error while communicating with the database!"})
		}

		return c.JSON(http.StatusOK, echo.Map{"success": true, "url": protocol + "://" + shortDomain.ShortDomain + "/" + urlMeta.ShortId})
	})

	e.GET("/url/:shortID/stats", func(c echo.Context) error {
		shortID := c.Param("shortID")

		var urlMeta = models.Url{ShortId: shortID}
		has, err := engine.Get(&urlMeta)
		fmt.Println(urlMeta)

		if !has || err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"success": false})
		}

		stats2 := models.UrlStatRet{}
		rows, err := engine.SQL("SELECT url.id,url.created_at,url.url,short_domains.short_domain,short_domains.secure,COUNT(visitors.id) as count FROM url LEFT JOIN visitors ON url.id=visitors.url LEFT JOIN short_domains ON url.short_domain=short_domains.id WHERE short_id=? LIMIT 1;", shortID).Rows(stats2)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "error": "Internal Server Error!"})
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&stats2)

			if err != nil {
				fmt.Println("Error", err)
			}
		}

		// if secure true
		if stats2.Secure == 1 {
			stats2.ShortDomain = "https://" + stats2.ShortDomain
		} else {
			stats2.ShortDomain = "http://" + stats2.ShortDomain
		}

		// stats per day
		singleDay := models.UrlStatDaily{}
		var daily []models.UrlStatDaily
		rows2, err := engine.SQL("SELECT DATE_FORMAT(visitors.created_at, \"%Y-%m-%d\") as day,COUNT(visitors.id) as clicks FROM url LEFT JOIN visitors ON url.id=visitors.url LEFT JOIN short_domains ON url.short_domain=short_domains.id WHERE short_id=? GROUP BY DAY(visitors.created_at);", shortID).Rows(singleDay)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "error": "Internal Server Error!"})
		}

		defer rows2.Close()
		for rows2.Next() {
			err = rows2.Scan(&singleDay)
			fmt.Println(singleDay)
			if err == nil {
				daily = append(daily, singleDay)
			}
		}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "msg": "There was a problem while communicating with the database!"})
		}

		return c.JSON(http.StatusOK, echo.Map{"success": true, "stats": stats2, "daily_stats": daily})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
