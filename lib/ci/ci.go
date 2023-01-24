package ci

import (
	"os"

	"dagger.io/dagger"
	"github.com/sagikazarmark/goci/lib"
)

type ciDetectorClient struct {
	lib.Client
}

// NewCIDetectorClient encapsulates a [lib.Client] and adds CI environment variables if they are detected.
func NewCIDetectorClient(client lib.Client) lib.Client {
	return ciDetectorClient{client}
}

func (c ciDetectorClient) Container(opts ...dagger.ContainerOpts) *dagger.Container {
	container := c.Client.Container(opts...)

	if os.Getenv("GITHUB_ACTIONS") != "" {
		container = container.
			WithEnvVariable("GITHUB_ACTIONS", os.Getenv("GITHUB_ACTIONS")).
			WithEnvVariable("GITHUB_HEAD_REF", os.Getenv("GITHUB_HEAD_REF")).
			WithEnvVariable("GITHUB_REF", os.Getenv("GITHUB_REF")).
			WithEnvVariable("GITHUB_REPOSITORY", os.Getenv("GITHUB_REPOSITORY")).
			WithEnvVariable("GITHUB_RUN_ID", os.Getenv("GITHUB_RUN_ID")).
			WithEnvVariable("GITHUB_SERVER_URL", os.Getenv("GITHUB_SERVER_URL")).
			WithEnvVariable("GITHUB_SHA", os.Getenv("GITHUB_SHA")).
			WithEnvVariable("GITHUB_WORKFLOW", os.Getenv("GITHUB_WORKFLOW"))
	}

	return container
}
