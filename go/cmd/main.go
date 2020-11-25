package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

// Settings for the whole project
type Settings struct {
	debug  bool
	port   string
	domain string
	apiKey string
}

// Response type to all requests
type Response struct {
	Msg       string    `json:"msg"`
	Timestamp time.Time `json:"timestamp"`
}

var settings Settings

func loadSettings() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	debugEnv := os.Getenv("DEBUG")
	if debugEnv == "true" || debugEnv == "TRUE" {
		settings.debug = true
	} else {
		settings.debug = false
	}

	settings.port = os.Getenv("PORT")
	if !settings.debug {
		settings.port = "443"
	} else if settings.port == "" {
		settings.port = "1337"
	}

	settings.domain = os.Getenv("DOMAIN")
	if settings.domain == "" {
		settings.domain = "example.com"
	}

	settings.apiKey = os.Getenv("API_KEY")
	if settings.apiKey == "" {
		log.Fatalln("You didn't inform your api key")
	}
}

func main() {

	loadSettings()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/secure-actions", secureActions)

	// Start server
	if !settings.debug {

		// Listen for HTTP requests on port 80 in a new goroutine. Use
		// autocertManager.HTTPHandler(nil) as the handler. This will send ACME
		// "http-01" challenge responses as necessary, and 302 redirect all other
		// requests to HTTPS.
		go func() {
			srv := &http.Server{
				Addr:         ":80",
				Handler:      e.AutoTLSManager.HTTPHandler(nil),
				IdleTimeout:  time.Minute,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			err := srv.ListenAndServe()
			e.Logger.Fatal(err)
		}()

		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.Prompt = autocert.AcceptTOS
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(settings.domain)
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Logger.Fatal(e.StartAutoTLS(":443"))
	}
	e.Logger.Fatal(e.Start(":" + settings.port))
}

func hello(c echo.Context) error {

	resp := &Response{
		Msg:       "Aplicação de Exemplo para Integração com API Pública KeyAPP",
		Timestamp: time.Now(),
	}

	return c.JSON(http.StatusOK, resp)
}

func secureActions(c echo.Context) error {

	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Response{Msg: "Bad Request", Timestamp: time.Now()})
	}

	//token, err := getToken()

	// TODO: ask to keyapp's public api with the token received in the header is valid

	return c.JSON(http.StatusCreated, jsonMap)

}
