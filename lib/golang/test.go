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

type testOptions struct {
	baseOptions baseOptions

	Verbose             bool
	RaceDetectorEnabled bool
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

	args = append(args, "./...")

	return base(client, options.baseOptions).WithExec(args)
}
