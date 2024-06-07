package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/healthcheck"
	"github.com/bancodobrasil/healthcheck/checks/goroutine"
	"github.com/gin-gonic/gin"
	"github.com/gsdenys/healthcheck/checks"
)

// HealthController the health endpoints controller
type HealthController struct {
	health healthcheck.Handler
}

// NewHealthController returns a new instance of the HealthController struct with a newHandler.
func NewHealthController() *HealthController {
	return &HealthController{
		health: newHandler(),
	}
}

// This instance will be used to register liveness and readiness checks for the application's
// health endpoints.
var health = healthcheck.NewHandler()

// newHandler creates a new healthcheck handler and adds liveness and readiness checks based on the
// configuration.
func newHandler() healthcheck.Handler {
	cfg := config.GetConfig()
	health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))

	if cfg.ResourceLoader.Type == "http" && cfg.ResourceLoader.HTTP.URL != "" {
		rawResourceLoaderURL := cfg.ResourceLoader.HTTP.URL
		resourceLoader, _ := url.Parse(rawResourceLoaderURL)

		if resourceLoader.Scheme == "" {
			log.Fatal("ResourceLoaderURL must have a scheme: http:// or https://")
		}

		if resourceLoader.Host == "" {
			log.Fatal("ResourceLoaderURL must have a host: example.com")
		}
		finalResourceLoader := resourceLoader.Scheme + "://" + resourceLoader.Host
		health.AddReadinessCheck("resource-loader", Get(finalResourceLoader, 1*time.Second))
	}

	if cfg.ResolverBridgeURL != "" {
		resolverBridgeURL := cfg.ResolverBridgeURL
		health.AddReadinessCheck("resolver-bridge", Get(resolverBridgeURL, 1*time.Second))
	}

	return health
}

// Get returns a check function that performs an HTTP GET request to a specified URL with a
// timeout and returns an error if the response status code isn't 200.
func Get(url string, timeout time.Duration) checks.Check {
	client := http.Client{
		Timeout: timeout,
	}

	return func() error {
		resp, err := client.Get(url)

		if err != nil {
			return err
		}

		resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("returned status %d", resp.StatusCode)
		}

		return nil
	}
}

// HealthLiveHandler is a Gin HTTP handler function that wraps the LiveEndpoint
// method of the health instance of the HealthController struct. The LiveEndpoint
// method is a handler function that returns a 200 status code if the application is live.
func (c *HealthController) HealthLiveHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.LiveEndpoint))
}

// HealthReadyHandler returns a Gin HTTP handler function that wraps the ReadyEndpoint method
// of the health instance in the HealthController struct. The ReadyEndpoint method is a handler
// function that returns a 200 status code if the application is ready to receive traffic. The gin.WrapH
// function is used to convert the http.HandlerFunc returned by the ReadyEndpoint method into a gin.HandlerFunc
func (c *HealthController) HealthReadyHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.ReadyEndpoint))
}
