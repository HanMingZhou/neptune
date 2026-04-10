package config

type Notebook struct {
	EphemeralStorageSize string `mapstructure:"ephemeral-storage-size" json:"ephemeral-storage-size" yaml:"ephemeral-storage-size"`
}
