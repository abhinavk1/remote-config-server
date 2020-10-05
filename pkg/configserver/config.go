package configserver

type ServerConfig struct {
	PrivateKeyPath   string
	RepoUrl          string
	WorkingDirectory string
	Port             int
}
