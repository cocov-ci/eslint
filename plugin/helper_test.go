package plugin

import (
	"strings"
	"testing"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/cocov-ci/eslint/mocks"
)

type testHelper struct {
	ctx  *mocks.MockContext
	exec *mocks.MockExec
}

func newTestHelper(t *testing.T) *testHelper {
	ctrl := gomock.NewController(t)
	ctx := mocks.NewMockContext(ctrl)
	exec := mocks.NewMockExec(ctrl)

	ctx.EXPECT().L().
		DoAndReturn(func() *zap.Logger { return zap.NewNop() }).
		AnyTimes()

	return &testHelper{ctx, exec}
}

func findRepositoryRoot(t *testing.T) string {
	out, err := cocov.Exec("git", []string{"rev-parse", "--show-toplevel"}, nil)
	require.NoError(t, err)
	return strings.TrimSpace(string(out))
}
