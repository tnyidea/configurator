# configurator

A simple Go configurator library

## Example
Create your own struct and optionally tag the fields with default values, environment variables and
if the value is required
```go
package myserver

import "github.com/tnyidea/configurator"

type Server struct {
	Port    string `default:"80" env:"PORT" config:"required"`
	AppPath string `default:"ui/app/build" env:"APP_PATH" config:"required"`
	DbUrl   string `env:"DB_URL" config:"required"`
}

func NewServer(configFile string) (Server, error) {
	var server Server
	err := configurator.ParseEnvConfig(&server, configFile)
	if err != nil {
		return Server{}, err
	}
	return server, nil
}
```
