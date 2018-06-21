package application

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/logging"
	"github.com/EarvinKayonga/rider/messaging"
	"github.com/EarvinKayonga/rider/transport"
)

// makeCancellable attachs context.Done to os signals.
// Make Ctrl+C cancel the context.Context.
func makeCancellable(callback func(context.Context) error) error {
	osSignals := make(chan os.Signal, 1)

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	signal.Notify(osSignals, os.Interrupt)
	defer func() {
		signal.Stop(osSignals)
		cancel()
	}()

	go func() {
		select {
		case <-osSignals:
			cancel()
		case <-ctx.Done():
		}
	}()

	return callback(ctx)
}

// runService handles graceful shutdown of the http server.
func runService(ctx context.Context, conf configuration.Server, server *http.Server, logger logging.Logger,
	onClose func(ctx context.Context)) error {
	socket, err := net.Listen("tcp", conf.String())
	if err != nil {
		return errors.Wrapf(err, "cannot listen port: %d", conf.Port)
	}

	errChan := make(chan error, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		logger.Infof("running service on %s", conf)

		if conf.Certificate != "" && conf.PrivateKey != "" {
			server.TLSConfig = transport.TLSConfiguration
			server.TLSNextProto = transport.TLSNextProto

			err = server.ServeTLS(socket, conf.Certificate, conf.PrivateKey)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot start https server")
			}
		} else {
			err = server.Serve(socket)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot start http server")
			}
		}
	}()

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = server.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Info("an error occured while shutting down server")
		}
		logger.Info("server gracefully stopping")

	case err := <-errChan:
		return err
	}

	logger.Warning("the server is shutting everything down")
	onClose(ctx)

	logger.Warning("complete shutdown")
	return nil
}

// runService handles graceful shutdown of the http server with an additionnal background listener.
func runServiceWithListener(ctx context.Context, conf configuration.Server, server *http.Server, logger logging.Logger,
	listener messaging.Consumer,
	onClose func(ctx context.Context)) error {
	socket, err := net.Listen("tcp", conf.String())
	if err != nil {
		return errors.Wrapf(err, "cannot listen port: %d", conf.Port)
	}

	errChan := make(chan error, 2)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	logger.Infof("running service on %s", conf)

	go func() {
		if conf.Certificate != "" && conf.PrivateKey != "" {
			server.TLSConfig = transport.TLSConfiguration
			server.TLSNextProto = transport.TLSNextProto

			err = server.ServeTLS(socket, conf.Certificate, conf.PrivateKey)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot start https server")
			}
		} else {
			err = server.Serve(socket)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot start http server")
			}
		}
	}()

	logger.Infof("launching background listener")
	go func() {
		err := listener.Run(ctx)
		if err != nil {
			errChan <- errors.Wrap(err, "failed launching background listener")
		}
	}()

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = server.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Info("an error occured while shutting down server")
		}
		logger.Info("server gracefully stopping")

	case err := <-errChan:
		return err
	}

	logger.Warning("the server is shutting everything down")
	onClose(ctx)

	logger.Warning("complete shutdown")
	return nil
}
