package application

import (
	"fmt"

	"github.com/urfave/cli"
)

// RunBike is a wrapper in order to keep the the main function tidy.
func RunBike(args []string, m Metadata) error {
	application := &Application{
		&cli.App{
			Name:      AppNameSpace + Bike,
			Author:    AppAuthor,
			Copyright: AppCopyright,

			EnableBashCompletion: true,

			UsageText: "the API for Bikes",
			Usage:     "a Service for Handling The Fleet of Bikes",

			Version: fmt.Sprintf(
				"Branch: %s, Compiler: %s, CompiledAt: %s, Commit: %s",
				m.Branch, m.Compiler, m.CompiledAt, m.Sha),

			Metadata: m.ToMap(),

			Flags: []cli.Flag{
				configurationFlag(),
				databaseFlag(),
				nsqLookUpFlag(),
			},

			Action: func(c *cli.Context) error {
				return bike(c, m)
			},
		},
	}

	return application.Run(args)
}

// RunGateway is a wrapper in order to keep the the main function tidy.
func RunGateway(args []string, m Metadata) error {
	application := &Application{
		&cli.App{
			Name:      AppNameSpace + Gateway,
			Author:    AppAuthor,
			Copyright: AppCopyright,

			EnableBashCompletion: true,

			UsageText: "Rider's API Gateway",
			Usage:     "Rider's API Gateway",

			Version: fmt.Sprintf(
				"Branch: %s, Compiler: %s, CompiledAt: %s, Commit: %s",
				m.Branch, m.Compiler, m.CompiledAt, m.Sha),

			Metadata: m.ToMap(),

			Flags: []cli.Flag{
				configurationFlag(),

				cli.StringFlag{
					Name:   "bike",
					Usage:  "the base url to bike service",
					EnvVar: "BIKE_URL",
					Value:  "",
				},

				cli.StringFlag{
					Name:   "trip",
					Usage:  "the base url to trip service",
					EnvVar: "TRIP_URL",
					Value:  "",
				},

				cli.StringFlag{
					Name:   "queue",
					Usage:  "the socket to nsq service",
					EnvVar: "NSQ_SOCKET",
					Value:  "0.0.0.0:4150",
				},
			},

			Action: func(c *cli.Context) error {
				return gateway(c, m)
			},
		},
	}

	return application.Run(args)
}

// RunTrip is a wrapper in order to keep the the main function tidy.
func RunTrip(args []string, m Metadata) error {
	application := &Application{
		&cli.App{
			Name:      AppNameSpace + Trip,
			Author:    AppAuthor,
			Copyright: AppCopyright,

			EnableBashCompletion: true,

			UsageText: "the API for Trip",
			Usage:     "a Service for Trip Handling",

			Version: fmt.Sprintf(
				"Branch: %s, Compiler: %s, CompiledAt: %s, Commit: %s",
				m.Branch, m.Compiler, m.CompiledAt, m.Sha),

			Metadata: m.ToMap(),

			Flags: []cli.Flag{
				configurationFlag(),
				databaseFlag(),
				nsqLookUpFlag(),
			},

			Action: func(c *cli.Context) error {
				return trip(c, m)
			},
		},
	}

	return application.Run(args)
}

// Application represents the cli application which will run when the binary is run.
type Application struct {
	*cli.App
}
