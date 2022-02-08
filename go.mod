module github.com/bancodobrasil/featws-ruller

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/hyperjumptech/grule-rule-engine v1.10.4
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/hyperjumptech/grule-rule-engine => github.com/bancodobrasil/grule-rule-engine v1.10.5-0.20220131123643-3eab2150780b
