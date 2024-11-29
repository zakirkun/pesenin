package bootstrap

import (
	"context"
	_log "log"
	"os"

	"github.com/zakirkun/pesnin/infrastructure"
	"github.com/zakirkun/pesnin/pkg/cache"
	"github.com/zakirkun/pesnin/pkg/database"
	"github.com/zakirkun/pesnin/pkg/logstash"
	"github.com/zakirkun/pesnin/pkg/minio"
	"github.com/zakirkun/pesnin/pkg/queue"
	"github.com/zakirkun/pesnin/pkg/tracer"
)

type BootstrapContext struct {
	builder *infrastructure.InfrastructureContext
}

func New(builder *infrastructure.InfrastructureContext) *BootstrapContext {
	return &BootstrapContext{
		builder: builder,
	}
}

func (b *BootstrapContext) InitApp() {

	// init logger
	b.logger()

	if b.builder.Logstash.Addr != "" {
		b.logstash()
	}

	if b.builder.Cache.Addr != "" {
		b.cache()
	}

	if b.builder.Queue.Address != "" {
		b.queue()
	}

	if b.builder.Minio.Endpoint != "" {
		b.minio()
	}

	if b.builder.Database.Host != "" {
		b.database()
	}

	if b.builder.Otel.Endpoint != "" {
		b.otel()
	}

	if b.builder.Server.Host != "" {
		b.serve()

	}

}

func (b *BootstrapContext) logger() {
	// init logger
	b.builder.Logger.
		OpenFileLogs(b.builder.Logger.FileName).
		Open()
}

func (b *BootstrapContext) database() {
	client, err := b.builder.Database.OpenDB()
	if err != nil {
		_log.Fatalf("Failed to open database connection: %v", err)
		os.Exit(1)
	}

	database.DB = client
	_log.Println("Database initialized successfully")
}

func (b *BootstrapContext) cache() {
	client := b.builder.Cache.Open()
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		_log.Fatalf("Failed to open database cache: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	cache.CACHE = client
}

func (b *BootstrapContext) queue() {
	_, err := b.builder.Queue.Open()
	if err != nil {
		_log.Fatalf("Failed to open rabbitmq: %v", err)
		os.Exit(1)
	}

	queue.RMQ = &b.builder.Queue
}

func (b *BootstrapContext) minio() {
	_, err := b.builder.Minio.Open()
	if err != nil {
		_log.Fatalf("Failed to open minio: %v", err)
		os.Exit(1)
	}

	minio.Object = &b.builder.Minio

}

func (b *BootstrapContext) logstash() {
	_, err := b.builder.Logstash.Open()
	if err != nil {
		_log.Fatalf("Failed to open minio: %v", err)
		os.Exit(1)
	}

	logstash.LOGSTASH = &b.builder.Logstash
}

func (b *BootstrapContext) serve() {
	b.builder.Server.Run()
}

func (b *BootstrapContext) otel() {
	tp, err := b.builder.Otel.Open()
	if err != nil {
		_log.Fatalf("Failed to open minio: %v", err)
		os.Exit(1)
	}

	tracer.Tracer = tp.Tracer(b.builder.ServicesName)

	defer func() {
		_ = tp.Shutdown(context.Background())
	}()
}
