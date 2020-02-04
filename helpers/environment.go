package helpers

import (
	"fmt"
	"os"
)

// ServerConfiguration contains server config properties
type ServerConfiguration struct {
	Port         string
	ClientID     string
	ClientSecret string
}

// GetConfiguration reads environment variables to get server config properties
func GetConfiguration() (ServerConfiguration, error) {
	port, err := getEnv("PORT")
	if err != nil {
		return ServerConfiguration{}, err
	}
	port = getPort(port)

	clientID, err := getEnv("CLIENT_ID")
	if err != nil {
		return ServerConfiguration{}, err
	}

	clientSecret, err := getEnv("CLIENT_SECRET")
	if err != nil {
		return ServerConfiguration{}, err
	}

	return ServerConfiguration{
		Port:         port,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

func getEnv(name string) (string, error) {
	envValue := os.Getenv(name)
	if envValue == "" {
		return "", fmt.Errorf("Env '%v' not set", name)
	}

	return envValue, nil
}

func getPort(portNumber string) string {
	return fmt.Sprintf(":%v", portNumber)
}
