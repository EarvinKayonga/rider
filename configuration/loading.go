package configuration

import (
	"net"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// ConfigType holds the supported type.
	ConfigType = "yaml"
)

func loadConfiguration(path string) error {
	viper.Set("Verbose", false)
	viper.SetConfigType(ConfigType)
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "an error occured while reading %s", path)
	}

	return nil
}

// GetGatewayConfiguration returns valid configuration to run the API gateway.
func GetGatewayConfiguration(path, bikeURL, tripURL string, nsq string) (*GatewayConfiguration, error) {
	err := loadConfiguration(path)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while loading configuration")
	}

	config := &GatewayConfiguration{
		Server: Server{
			Port: 8080,
		},

		Logging: Logging{
			Level: "debug",
		},
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrapf(err,
			"an error occured while unmarshalling file: %s", path)
	}

	// overwriting values in configuration file
	// with ENV values.
	if bikeURL != "" {
		config.BikeURL = bikeURL
	}

	if tripURL != "" {
		config.TripURL = tripURL
	}

	if nsq != "" {
		config.Messaging.Emission.Address = nsq
	}

	return config, nil
}

// GetBikeConfiguration returns valid configuration to run a bike service.
func GetBikeConfiguration(path, databaseURL, consumerSOCKET string) (*BikeConfiguration, error) {
	err := loadConfiguration(path)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while loading configuration")
	}

	config := &BikeConfiguration{
		Server: Server{
			Port: 8081,
		},

		Logging: Logging{
			Level: "debug",
		},
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrapf(err,
			"an error occured while unmarshalling file: %s", path)
	}

	if databaseURL != "" {
		database, err := parseDatabaseURL(databaseURL)
		if err != nil {
			return nil, errors.Wrapf(err, "an error occured while parsing database url: %s", databaseURL)
		}

		config.Database = *database
	}

	if consumerSOCKET != "" {
		config.Messaging.Consumption.Address = consumerSOCKET
	}

	return config, nil
}

// GetTripConfiguration returns valid configuration to run a trip service.
func GetTripConfiguration(path, databaseURL, consumerSOCKET string) (*TripConfiguration, error) {
	err := loadConfiguration(path)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while loading configuration")
	}

	config := &TripConfiguration{
		Server: Server{
			Port: 8082,
		},

		Logging: Logging{
			Level: "debug",
		},
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrapf(err,
			"an error occured while unmarshalling file: %s", path)
	}

	if databaseURL != "" {
		database, err := parseDatabaseURL(databaseURL)
		if err != nil {
			return nil, errors.Wrapf(err, "an error occured while parsing database url: %s", databaseURL)
		}

		config.Database = *database
	}

	if consumerSOCKET != "" {
		config.Messaging.Consumption.Address = consumerSOCKET
	}

	return config, nil
}

func parseDatabaseURL(databaseURL string) (*Database, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return nil, errors.Wrapf(err, "an error occured while parsing database url: %s", databaseURL)
	}

	password, ok := u.User.Password()
	if !ok {
		password = ""
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil,
			errors.Wrapf(err, "an error occured while splitting host and port in database url: %s",
				u.Host)
	}

	if port == "" {
		port = "5432"
	}

	return &Database{
		Host:     host,
		Port:     port,
		User:     u.User.Username(),
		Password: password,
		Name:     strings.TrimLeft(u.Path, "/"),
	}, nil
}
