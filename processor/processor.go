package processor

import (
	"strconv"

	"github.com/bancodobrasil/featws-ruller/types"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Evaluate(ctx *types.Context, expression string) string {
	return "true"
}

func (p *Processor) Boolean(value bool) string {
	return strconv.FormatBool(value)
}
