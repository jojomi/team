package job

import "time"

const (
	defaultWeight = 0
)

type Metadata struct {
	JobType     string `yaml:"type"`
	Filename    string
	Name        string   `yaml:"name,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Unattended  *bool    `yaml:"unattended,omitempty"`
	Enabled     *bool    `yaml:"enabled,omitempty"`
	Next        string   `yaml:"next,omitempty"`
	Possible    []string `yaml:"possible,omitempty"`

	Timeout time.Duration // TODO implement

	Weight *int `yaml:"weight,omitempty"`

	Output string `yaml:"output,omitempty"`
	Pre    struct {
		Command string `yaml:"command,omitempty"`
	} `yaml:"pre,omitempty"`
	Post struct {
		Command string `yaml:"command,omitempty"`
	} `json:"post,omitempty"`
}

func (x Metadata) GetWeight() int {
	if x.Weight == nil {
		return defaultWeight
	}
	return *x.Weight
}

func (x Metadata) IsEnabled() bool {
	return x.Enabled == nil || *x.Enabled
}

func (x Metadata) IsDisabled() bool {
	return !x.IsEnabled()
}

func (x Metadata) IsUnattended() bool {
	return x.Unattended != nil && *x.Unattended
}
