package main

import (
	"bytes"
	"errors"
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
	debug     bool
	port      string
	domain    string
	apiHost   string
	apiKey    string
	apiSecret string
}

// Response type to all requests
type Response struct {
	Msg       string    `json:"msg"`
	Timestamp time.Time `json:"timestamp"`
}

// ActionToken holds the incoming token
type ActionToken struct {
	Token string `json:"token"`
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

	settings.apiHost = os.Getenv("API_HOST")
	if settings.apiHost == "" {
		log.Fatalln("You didn't inform your api host")
	}

	settings.apiKey = os.Getenv("API_KEY")
	if settings.apiKey == "" {
		log.Fatalln("You didn't inform your api key")
	}

	settings.apiSecret = os.Getenv("API_SECRET")
	if settings.apiSecret == "" {
		log.Fatalln("You didn't inform your api secret")
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
	e.POST("/signed-actions", signedActions)

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

func signedActions(c echo.Context) error {
	actionToken := &ActionToken{}
	c.Bind(actionToken)

	if actionToken.Token == "" {
		return errors.New("Empty token. Please provide a token in the request body")
	}

	request, _ := http.NewRequest("GET", settings.apiHost+"/apptokens", nil)
	request.Header.Add("authorization", "Basic "+settings.apiKey+":"+settings.apiSecret)
	request.Header.Add("token", actionToken.Token)

	cli := &http.Client{}
	resp, err := cli.Do(request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	responseBody := buf.String()

	if err != nil {
		return err
	}

	return c.JSON(resp.StatusCode, responseBody)
}
