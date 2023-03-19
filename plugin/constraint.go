package plugin

import "github.com/Masterminds/semver"

type constraints []*semver.Constraints

func (c constraints) eval(versions ...*semver.Version) bool {
	var conforms bool
	for _, v := range versions {
		for _, constraint := range c {
			conforms = constraint.Check(v)
		}
	}
	return conforms
}
