package config

import (
	"os"
	"strconv"
)

// Config holds the main application configuration settings.
type Config struct {
	Proxy         ProxyConfig         // Proxy configuration.
	Mongo         MongoDBConfig       // MongoDB configuration.
	SourceHandler SourceHandlerConfig // Source handler configuration.
	AuthServer    AuthServerConfig    // AuthServer holds configuration details for the Auth service.
	VacancyServer VacancyServerConfig // VacancyServer holds configuration details for the Vacancy service.
	Env           string              // Environment type (e.g., dev, prod).
}

// VacancyServerConfig contains connection details for the Vacancy service.
type VacancyServerConfig struct {
	Address string // Address is the address the Vacancy service listens on.
}

// AuthServerConfig contains connection details for the Auth service.
type AuthServerConfig struct {
	Address string // Address is the address the Auth service listens on.
	Issuer  string // Issuer is the identifier of the entity issuing tokens (e.g., "grpc.pulse-finder.bot").
}

// ProxyConfig holds configuration settings for Proxy.
type ProxyConfig struct {
	Host            string // Host is the hostname or IP address of the Proxy server.
	Port            string // Port is the port number of the Proxy server.
	ControlPassword string // ControlPassword is the password for the Proxy control port.
	ControlPort     string // ControlPort is the port number of the Proxy control port.
	PingUrl         string // PingUrl is the URL used to check the proxy's status or connectivity.
}

// MongoDBConfig holds configuration settings for MongoDB.
type MongoDBConfig struct {
	Host              string // Host is the hostname or IP address of the MongoDB server.
	Port              string // Port is the port number of the MongoDB server.
	User              string // User is the username for connecting to the MongoDB server.
	Pass              string // Pass is the password for connecting to the MongoDB server.
	DB                string // DB is the name of the MongoDB database.
	UrlsCollection    string // UrlsCollection is the name of the MongoDB collection.
	VacancyCollection string // VacancyCollection is the name of the MongoDB collection.
}

// SourceHandlerConfig holds configuration settings for Source Handlers.
type SourceHandlerConfig struct {
	Alfa      SourceConfig // Alfa source handler configuration.
	Beta      SourceConfig // Beta source handler configuration.
	Gamma     SourceConfig // Gamma source handler configuration.
	BatchSize int          // BatchSize is the batch size for processing.
}

// SourceConfig represents configuration for a single source.
type SourceConfig struct {
	SitemapURL string // URL of the sitemap or RSS feed.
}

// LoadConfig loads the configuration settings from environment variables, falling back to default values.
func LoadConfig() *Config {
	config := &Config{
		Proxy: ProxyConfig{
			Host:            getEnv("PROXY_HOST", ""),
			Port:            getEnv("PROXY_PORT", ""),
			ControlPassword: getEnv("PROXY_CONTROL_PASSWORD", ""),
			ControlPort:     getEnv("PROXY_CONTROL_PORT", ""),
			PingUrl:         getEnv("PROXY_PING_URL", ""),
		},
		Mongo: MongoDBConfig{
			Host:              getEnv("MONGO_HOST", ""),
			Port:              getEnv("MONGO_PORT", ""),
			User:              getEnv("MONGO_USER", ""),
			Pass:              getEnv("MONGO_PASS", ""),
			DB:                getEnv("MONGO_DB", ""),
			UrlsCollection:    getEnv("MONGO_URLS_COLLECTION", ""),
			VacancyCollection: getEnv("MONGO_VACANCY_COLLECTION", ""),
		},
		SourceHandler: SourceHandlerConfig{
			Alfa: SourceConfig{
				SitemapURL: getEnv("SOURCE_ALFA_SITEMAP_URL", "example.com"),
			},
			Beta: SourceConfig{
				SitemapURL: getEnv("SOURCE_BETA_SITEMAP_URL", ""),
			},
			Gamma: SourceConfig{
				SitemapURL: getEnv("SOURCE_GAMMA_SITEMAP_URL", ""),
			},
			BatchSize: getEnvAsInt("SOURCE_BATCH_SIZE", 1),
		},
		AuthServer: AuthServerConfig{
			Address: getEnv("AUTH_SERVER_ADDRESS", ""),
			Issuer:  getEnv("AUTH_ISSUER", ""),
		},
		VacancyServer: VacancyServerConfig{
			Address: getEnv("VACANCY_SERVER_ADDRESS", ""),
		},
		Env: getEnv("ENV", "dev"),
	}

	return config
}

// getEnv fetches the value of an environment variable or returns a fallback.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

// getEnvAsInt fetches the value of an environment variable as an integer or returns a fallback.
func getEnvAsInt(key string, fallback int) int {
	v := getEnv(key, "")
	if value, err := strconv.Atoi(v); err == nil {
		return value
	}
	return fallback
}
