package test

import (
	"context"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer represents the postgres container type used in the module
type PostgresContainer struct {
	testcontainers.Container
}

type postgresContainerOption func(req *testcontainers.ContainerRequest)

func WithWaitStrategy(strategies ...wait.Strategy) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.WaitingFor = wait.ForAll(strategies...).WithDeadline(1 * time.Minute)
	}
}

func WithPort(port string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.ExposedPorts = append(req.ExposedPorts, port)
	}
}

func WithInitialDatabase(user string, password string, dbName string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = password
		req.Env["POSTGRES_DB"] = dbName
	}
}

// SetupPostgres creates an instance of the postgres container type
func SetupPostgres(ctx context.Context, opts ...postgresContainerOption) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:11-alpine",
		Env:          map[string]string{},
		ExposedPorts: []string{},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
	}

	for _, opt := range opts {
		opt(&req)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{Container: container}, nil
}

const PostgresTestUserName = "test"
const PostgresTestPassword = "test"
const PostgresTestDatabase = "test"

func SetupTestDatabase(ctx context.Context) (testcontainers.Container, func(), error) {

	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return nil, nil, err
	}

	container, err := SetupPostgres(ctx,
		WithPort(port.Port()),
		WithInitialDatabase(PostgresTestUserName, PostgresTestPassword, PostgresTestDatabase),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	closeContainer := func() {
		log.Info("terminating container")
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("error terminating postgres container: %s", err)
		}
	}

	return container, closeContainer, nil
}
