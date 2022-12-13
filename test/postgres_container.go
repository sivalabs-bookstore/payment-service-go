package test

import (
	"context"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

const postgresImage = "postgres:15.0-alpine"
const postgresPort = "5432"
const postgresUserName = "postgres"
const postgresPassword = "postgres"
const postgresDbName = "postgres"

type PostgresContainer struct {
	Container testcontainers.Container
	CloseFn   func()
	Host      string
	Port      string
	Database  string
	Username  string
	Password  string
}

// SetupPostgres creates an instance of the postgres container type
func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	port, err := nat.NewPort("tcp", postgresPort)
	if err != nil {
		return nil, err
	}
	req := testcontainers.ContainerRequest{
		Image: postgresImage,
		Env: map[string]string{
			"POSTGRES_USER":     postgresUserName,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDbName,
		},
		ExposedPorts: []string{port.Port()},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second)).
			WithDeadline(1 * time.Minute),
		Cmd: []string{"postgres", "-c", "fsync=off"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, _ := container.Host(ctx)
	hostPort, _ := container.MappedPort(ctx, postgresPort)

	return &PostgresContainer{
		Container: container,
		CloseFn: func() {
			log.Info("terminating container")
			if err := container.Terminate(ctx); err != nil {
				log.Fatalf("error terminating postgres container: %s", err)
			}
		},
		Host:     host,
		Port:     hostPort.Port(),
		Database: postgresDbName,
		Username: postgresUserName,
		Password: postgresPassword,
	}, nil
}
