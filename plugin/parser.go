package plugin

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Masterminds/semver"
)

type parser struct {
	constraints []string
	buffer      []byte
}

func newParser() *parser {
	return &parser{constraints: []string{}, buffer: []byte{}}
}

func (p *parser) parse(s string) ([]*semver.Constraints, error) {
	for _, b := range s {
		p.feed(byte(b))
	}

	p.finish()

	cs := p.constraints
	res := make(constraints, 0, len(cs))

	for _, c := range cs {
		c = trimXVersion(c)
		constraint, err := semver.NewConstraint(c)
		if err != nil {
			return nil, constraintRenderErr(c, err)
		}
		res = append(res, constraint)
	}

	return res, nil
}

func (p *parser) feed(b byte) {
	r := rune(b)
	switch {
	case unicode.IsSpace(r) || r == ',':
		if len(p.buffer) <= 1 {
			return
		}

		p.constraints = append(p.constraints, string(p.buffer))
		p.buffer = []byte{}
		return

	default:
		p.buffer = append(p.buffer, b)
		return
	}

}

func (p *parser) finish() {
	if len(p.buffer) != 0 {
		p.constraints = append(p.constraints, string(p.buffer))
	}

	p.buffer = []byte{}
}

func constraintRenderErr(constraint string, e error) error {
	return fmt.Errorf("error parsing constraint \"%s\": %s", constraint, e)
}

func trimXVersion(s string) string {
	xVersion := ".x"
	if !strings.Contains(s, xVersion) {
		return s
	}
	return strings.Join(strings.Split(s, xVersion), "")
}
