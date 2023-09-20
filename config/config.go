package config

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"
)

func MustGetEnv(key string) string {
	value, success := os.LookupEnv(key)
	if !success {
		log.Fatalln("failed to get: ", key)
	}
	return value
}

func GetEnvOrDefault(key, d string) string {
	value, success := os.LookupEnv(key)
	if !success {
		return d
	}
	return value
}

func Connect() *nats.Conn {
	url := GetEnvOrDefault("NATS_URL", nats.DefaultURL)
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalln(err)
	}
	return nc
}

func JetStream() nats.JetStreamContext {
	nc := Connect()
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalln(err)
	}
	return js
}

func ReadConfig(filename string, conf interface{}) {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("err reading yaml: ", err)
	}
	err = yaml.Unmarshal(f, conf)
	if err != nil {
		log.Fatalln("err unmarshal: ", err)
	}
}
