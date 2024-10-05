package bundle

import (
	logs "log"
	"os"
	"time"

	"github.com/zakirkun/pesnin/pkg/config"
	"github.com/zakirkun/pesnin/pkg/database"
	"github.com/zakirkun/pesnin/pkg/log"
	"github.com/zakirkun/pesnin/pkg/logstash"
	"github.com/zakirkun/pesnin/pkg/minio"
	"github.com/zakirkun/pesnin/pkg/queue"
	"github.com/zakirkun/pesnin/pkg/server"
)

func SetConfig(configFile *string) {
	cfg := config.NewConfig(*configFile)
	if err := cfg.Initialize(); err != nil {
		logs.Fatalf("Error reading config : %v", err)
		os.Exit(1)
	}
}

func SetLogger() log.LoggerContext {
	return log.LoggerContext{
		Debug:    config.GetString("server.mode"),
		FileName: config.GetString("server.filename"),
	}
}

func SetDatabase() database.DBModel {
	return database.DBModel{
		ServerMode:   config.GetString("server.mode"),
		Driver:       config.GetString("database.db_driver"),
		Host:         config.GetString("database.db_host"),
		Port:         config.GetString("database.db_port"),
		Name:         config.GetString("database.db_name"),
		Username:     config.GetString("database.db_username"),
		Password:     config.GetString("database.db_password"),
		MaxIdleConn:  config.GetInt("pool.conn_idle"),
		MaxOpenConn:  config.GetInt("pool.conn_max"),
		ConnLifeTime: config.GetInt("pool.conn_lifetime"),
	}
}

func SetWebServer() server.ServerContext {
	return server.ServerContext{
		Host:         ":" + config.GetString("server.port"),
		Handler:      nil,
		ReadTimeout:  time.Duration(config.GetInt("server.http_timeout")),
		WriteTimeout: time.Duration(config.GetInt("server.http_timeout")),
	}
}

func SetQueue() queue.RabbitMQ {
	return queue.RabbitMQ{
		Address: config.GetString("queue.rabbimq_url"),
	}
}

func SetLogstash() logstash.LogstashModel {
	return logstash.LogstashModel{
		Network: config.GetString("logstash.network"),
		Addr:    config.GetString("logstash.addr"),
	}
}

func SetMinio() minio.MinioContext {
	return minio.MinioContext{
		Endpoint:        config.GetString("minio.endpoint"),
		AccessKeyID:     config.GetString("minio.access_key"),
		SecretAccessKey: config.GetString("minio.access_secret"),
		UseSSL:          config.GetBool("minio.use_ssl"),
	}
}
