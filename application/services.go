package application

import (
	"context"

	"github.com/EarvinKayonga/rider/domain"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/entropy"
	"github.com/EarvinKayonga/rider/messaging"
	"github.com/EarvinKayonga/rider/storage"
	"github.com/EarvinKayonga/rider/transport"
)

func trip(c *cli.Context, m Metadata) error {
	return makeCancellable(func(ctx context.Context) error {
		config, err := configuration.GetTripConfiguration(configFromContext(c),
			databaseFromContext(c),
			nsqLookupFromContext(c))
		if err != nil {
			return errors.Wrap(err,
				"an error occured while reading trip configuration")
		}

		logger, statsd, err := createInfraTools(ctx, config.Logging, config.Monitoring)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while creating tools for infra")
		}

		database, err := storage.NewPostgresDatabase(ctx, config.Database, *logger)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while contacting database")
		}

		ctx = storage.NewContext(ctx, database)
		service, err := transport.NewTripService(ctx, m.ToMap(), *config, *logger, statsd)
		if err != nil {
			return errors.Wrap(
				err, "an error occured while initialising trip service")
		}

		listener, err := domain.ListenerToTripEvent(ctx, *config, *logger, database)
		if err != nil {
			return errors.Wrap(
				err, "an error occured while initialising background listener service")
		}

		ctx = entropy.NewContext(ctx, entropy.NewIDGenerator())
		return runServiceWithListener(ctx, config.Server, service, *logger, listener, func(ctx context.Context) {
			database.Close(ctx)
		})
	})
}

func bike(c *cli.Context, m Metadata) error {
	return makeCancellable(func(ctx context.Context) error {
		config, err := configuration.GetBikeConfiguration(configFromContext(c),
			databaseFromContext(c),
			nsqLookupFromContext(c))
		if err != nil {
			return errors.Wrap(err,
				"an error occured while reading bike configuration")
		}

		logger, statsd, err := createInfraTools(ctx, config.Logging, config.Monitoring)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while creating tools for infra")
		}

		database, err := storage.NewPostgresDatabase(ctx, config.Database, *logger)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while contacting database")
		}

		err = domain.PopulateDatabase(ctx, database)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while populate database")
		}

		ctx = storage.NewContext(ctx, database)
		service, err := transport.NewBikeService(ctx, m.ToMap(), *config, *logger, statsd)
		if err != nil {
			return errors.Wrap(
				err, "an error occured while initialising bike service")
		}

		listener, err := domain.ListenerToBikeEvent(ctx, *config, *logger, database)
		if err != nil {
			return errors.Wrap(
				err, "an error occured while initialising background listener service")
		}

		ctx = entropy.NewContext(ctx, entropy.NewIDGenerator())
		return runServiceWithListener(ctx, config.Server, service, *logger, listener, func(ctx context.Context) {
			database.Close(ctx)
		})
	})
}

func gateway(c *cli.Context, m Metadata) error {
	return makeCancellable(func(ctx context.Context) error {
		config, err := configuration.GetGatewayConfiguration(
			configFromContext(c),
			c.GlobalString("bike"),
			c.GlobalString("trip"),
			c.GlobalString("queue"))
		if err != nil {
			return errors.Wrap(err,
				"an error occured while reading gateway configuration")
		}

		logger, statsd, err := createInfraTools(ctx, config.Logging, config.Monitoring)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while creating tools for infra")
		}

		messenger, err := messaging.NewEmitter(ctx, config.Messaging.Emission, *logger)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while creating messenger")
		}

		test := "test"
		err = messenger.Emit(ctx, &test)
		if err != nil {
			return errors.Wrap(err,
				"an error occured while publishing test message")
		}

		service, err := transport.NewGatewayService(ctx, m.ToMap(), *config, *logger, statsd, messenger)
		if err != nil {
			return errors.Wrap(
				err, "an error occured while initialising gateway service")
		}

		ctx = entropy.NewContext(ctx, entropy.NewIDGenerator())
		return runService(ctx, config.Server, service, *logger,
			func(ctx context.Context) {})
	})
}
