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

// NewHealthController ...
func NewHealthController() *HealthController {
	return &HealthController{
		health: newHandler(),
	}
}

var health = healthcheck.NewHandler()

func newHandler() healthcheck.Handler {
	cfg := config.GetConfig()
	health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))

	if cfg.ResourceLoaderURL != "" {
		rawResourceLoaderURL := cfg.ResourceLoaderURL
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

// HealthLiveHandler ...
func (c *HealthController) HealthLiveHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.LiveEndpoint))
}

// HealthReadyHandler ...
func (c *HealthController) HealthReadyHandler() gin.HandlerFunc {
	return gin.WrapH(http.HandlerFunc(c.health.ReadyEndpoint))
}
