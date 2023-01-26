package golang_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/goci/lib/golang"
)

func TestTest(t *testing.T) {
	t.Parallel()

	ctx, c, buf := setupClient(t)
	t.Cleanup(func() { t.Log(buf.String()) })

	_, err := golang.Test(c, golang.ProjectRoot("./testdata/test")).ExitCode(ctx)
	require.NoError(t, err)
}

func TestCoverMode(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			value    golang.CoverMode
			expected bool
		}{
			{
				value:    golang.SetCoverMode,
				expected: true,
			},
			{
				value:    golang.CountCoverMode,
				expected: true,
			},
			{
				value:    golang.AtomicCoverMode,
				expected: true,
			},

			{
				value:    golang.CoverModeUndefined,
				expected: false,
			},
			{
				value:    0,
				expected: false,
			},
			{
				value:    99,
				expected: false,
			},
		}

		for _, testCase := range testCases {
			testCase := testCase

			t.Run("", func(t *testing.T) {
				t.Parallel()

				if got, want := testCase.value.Valid(), testCase.expected; got != want {
					t.Errorf("validation does not match the expected result\nactual:   %t\nexpected: %t", got, want)
				}
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			value    golang.CoverMode
			expected string
		}{
			{
				value:    golang.SetCoverMode,
				expected: "set",
			},
			{
				value:    golang.CountCoverMode,
				expected: "count",
			},
			{
				value:    golang.AtomicCoverMode,
				expected: "atomic",
			},

			{
				value:    golang.CoverModeUndefined,
				expected: "",
			},
			{
				value:    0,
				expected: "",
			},
			{
				value:    99,
				expected: "",
			},
		}

		for _, testCase := range testCases {
			testCase := testCase

			t.Run("", func(t *testing.T) {
				t.Parallel()

				if got, want := testCase.value.String(), testCase.expected; got != want {
					t.Errorf("validation does not match the expected result\nactual:   %s\nexpected: %s", got, want)
				}
			})
		}
	})
}
