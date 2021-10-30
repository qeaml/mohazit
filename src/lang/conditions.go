package lang

import (
	"errors"
	"strings"
)

type Conditional struct {
	Comp string
	Args []*Object
}

func ParseConditional(s string, p *Parser) (*Conditional, error) {
	ctx := ""
	params := []string{}
	hasParams := false
	comp := ""
	hasComp := false
	for _, c := range s {
		if hasComp || !hasParams {
			if c == '(' {
				if len(params) > 1 && p.typeOf(params[len(params)-1]) == ObjStr {
					params = append(params, "\\")
				}
				hasParams = true
			} else if c == ' ' {
				a := strings.TrimSpace(ctx)
				if len(a) == 0 {
					continue
				}
				params = append(params, a)
				ctx = ""
			} else {
				ctx += string(c)
			}
		} else {
			if c == ')' {
				comp = strings.ToLower(strings.TrimSpace(comp))
				hasComp = true
			} else {
				comp += string(c)
			}
		}
	}
	if len(strings.TrimSpace(ctx)) != 0 && hasComp {
		params = append(params, strings.TrimSpace(ctx))
	}
	if !hasComp {
		return nil, errors.New("no comparator specified")
	}
	args, err := p.parseArgs(params)
	if err != nil {
		return nil, err
	}
	return &Conditional{Comp: comp, Args: args}, nil
}
