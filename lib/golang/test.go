package golang

import (
	"dagger.io/dagger"
	"github.com/sagikazarmark/go-option"
	"github.com/sagikazarmark/goci/lib"
)

// TestOption configures test parameters.
type TestOption interface {
	applyTest(o *testOptions)
}

type testOptionFunc func(o *testOptions)

func (f testOptionFunc) applyTest(o *testOptions) {
	f(o)
}

// Verbose enables verbose logging.
func Verbose(v bool) TestOption {
	return testOptionFunc(func(o *testOptions) {
		o.Verbose = v
	})
}

// EnableRaceDetector enables the race detector.
func EnableRaceDetector() TestOption {
	return testOptionFunc(func(o *testOptions) {
		o.RaceDetectorEnabled = true
	})
}

// CoverMode sets the coverage mode for Go test.
type CoverMode uint32

func (v CoverMode) applyTest(o *testOptions) {
	o.CoverMode = v
}

func (v CoverMode) Valid() bool {
	return CoverModeUndefined < v && v <= AtomicCoverMode
}

func (v CoverMode) String() string {
	switch v {
	case SetCoverMode:
		return "set"

	case CountCoverMode:
		return "count"

	case AtomicCoverMode:
		return "atomic"
	}

	return ""
}

const (
	CoverModeUndefined CoverMode = iota
	SetCoverMode
	CountCoverMode
	AtomicCoverMode
)

// CoverProfile specifies the output file for coverage information.
func CoverProfile(v string) TestOption {
	return testOptionFunc(func(o *testOptions) {
		o.CoverProfile = v
	})
}

type testOptions struct {
	baseOptions baseOptions

	Verbose             bool
	RaceDetectorEnabled bool
	CoverMode           CoverMode
	CoverProfile        string
}

func Test(client lib.Client, opts ...TestOption) *dagger.Container {
	var options testOptions

	for _, opt := range opts {
		opt.applyTest(&options)
	}

	args := []string{"go", "test"}

	if options.Verbose {
		args = append(args, "-v")
	}

	if options.RaceDetectorEnabled {
		args = append(args, "-race")

		options.baseOptions.CgoEnabled = option.Some(true)
	}

	if options.CoverMode.Valid() {
		args = append(args, "-covermode", options.CoverMode.String())
	}

	if options.CoverProfile != "" {
		args = append(args, "-coverprofile", options.CoverProfile)
	}

	args = append(args, "./...")

	return base(client, options.baseOptions).WithExec(args)
}
