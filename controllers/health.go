package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

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

func newHandler() healthcheck.Handler {
	cfg := config.GetConfig()
	rawResourceLoaderUrl := cfg.ResourceLoaderURL
	resolverBridgeUrl := cfg.ResolverBridgeURL
	resourceLoader, _ := url.Parse(rawResourceLoaderUrl)
	finalResourceLoader := resourceLoader.Scheme + "://" + resourceLoader.Host
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))
	// log.Println("resourceLoader: ", resourceLoader)
	health.AddReadinessCheck("remote-resources", Get(finalResourceLoader, 1*time.Second))
	health.AddReadinessCheck("resolver-bridge", Get(resolverBridgeUrl, 1*time.Second))
	return health
}

// Get was the function that allow follow the url
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
