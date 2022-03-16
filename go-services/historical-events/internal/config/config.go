package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

// Config Can be used as dependency in other packages
type Config struct {
	App app
	DB  Db
}

var (
	c *Config
)

// App config
type app struct {
	AppEnv   string
	GRPCPort string
	HTTPort  string
}

// DB config
type Db struct {
	Uri          string
	ReadTimeOut  uint8
	WriteTimeOut uint8
}

// LoadConfig loads the config ensuring all the required variables are there
func LoadConfig() (*Config, error) {

	if c != nil {
		return c, nil
	}

	appInstance := app{
		AppEnv:   os.Getenv("APP_ENV"),
		GRPCPort: os.Getenv("GRPC_PORT"),
		HTTPort:  os.Getenv("HTTP_PORT"),
	}

	mongoReadTimeOut, err := strconv.Atoi(os.Getenv("MONGODB_READ_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	mongoWriteTimeOut, err := strconv.Atoi(os.Getenv("MONGODB_WRITE_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	db := Db{
		Uri:          os.Getenv("MONGODB_URI"),
		ReadTimeOut:  uint8(mongoReadTimeOut),
		WriteTimeOut: uint8(mongoWriteTimeOut),
	}

	// check all values for empty strings / return appropriate error
	err = checkMissingValues(appInstance, db)

	if err != nil {
		log.Printf("err in config %v", err)

		return nil, err
	}

	c = &Config{
		App: appInstance,
		DB:  db,
	}

	return c, nil
}

func checkMissingValues(configs ...interface{}) error {
	for _, config := range configs {
		v := reflect.ValueOf(config)
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).String() == "" {
				return fmt.Errorf("missing required env variable :: %s in configuration %s", v.Type().Field(i).Name, reflect.TypeOf(config).String())
			}
		}
	}
	return nil
}
