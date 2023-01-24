package lib

import "dagger.io/dagger"

// Client encapsulates [dagger.Client] and allows injecting custom logic into client calls.
//
// Common use cases are injecting information into every container (eg. CI/CD env vars).
type Client interface {
	Container(opts ...dagger.ContainerOpts) *dagger.Container
	Host() *dagger.Host
	CacheVolume(key string) *dagger.CacheVolume
}

// BaseClient is a base implementation of the [Client] interface.
//
// It can be embedded into alternative client implementations to provide default implementations for supported calls.
type BaseClient struct {
	client *dagger.Client
}

// NewBaseClient returns a new [BaseClient].
func NewBaseClient(client *dagger.Client) BaseClient {
	return BaseClient{client}
}

func (c BaseClient) Container(opts ...dagger.ContainerOpts) *dagger.Container {
	return c.client.Container(opts...)
}

func (c BaseClient) Host() *dagger.Host {
	return c.client.Host()
}

func (c BaseClient) CacheVolume(key string) *dagger.CacheVolume {
	return c.client.CacheVolume(key)
}
