package golang_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/goci/lib/golang"
)

func setupClient(t *testing.T) (context.Context, *dagger.Client, *bytes.Buffer) {
	t.Helper()

	ctx := context.Background()

	var buf bytes.Buffer

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(&buf))
	require.NoError(t, err)

	t.Cleanup(func() {
		c.Close()
	})

	return ctx, c, &buf
}

func TestBase_Default(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	_, err := golang.Base(c).WithExec([]string{"go", "version"}).ExitCode(ctx)
	require.NoError(t, err)
}

func TestBase_BaseImageRepository(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	_, err := golang.Base(c, golang.BaseImageRepository("index.docker.io/library/golang")).WithExec([]string{"go", "version"}).ExitCode(ctx)
	require.NoError(t, err)
}

func TestBase_BaseImageTag(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	platform, err := c.DefaultPlatform(ctx)
	require.NoError(t, err)

	const version = "1.19.5"

	output, err := golang.Base(c, golang.BaseImageTag("1.19.5")).WithExec([]string{"go", "version"}).Stdout(ctx)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("go version go%s %s", version, platform), strings.TrimSpace(output))
}

func TestBase_Version(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	platform, err := c.DefaultPlatform(ctx)
	require.NoError(t, err)

	const version = "1.19.5"

	output, err := golang.Base(c, golang.Version("1.19.5")).WithExec([]string{"go", "version"}).Stdout(ctx)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("go version go%s %s", version, platform), strings.TrimSpace(output))
}

func TestBase_BaseImage(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	platform, err := c.DefaultPlatform(ctx)
	require.NoError(t, err)

	const version = "1.19.5"

	output, err := golang.Base(c, golang.BaseImage(fmt.Sprintf("index.docker.io/library/golang:%s", version))).WithExec([]string{"go", "version"}).Stdout(ctx)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("go version go%s %s", version, platform), strings.TrimSpace(output))
}

func TestBase_Cgo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		opts     []golang.BaseOption
		expected string
	}{
		{
			name:     "default",
			expected: "",
		},
		{
			name:     "enabled",
			opts:     []golang.BaseOption{golang.EnableCgo()},
			expected: "1",
		},
		{
			name:     "disabled",
			opts:     []golang.BaseOption{golang.DisableCgo()},
			expected: "0",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx, c, buf := setupClient(t)
			t.Cleanup(func() { t.Log(buf.String()) })

			output, err := golang.Base(c, testCase.opts...).WithExec([]string{"sh", "-c", "echo $CGO_ENABLED"}).Stdout(ctx)
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, strings.TrimSpace(output))
		})
	}
}

func TestBase_ProjectRoot(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	output, err := golang.Base(c, golang.ProjectRoot("./testdata/empty")).WithExec([]string{"go", "list", "-m"}).Stdout(ctx)
	require.NoError(t, err)

	assert.Equal(t, "github.com/sagikazarmark/goci/lib/golang/testdata/empty", strings.TrimSpace(output))
}
