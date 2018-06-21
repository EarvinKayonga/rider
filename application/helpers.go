package application

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/logging"
	"github.com/EarvinKayonga/rider/stats"
)

const (
	// AppAuthor is the app's author.
	AppAuthor = "Gilbert"

	// AppNameSpace is a prefix for the app's name.
	AppNameSpace = "rider/"

	// AppCopyright specifies the app's copyright.
	AppCopyright = "Gilbert and Co"

	// DefaultConfigurationPath specifies the default path of the configuration file.
	DefaultConfigurationPath = "configuration.yml"

	// ConfigurationFlag specifies the path of the configuration file.
	ConfigurationFlag = "RIDER_CONF"

	// DatabaseFlag specifies the postgres URL to an instance.
	DatabaseFlag = "DATABASE_URL"
)

// AppNames for binaries.
const (
	Bike    = "bike"
	Trip    = "trip"
	Gateway = "gateway"
)

// configurationFlag creates a Flag for the cli app.
func configurationFlag() cli.StringFlag {
	return cli.StringFlag{
		Name:   "configuration",
		Usage:  "path to configuration file (yml)",
		EnvVar: ConfigurationFlag,
		Value:  DefaultConfigurationPath,
	}
}

// databaseFlag  creates a Flag for the cli app.
func databaseFlag() cli.StringFlag {
	return cli.StringFlag{
		Name:   "database",
		Usage:  "postgres URL, example: `postgresql://user:password@host:port/name`",
		EnvVar: DatabaseFlag,
	}
}

// nsqLookUpFlag  creates a Flag for the cli app.
func nsqLookUpFlag() cli.StringFlag {
	return cli.StringFlag{
		Name:   "lookup",
		Usage:  "socket for nsq lookup example: `0.0.0.0:4161`",
		EnvVar: "LOOKUP",
	}
}

// configFromContext reads the configuration file from context.
func configFromContext(ctx *cli.Context) string {
	path := ctx.GlobalString("configuration")
	if path == "" {
		return DefaultConfigurationPath
	}

	return path
}

func databaseFromContext(ctx *cli.Context) string {
	return ctx.GlobalString("database")
}

func nsqLookupFromContext(ctx *cli.Context) string {
	return ctx.GlobalString("lookup")
}

// createInfraTools returns a logger and statter.
func createInfraTools(ctx context.Context, logInfo configuration.Logging,
	monitInfo configuration.Monitoring) (*logging.Logger, stats.Statter, error) {

	logger := logging.NewLogger(logInfo)
	statsd, err := stats.NewStatsdClient(monitInfo)
	if err != nil {
		return nil, nil, errors.Wrapf(err,
			"an error occured while creating statter with conf %v", monitInfo)
	}

	return logger, statsd, nil
}
