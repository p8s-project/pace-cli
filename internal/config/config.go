package config

// Config represents the structure of the ~/.pace/config.yaml file.
type Config struct {
	ActiveTap string `yaml:"active_tap"`
	Taps      []Tap  `yaml:"taps"`
}

// Tap represents a single tap configuration.
type Tap struct {
	Name    string  `yaml:"name"`
	URL     string  `yaml:"url"`
	Adapter Adapter `yaml:"adapter"`
}

// Adapter represents a platform adapter configuration.
type Adapter struct {
	RemoteBackend RemoteBackend `yaml:"remote_backend"`
}

// RemoteBackend represents a Terraform remote backend configuration.
type RemoteBackend struct {
	S3 S3Backend `yaml:"s3"`
}

// S3Backend represents an S3 remote backend configuration.
type S3Backend struct {
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}
