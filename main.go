package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl-localization"
	"github.com/calebhiebert/gobbl/context"
	"github.com/calebhiebert/gobbl/luis"
	"github.com/calebhiebert/gobbl/messenger"
	"github.com/calebhiebert/gobbl/session"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// StaticBox is a packr box containing static files
var StaticBox *packr.Box

func main() {
	bundle := i18n.Bundle{
		DefaultLanguage: language.English,
	}
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("./assets/lang/en-US.json")

	localizationConfig := &glocalize.LocalizationConfig{
		Bundle: &bundle,
	}

	StaticBox = packr.New("StaticBox", "assets/static")

	gobblr := gbl.New()

	/*
		STANDARD MIDDLEWARE SETUP
		****************************************
		The MarkSeenMiddleware is specific to facebook messenger
	*/
	gobblr.Use(gbl.ResponderMiddleware())
	gobblr.Use(gbl.UserExtractionMiddleware())
	gobblr.Use(gbl.RequestExtractionMiddleware())
	gobblr.Use(fb.MarkSeenMiddleware())
	gobblr.Use(sess.Middleware(sess.MemoryStore()))
	gobblr.Use(bctx.Middleware())
	gobblr.Use(glocalize.Middleware(localizationConfig))

	/*
		LUIS SETUP
		****************************************
		Uncomment to enable LUIS integration
	*/
	louie, err := luis.New(os.Getenv("LUIS_ENDPOINT"))
	if err != nil {
		panic(err)
	}

	/*
		ROUTER SETUP
		****************************************
		Routers in this project are package-global and used here
	*/
	textRouter := gbl.TextRouter()
	ictxRouter := bctx.ContextIntentRouter()
	intentRouter := gbl.IntentRouter()
	customRouter := gbl.CustomRouter()

	gobblr.Use(textRouter.Middleware())
	gobblr.Use(customRouter.Middleware())

	// LUIS is added at this point so that if any of our text routes match
	// we can skip the NLP process becuase we don't need to know the intent
	gobblr.Use(luis.Middleware(louie))

	gobblr.Use(intentRouter.Middleware())
	gobblr.Use(ictxRouter.Middleware())

	// Fallback handler
	gobblr.Use(HDefaultFallback)

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
		USetMessengerProfile(mapi)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r := gin.Default()

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
