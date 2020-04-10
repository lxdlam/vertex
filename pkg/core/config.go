package core

// Config represents the core config of vertex
// It can be parsed from a file or passed from cli options
type Config struct {
	DbPath  string
	LogPath string
}

func NewConfig() *Config {

}
