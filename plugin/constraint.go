package plugin

import "github.com/Masterminds/semver"

type constraints []*semver.Constraints

func (c constraints) eval(versions ...*semver.Version) bool {
	for _, v := range versions {
		for _, constraint := range c {
			if constraint.Check(v) {
				return true
			}
		}
	}
	return false
}
