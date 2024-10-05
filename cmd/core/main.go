package main

import (
	"flag"

	"github.com/zakirkun/pesnin/bootstrap"
	"github.com/zakirkun/pesnin/bundle"
	"github.com/zakirkun/pesnin/infrastructure"
)

var configFile *string

func init() {
	// parse config
	configFile = flag.String("c", "config.toml", "configuration file")
	flag.Parse()

	// load config
	bundle.SetConfig(configFile)
}

func main() {

	builder := infrastructure.NewInfrastructureBuilder().
		WithServiceName("services.core.app").
		InitLogger(bundle.SetLogger()).
		WithQueue(bundle.SetQueue()).
		WithLogstash(bundle.SetLogstash()).
		WithDatabase(bundle.SetDatabase()).
		WithMinio(bundle.SetMinio()).
		WithServer(bundle.SetWebServer())

	runApp := bootstrap.New(builder.Build())
	runApp.InitApp()
}
