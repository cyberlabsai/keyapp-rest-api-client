package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cyberlabsai/exemplo-integracao-api/settings"

	"github.com/cyberlabsai/exemplo-integracao-api/actions"
	"github.com/cyberlabsai/exemplo-integracao-api/apiclient"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

// Response type to all requests
type Response struct {
	Msg       string    `json:"msg"`
	Timestamp time.Time `json:"timestamp"`
}

var envSettings settings.Settings

func loadSettings() settings.Settings {
	envSettings := settings.Settings{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	debugEnv := os.Getenv("DEBUG")
	if debugEnv == "true" || debugEnv == "TRUE" {
		envSettings.Debug = true
	} else {
		envSettings.Debug = false
	}

	envSettings.Port = os.Getenv("PORT")
	if !envSettings.Debug {
		envSettings.Port = "443"
	} else if envSettings.Port == "" {
		envSettings.Port = "1337"
	}

	envSettings.Domain = os.Getenv("DOMAIN")
	if envSettings.Domain == "" {
		envSettings.Domain = "example.com"
	}

	envSettings.ApiHost = os.Getenv("API_HOST")
	if envSettings.ApiHost == "" {
		log.Fatalln("You didn't inform your api host")
	}

	envSettings.ApiKey = os.Getenv("API_KEY")
	if envSettings.ApiKey == "" {
		log.Fatalln("You didn't inform your api key")
	}

	envSettings.ApiSecret = os.Getenv("API_SECRET")
	if envSettings.ApiSecret == "" {
		log.Fatalln("You didn't inform your api secret")
	}

	return envSettings
}

func main() {
	envSettings := loadSettings()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/signed-actions", signedActions)

	// Start server
	if !envSettings.Debug {

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
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(envSettings.Domain)
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Logger.Fatal(e.StartAutoTLS(":443"))
	}

	e.Logger.Fatal(e.Start(":" + envSettings.Port))
}

func hello(c echo.Context) error {

	resp := &Response{
		Msg:       "Aplicação de Exemplo para Integração com API Pública KeyAPP",
		Timestamp: time.Now(),
	}

	return c.JSON(http.StatusOK, resp)
}

func signedActions(c echo.Context) error {
	actionToken := &actions.ActionToken{}

	if err := c.Bind(actionToken); err != nil {
		return err
	}

	if !actionToken.Valid() {
		return errors.New("Empty token. Please provide a token in the request body")
	}

	cli := &apiclient.Client{
		Settings: settings.Settings{
			ApiHost:   os.Getenv("API_HOST"),
			ApiKey:    os.Getenv("API_KEY"),
			ApiSecret: os.Getenv("API_SECRET"),
		},
		HttpClient: &http.Client{},
	}

	resp, err := cli.VerifyActionToken(actionToken)
	if err != nil {
		return c.JSON(resp.StatusCode, map[string]string{
			"Error": err.Error(),
		})
	}

	body := cli.ResponseBody(resp)
	return c.JSON(resp.StatusCode, body)
}
