package plugin

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
)

const (
	stateFeed = iota
	stateOr
)

type parser struct {
	state       int
	constraints []string
	buffer      []byte
}

func newParser() *parser {
	return &parser{constraints: []string{}, buffer: []byte{}, state: stateFeed}
}

func (p *parser) parse(s string) (constraints, error) {
	for _, b := range s {
		p.feed(byte(b))
	}

	p.finish()

	return p.render()
}

func (p *parser) render() (constraints, error) {
	cs := make(constraints, 0, len(p.constraints))
	for _, s := range p.constraints {
		s = replaceXVersion(s)
		c, err := semver.NewConstraint(s)
		if err != nil {
			return nil, constraintRenderErr(s, err)
		}
		cs = append(cs, c)
	}
	return cs, nil
}

var symbols = map[rune]bool{
	'=': true,
	'!': true,
	'>': true,
	'<': true,
	'^': true,
	'~': true,
}

func isConstraintSym(r rune) bool { return symbols[r] }

func symBuffer(buff []byte) bool {
	for _, b := range buff {
		if !isConstraintSym(rune(b)) {
			return false
		}
	}
	return true
}

func (p *parser) feed(b byte) {
	r := rune(b)
	switch p.state {
	case stateFeed:
		if unicode.IsSpace(r) {
			return
		}

		if r == ',' {
			p.constraints = append(p.constraints, string(p.buffer))
			p.buffer = p.buffer[:0]
			return
		}

		if r == 'v' {
			if symBuffer(p.buffer) {
				p.buffer = append(p.buffer, b)
				return
			}

			p.constraints = append(p.constraints, string(p.buffer))
			p.buffer = p.buffer[:0]
			return
		}

		if isConstraintSym(r) {
			if len(p.buffer) <= 1 {
				p.buffer = append(p.buffer, b)
				return
			}

			p.constraints = append(p.constraints, string(p.buffer))
			p.buffer = p.buffer[:0]
			p.buffer = append(p.buffer, b)
			return
		}

		if r == '|' {
			p.buffer = append(p.buffer, b)
			p.state = stateOr
			return
		}

	case stateOr:
		if b == ',' {
			p.constraints = append(p.constraints, string(p.buffer))
			p.buffer = p.buffer[:0]
			p.state = stateFeed
			return
		}
	}

	p.buffer = append(p.buffer, b)
}

func (p *parser) finish() {
	if len(p.buffer) != 0 {
		p.constraints = append(p.constraints, string(p.buffer))
	}

	p.buffer = p.buffer[:0]
}

func constraintRenderErr(constraint string, e error) error {
	return fmt.Errorf("error parsing constraint \"%s\": %s", constraint, e)
}

func replaceXVersion(s string) string {
	return strings.Replace(s, "x", "*", -1)
}
