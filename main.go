package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl-localization"
	"github.com/calebhiebert/gobbl-redis-store"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/messenger"
	"github.com/calebhiebert/gobbl/session"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gobuffalo/packr/v2"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gobbl-starter-bot/bdb"
)

// StaticBox is a packr box containing static files
var StaticBox *packr.Box

// FileBox is a packr box containing static files
var FileBox *packr.Box

func main() {
	godotenv.Load(".env")

	StaticBox = packr.New("StaticBox", "assets/static")
	FileBox = packr.New("FileBox", "assets")

	dbase, err := bdb.New()
	if err != nil {
		panic(err)
	}

	sql, err := FileBox.FindString("schema.sql")
	if err != nil {
		panic(err)
	}

	_, err = dbase.Exec(sql)
	if err != nil {
		panic(err)
	}

	bundle := i18n.Bundle{
		DefaultLanguage: language.English,
	}
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	/*
		LOAD LANGUAGE FILES
		****************************************
		This is where language files are loaded. To add more files,
		just add more lines like this one. If you wanted to add lang/fr-CA.json:
		bundle.MustParseMessageFileBytes(FileBox.Bytes("lang/fr-CA.json"), "fr-CA.json")
	*/
	langENUS, err := FileBox.Find("lang/en-US.json")
	if err != nil {
		panic(err)
	}

	bundle.MustParseMessageFileBytes(langENUS, "en-US.json")

	localizationConfig := &glocalize.LocalizationConfig{
		Bundle: &bundle,
	}

	redisStore := gobblredis.New(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}, 20*time.Minute, "session:")

	gobblr := gbl.New()

	/*
		STANDARD MIDDLEWARE SETUP
		****************************************
		The MarkSeenMiddleware is specific to facebook messenger
	*/
	gobblr.Use(bdb.Middleware(dbase))
	gobblr.Use(gbl.ResponderMiddleware())
	gobblr.Use(gbl.UserExtractionMiddleware())
	gobblr.Use(gbl.RequestExtractionMiddleware())
	gobblr.Use(fb.MarkSeenMiddleware())
	gobblr.Use(sess.Middleware(redisStore))
	gobblr.Use(bctx.Middleware())
	gobblr.Use(glocalize.Middleware(localizationConfig))
	gobblr.Use(UUserExtractorMiddleware())
	gobblr.Use(UUserActionLoggerMiddleware())

	/*
		ROUTER SETUP
		****************************************
		Defined routers are here. This setup should work for most projects,
		but new routers can easily be added.
	*/
	textRouter := gbl.TextRouter()
	ictxRouter := bctx.ContextIntentRouter()
	intentRouter := gbl.IntentRouter()
	customRouter := gbl.CustomRouter()

	gobblr.Use(textRouter.Middleware())
	gobblr.Use(customRouter.Middleware())
	gobblr.Use(intentRouter.Middleware())
	gobblr.Use(ictxRouter.Middleware())

	// Fallback handler
	gobblr.Use(HDefaultFallback)

	/*
		DEV TOOLS
		****************************************
		These routes are only enabled when the ENABLE_DEV_TOOLS environment
		variable is set to "true"
	*/
	if os.Getenv("ENABLE_DEV_TOOLS") == "true" {
		textRouter.Text(TCID, HGetID)
		textRouter.Text(TCDeleteData, HDeleteUserData)
	}

	/*
		ROUTE SETUP
		****************************************
		All the project routes are defined here
	*/

	// Text Routes
	textRouter.Text("GET_STARTED", HGetStarted)

	/*
		FACEBOOK MESSENGER SETUP
		****************************************
	*/
	mapi := fb.CreateMessengerAPI(os.Getenv("PAGE_ACCESS_TOKEN"))
	messengerIntegration := fb.MessengerIntegration{
		API:         mapi,
		Bot:         gobblr,
		VerifyToken: "frogs",
		DevMode:     true,
	}

	if os.Getenv("MESSENGER_PROFILE") == "true" {
		res, err := USetMessengerProfile(mapi)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Messenger Profile", res)
		}
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	/*
		WORKER CHECK
		****************************************
		To enable "worker" mode (for long running tasks), just start the
		bot with the argument "worker"
	*/
	if len(os.Args) >= 2 && os.Args[1] == "worker" {
		// TODO add some lovely worker activities
		fmt.Println("Doing Work")
	} else {
		r := gin.Default()

		r.Use(bdb.GinMiddleware(dbase))

		r.GET("/webhook", func(c *gin.Context) {
			mode := c.Query("hub.mode")
			token := c.Query("hub.verify_token")
			challenge := c.Query("hub.challenge")

			if mode == "subscribe" && token == messengerIntegration.VerifyToken {
				c.String(200, challenge)
			} else {
				c.AbortWithStatus(401)
			}
		})

		r.POST("/webhook", func(c *gin.Context) {
			var webhookRequest fb.WebhookRequest

			err := c.ShouldBindJSON(&webhookRequest)
			if err != nil {
				fmt.Println("WEBHOOK PARSE ERR", err)
				c.JSON(500, gin.H{"error": "Invalid json"})
			} else {
				_ = messengerIntegration.ProcessWebhookRequest(&webhookRequest)

				c.JSON(200, gin.H{"o": "k"})
			}
		})

		r.StaticFS("/static", StaticBox)

		r.Run("0.0.0.0:" + port)
	}
}
