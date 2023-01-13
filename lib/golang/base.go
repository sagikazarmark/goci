package golang

import (
	"crypto/sha256"
	"fmt"

	"dagger.io/dagger"
	"github.com/sagikazarmark/go-option"
)

const (
	defaultBaseImageRepository = "docker.io/library/golang"
	defaultBaseImageTag        = "latest"
)

type baseOptions struct {
	BaseImageRepository string
	BaseImageTag        string
	BaseImage           string

	ProjectRoot string

	CgoEnabled option.Option[bool]
}

func (o baseOptions) HasCgoValue() bool {
	return o.CgoEnabled != nil && option.IsSome(o.CgoEnabled)
}

func (o baseOptions) CgoValue() string {
	return option.MapOr(o.CgoEnabled, "", func(v bool) string {
		if v {
			return "1"
		}

		return "0"
	})
}

// Option configures common parameters.
type Option interface {
	applyBase(o *baseOptions)
}

// BaseImageRepository specifies which image repository to use as a base image for Go operations.
// The value should follow the [OCI content addressable format] (without a tag or digest).
//
// BaseImageRepository is ignored when a full image reference is provided using BaseImage.
//
// [OCI content addressable format]: https://github.com/opencontainers/.github/blob/master/docs/docs/introduction/digests.md#unique-resource-identifiers
func BaseImageRepository(v string) Option {
	return baseImageRepository(v)
}

type baseImageRepository string

func (v baseImageRepository) applyBase(o *baseOptions) {
	o.BaseImageRepository = string(v)
}

// BaseImageTag specifies which tag from an image repository to use as a base image for Go operations.
//
// BaseImageTag is ignored when a full image reference is provided using BaseImage.
func BaseImageTag(v string) Option {
	return baseImageTag(v)
}

type baseImageTag string

func (v baseImageTag) applyBase(o *baseOptions) {
	o.BaseImageTag = string(v)
}

// Version is an alias to BaseImageTag to provide a more user-friendly alternative.
func Version(v string) Option {
	return BaseImageTag(v)
}

// BaseImage specifies which image to use as a base image for Go operations.
// The value should follow the [OCI content addressable format].
func BaseImage(v string) Option {
	return baseImage(v)
}

type baseImage string

func (v baseImage) applyBase(o *baseOptions) {
	o.BaseImage = string(v)
}

// EnableCgo sets CGO_ENABLED to 1.
//
// If not set, Go will fall back to the default value.
func EnableCgo() Option {
	return cgo(true)
}

// DisableCgo sets CGO_ENABLED to 0.
//
// If not set, Go will fall back to the default value.
func DisableCgo() Option {
	return cgo(false)
}

type cgo bool

func (v cgo) applyBase(o *baseOptions) {
	o.CgoEnabled = option.Some(bool(v))
}

// ProjectRoot sets the project root to an alternate path (relative or absolute).
func ProjectRoot(v string) Option {
	return projectRoot(v)
}

type projectRoot string

func (v projectRoot) applyBase(o *baseOptions) {
	o.ProjectRoot = string(v)
}

func Base(client *dagger.Client, opts ...Option) *dagger.Container {
	var options baseOptions

	for _, opt := range opts {
		opt.applyBase(&options)
	}

	return base(client, options)
}

func base(client *dagger.Client, options baseOptions) *dagger.Container {
	projectRoot := "."
	if options.ProjectRoot != "" {
		projectRoot = options.ProjectRoot
	}

	baseImageRepository := defaultBaseImageRepository
	if options.BaseImageRepository != "" {
		baseImageRepository = options.BaseImageRepository
	}

	baseImageTag := defaultBaseImageTag
	if options.BaseImageTag != "" {
		baseImageTag = options.BaseImageTag
	}

	baseImage := fmt.Sprintf("%s:%s", baseImageRepository, baseImageTag)
	if options.BaseImage != "" {
		baseImage = options.BaseImage
	}

	// Calculate a hash for the base image to use in cache volume names.
	h := sha256.New()
	h.Write([]byte(baseImage))

	baseImageHash := h.Sum(nil)

	container := client.Container().
		From(baseImage).
		WithMountedCache("/root/.cache/go-build", client.CacheVolume(fmt.Sprintf("go-build-%x", baseImageHash))).
		WithMountedCache("/go/pkg/mod", client.CacheVolume(fmt.Sprintf("go-mod-%x", baseImageHash))).
		WithMountedDirectory("/src", client.Host().Directory(projectRoot)).
		WithWorkdir("/src")

	if options.HasCgoValue() {
		container = container.WithEnvVariable("CGO_ENABLED", options.CgoValue())
	}

	return container
}
