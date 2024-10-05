package infrastructure

import (
	"github.com/zakirkun/pesnin/pkg/cache"
	"github.com/zakirkun/pesnin/pkg/database"
	"github.com/zakirkun/pesnin/pkg/log"
	"github.com/zakirkun/pesnin/pkg/logstash"
	"github.com/zakirkun/pesnin/pkg/minio"
	"github.com/zakirkun/pesnin/pkg/queue"
	"github.com/zakirkun/pesnin/pkg/server"
)

type InfrastructureContext struct {
	Database database.DBModel
	Cache    cache.Cache
	Queue    queue.RabbitMQ
	Minio    minio.MinioContext
	Logstash logstash.LogstashModel
	Server   server.ServerContext
	Logger   log.LoggerContext

	ServicesName string
}

// InfrastructureBuilder struct for creating Infrastructure with selected dependencies
type InfrastructureBuilder struct {
	Infra *InfrastructureContext
}

// NewInfrastructureBuilder initializes the builder
func NewInfrastructureBuilder() *InfrastructureBuilder {
	return &InfrastructureBuilder{
		Infra: &InfrastructureContext{},
	}
}

// Init Logger adds the logger component to the infrastructure
func (b *InfrastructureBuilder) InitLogger(logger log.LoggerContext) *InfrastructureBuilder {
	b.Infra.Logger = logger
	return b
}

// WithDatabase adds the database component to the infrastructure
func (b *InfrastructureBuilder) WithDatabase(db database.DBModel) *InfrastructureBuilder {
	b.Infra.Database = db
	return b
}

// WithCache adds the cache component to the infrastructure
func (b *InfrastructureBuilder) WithCache(c cache.Cache) *InfrastructureBuilder {
	b.Infra.Cache = c
	return b
}

// WithQueue adds the queue component to the infrastructure
func (b *InfrastructureBuilder) WithQueue(q queue.RabbitMQ) *InfrastructureBuilder {
	b.Infra.Queue = q
	return b
}

// WithMinio adds the MinIO component to the infrastructure
func (b *InfrastructureBuilder) WithMinio(m minio.MinioContext) *InfrastructureBuilder {
	b.Infra.Minio = m
	return b
}

// WithLogstash adds the Logstash component to the infrastructure
func (b *InfrastructureBuilder) WithLogstash(l logstash.LogstashModel) *InfrastructureBuilder {
	b.Infra.Logstash = l
	return b
}

// WithServer adds the server component to the infrastructure
func (b *InfrastructureBuilder) WithServer(s server.ServerContext) *InfrastructureBuilder {
	b.Infra.Server = s
	return b
}

// WithServiceName sets the service name in the infrastructure
func (b *InfrastructureBuilder) WithServiceName(name string) *InfrastructureBuilder {
	b.Infra.Server.AppName = name
	b.Infra.ServicesName = name
	return b
}

// Build finalizes and returns the complete infrastructure
func (b *InfrastructureBuilder) Build() *InfrastructureContext {
	return b.Infra
}
