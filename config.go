package jobs

import (
	poolImpl "github.com/roadrunner-server/sdk/v4/pool"
)

const (
	// name used to set pipeline name
	pipelineName string = "name"
	priorityKey  string = "priority"
)

// Config defines settings for job broker, workers and job-pipeline mapping.
type Config struct {
	// NumPollers configures number of priority queue pollers
	// Default - num logical cores
	NumPollers int `mapstructure:"num_pollers"`

	// PipelineSize is the limit of a main jobs queue which consume Items from the drivers pipeline
	// Driver pipeline might be much larger than a main jobs queue
	PipelineSize uint64 `mapstructure:"pipeline_size"`

	// Timeout in seconds is the per-push limit to put the job into queue
	Timeout int `mapstructure:"timeout"`

	// Pool configures roadrunner workers pool.
	Pool *poolImpl.Config `mapstructure:"pool"`

	// Pipelines defines mapping between PHP job pipeline and associated job broker.
	Pipelines map[string]Pipeline `mapstructure:"pipelines"`

	// Consuming specifies names of pipelines to be consumed on service start.
	Consume []string `mapstructure:"consume"`
}

func (c *Config) InitDefaults() {
	if c.Pool == nil {
		c.Pool = &poolImpl.Config{}
	}

	if c.PipelineSize == 0 {
		c.PipelineSize = 1_000_000
	}

	for k := range c.Pipelines {
		// set the pipeline name
		c.Pipelines[k].With(pipelineName, k)
		c.Pipelines[k].With(priorityKey, int64(c.Pipelines[k].Int(priorityKey, 10)))
	}

	if c.Timeout == 0 {
		c.Timeout = 60
	}

	c.Pool.InitDefaults()

	// NumPollers is hardcoded because it should be slightly more than the number of workers
	// to properly load all workers
	if c.NumPollers == 0 {
		c.NumPollers = int(c.Pool.NumWorkers) + 2
	}
}
