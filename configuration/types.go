package configuration

import (
	"fmt"
	"strconv"
)

var (
	// ConfigType holds the supported type.
	ConfigType = "yaml"
)

// TripConfiguration specifies general configurations
// for the trip service.
type TripConfiguration struct {
	Server     Server
	Logging    Logging
	Monitoring Monitoring
}

// BikeConfiguration specifies general configurations
// for the bike service.
type BikeConfiguration struct {
	Server     Server
	Logging    Logging
	Monitoring Monitoring
}

// GatewayConfiguration specifies general configurations
// for the gateway service.
type GatewayConfiguration struct {
	Server     Server
	Logging    Logging
	Monitoring Monitoring
	Limiter    Limiter

	// BikeURL is the base url to bike service.
	BikeURL string

	// TripURL is the base url to trip service.
	TripURL string
}

// Server specifies http based configuration for the underlying server.
type Server struct {
	Port        int
	Certificate string
	PrivateKey  string
	Host        string
}

// String for Stringer interface.
func (s Server) String() string {
	return fmt.Sprintf("%s:%s", s.Host, strconv.FormatInt(int64(s.Port), 10))
}

// Logging holds log configuration.
type Logging struct {
	// Standard log level
	// Example:
	// 			panic,
	//			fatal,
	//			error,
	//			warn, warning,
	//			info,
	//			debug
	Level string

	// logging format
	// only json or plain text are supported
	Format string
}

const (
	// JSONFormat is for the logging format.
	JSONFormat = "json"
)

// Monitoring holds monitoring configuration.
// (for things, like statsd or else).
type Monitoring struct {
	Addr   string
	Prefix string
}

// Limiter for rate limit features.
type Limiter struct {
	Limit float64
	Burst int
}
