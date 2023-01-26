package golang_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/goci/lib/golang"
)

func TestLint(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	_, err := golang.Lint(c, golang.ProjectRoot("./testdata/test")).ExitCode(ctx)
	require.NoError(t, err)
}
