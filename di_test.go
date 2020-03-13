package di

import (
	"context"
	"testing"
	"github.com/jackc/pgx/v4/pgxpool"
)

func newPool() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), "postgres://user:pass@localhost:5434/app")
	if err != nil {
		panic(err)
	}
	return pool
}

func TestCreateDependencies(t *testing.T){

	conainer := NewContainer()
	conainer.Provide(
		newPool,
	)

	if len(conainer.components) != 1 {
		t.Error("Can't create dependencies")
	}
}


func TestErrorDependencies(t *testing.T) {
	container := NewContainer()

	container.Provide(
		newPool,
	)
}