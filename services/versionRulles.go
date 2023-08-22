package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/featws-ruller/types"
)

var (
	VersionRulesService     IVersionRules = NewVersionRules(config.GetConfig())
	URLResouceTimeoutSecond               = 10
)

type Pipeline struct {
	Id        int64  `json:"id"`
	Sha       string `json:"sha"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type PackageVersion struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	CreatedAt string `json:"created_at"`
	Pipeline  `json:"pipeline"`
	Pipelines []Pipeline `json:"pipelines"`
}

type IVersionRules interface {
	GetVersionsFromRuleSheet(string) ([]PackageVersion, error)
	GetLatestVersionFromRulesheet(knowledge string) (string, error)
}

type VersionRules struct {
	resourceUrl     string
	resourceHeaders http.Header
	cacheVersions   *types.Cache[string, PackageVersion]
	versionTTL      int64
}

func NewVersionRules(config *config.Config) VersionRules {
	rawResourceUrl, err := getPackageRegistryUrl(config.ResourceLoaderURL)
	if err != nil {
		panic("Resource loader URL is invalid!")
	}
	return VersionRules{
		resourceUrl:     rawResourceUrl,
		resourceHeaders: config.ResourceLoaderHeaders,
		cacheVersions:   types.NewCache[string, PackageVersion](),
		versionTTL:      config.KnowledgeBaseVersionTTL,
	}
}

func (v VersionRules) GetLatestVersionFromRulesheet(knowledgeBase string) (string, error) {
	if cachedData, ok := v.cacheVersions.Get(knowledgeBase); ok {
		return cachedData.Version, nil
	}

	packageVersions, err := v.GetVersionsFromRuleSheet(knowledgeBase)
	if err != nil {
		return "", fmt.Errorf("Error on getting the latested version: %s", err)
	}
	highestVersion := -1
	pvMap := make(map[string]PackageVersion)
	for _, pv := range packageVersions {
		pvMap[pv.Version] = pv
		pvVersion, err := strconv.Atoi(pv.Version)
		if err != nil {
			continue
		}
		if pvVersion > highestVersion {
			highestVersion = pvVersion
		}
	}
	highestPv := pvMap[strconv.Itoa(highestVersion)]
	v.cacheVersions.Set(knowledgeBase, highestPv, time.Duration(v.versionTTL)*time.Second)
	return highestPv.Version, nil
}

func (v VersionRules) GetVersionsFromRuleSheet(knowledBase string) ([]PackageVersion, error) {
	rawPackageResourceUrl := strings.Replace(v.resourceUrl, "{knowledgeBase}", knowledBase, 1)
	packageResourceUrl, err := url.Parse(rawPackageResourceUrl)
	if err != nil {
		return nil, err
	}

	queryString := packageResourceUrl.Query()
	queryString.Set("package_type", "generic")
	packageResourceUrl.RawQuery = queryString.Encode()

	client := types.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(URLResouceTimeoutSecond)*time.Second)
	defer cancel()
	log.Debugf("Requesting packages resources from %s", packageResourceUrl.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, packageResourceUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	if len(v.resourceHeaders) > 0 {
		req.Header = v.resourceHeaders
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed with status: %s", resp.Status)
	}

	var packageVersions []PackageVersion
	err = json.NewDecoder(resp.Body).Decode(&packageVersions)
	if err != nil {
		return nil, fmt.Errorf("Error encoding JSON: %s", err)
	}
	return packageVersions, nil
}

func getPackageRegistryUrl(resourceUrl string) (string, error) {
	pattern := `https.*?packages`
	re := regexp.MustCompile(pattern)
	match := re.FindString(resourceUrl)

	if match != "" {
		return match, nil
	} else {
		return "", fmt.Errorf("ResoucerUrl '%s' is incompatible", resourceUrl)
	}
}
