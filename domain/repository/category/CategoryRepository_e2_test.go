//go:build e2e
// +build e2e

package category

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
	"microservice/cmd/infra/config"
	"microservice/cmd/infra/database"
	"microservice/cmd/infra/log"
	"microservice/domain/entitie"
	"reflect"
	"testing"
)

func TestUnit_categoryRepository_GetByID(t *testing.T) {

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "aula",
			"POSTGRES_PASSWORD": "aula",
			"POSTGRES_DB":       "aula",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("PostgreSQL init process complete; ready for start up."),
			wait.ForListeningPort("5432/tcp"),
		),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := postgresC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	db := database.ProvideDatabase(log.ProvideLogger(), config.Config{
		Database: struct {
			Host     string
			Port     string
			User     string
			Name     string
			Password string
			SslMode  string
		}{Host: host, Port: port.Port(), User: "aula", Name: "aula", Password: "aula", SslMode: "disable"},
	})

	db.AutoMigrate(&entitie.Category{})

	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entitie.Category
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				DB: db,
			},
			args: args{
				id:  1,
				ctx: ctx,
			},
			want:    entitie.Category{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &categoryRepository{
				DB: tt.fields.DB,
			}
			got, err := m.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
