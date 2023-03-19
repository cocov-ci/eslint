package plugin

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	base := ">=0.10.3 <0.12"
	p := newParser()

	cs, err := p.parse(base)
	require.NoError(t, err)
	assert.Len(t, cs, 2)

	c1, c2 := cs[0], cs[1]

	checkVersion, err := semver.NewVersion("0.10.4")

	require.NoError(t, err)
	assert.True(t, c1.Check(checkVersion))
	assert.True(t, c2.Check(checkVersion))

}
