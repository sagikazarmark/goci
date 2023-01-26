package golang

import (
	"fmt"

	"dagger.io/dagger"
)

const (
	defaultGolangciLintImageRepository = "docker.io/golangci/golangci-lint"
	defaultGolangciLintImageTag        = "latest"
)

// LintOption configures lint parameters.
type LintOption interface {
	applyLint(o *lintOptions)
}

type lintOptionFunc func(o *lintOptions)

func (f lintOptionFunc) applyLint(o *lintOptions) {
	f(o)
}

// SourceImageRepository specifies which image repository to use as a source for GolangCI Lint.
// The value should follow the [OCI content addressable format] (without a tag or digest).
//
// SourceImageRepository is ignored when a full image reference is provided using [SourceImage].
//
// [OCI content addressable format]: https://github.com/opencontainers/.github/blob/master/docs/docs/introduction/digests.md#unique-resource-identifiers
func SourceImageRepository(v string) LintOption {
	return sourceImageRepository(v)
}

type sourceImageRepository string

func (v sourceImageRepository) applyLint(o *lintOptions) {
	o.SourceImageRepository = string(v)
}

// SourceImageTag specifies which tag from an image repository to use as a source for GolangCI Lint.
//
// SourceImageTag is ignored when a full image reference is provided using [SourceImage].
func SourceImageTag(v string) LintOption {
	return sourceImageTag(v)
}

type sourceImageTag string

func (v sourceImageTag) applyLint(o *lintOptions) {
	o.SourceImageTag = string(v)
}

// LinterVersion is an alias to [SourceImageTag] to provide a more user-friendly alternative.
func LinterVersion(v string) LintOption {
	return SourceImageTag(v)
}

// SourceImage specifies which image to use as a source for GolangCI Lint.
// The value should follow the [OCI content addressable format].
//
// [OCI content addressable format]: https://github.com/opencontainers/.github/blob/master/docs/docs/introduction/digests.md#unique-resource-identifiers
func SourceImage(v string) LintOption {
	return sourceImage(v)
}

type sourceImage string

func (v sourceImage) applyLint(o *lintOptions) {
	o.SourceImage = string(v)
}

type lintOptions struct {
	baseOptions baseOptions

	SourceImageRepository string
	SourceImageTag        string
	SourceImage           string
}

func Lint(client *dagger.Client, opts ...LintOption) *dagger.Container {
	var options lintOptions

	for _, opt := range opts {
		opt.applyLint(&options)
	}

	sourceImageRepository := defaultGolangciLintImageRepository
	if options.SourceImageRepository != "" {
		sourceImageRepository = options.SourceImageRepository
	}

	sourceImageTag := defaultGolangciLintImageTag
	if options.SourceImageTag != "" {
		sourceImageTag = options.SourceImageTag
	}

	sourceImage := fmt.Sprintf("%s:%s", sourceImageRepository, sourceImageTag)
	if options.SourceImage != "" {
		sourceImage = options.SourceImage
	}

	bin := client.Container().
		From(sourceImage).
		File("/usr/bin/golangci-lint")

	args := []string{"golangci-lint", "run", "--verbose"}

	return base(client, options.baseOptions).
		WithFile("/usr/local/bin/golangci-lint", bin).
		WithExec(args)
}
