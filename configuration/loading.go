package configuration

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

func GetGatewayConfiguration(path, bikeURL, tripURL string) (*GatewayConfiguration, error) {
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

	return config, nil
}

func GetBikeConfiguration(path string) (*BikeConfiguration, error) {
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

	return config, nil
}

func GetTripConfiguration(path string) (*TripConfiguration, error) {
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

	return config, nil
}
