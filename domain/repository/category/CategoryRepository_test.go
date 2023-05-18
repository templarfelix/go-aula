package category

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"microservice/domain/entitie"
	"regexp"
	"testing"
	"time"
)

type RepositorySuite struct {
	suite.Suite
	ctx        context.Context
	conn       *sql.DB
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository *categoryRepository
	entitie    *entitie.Category
}

func (rs *RepositorySuite) SetupSuite() {
	var (
		err error
	)
	rs.conn, rs.mock, err = sqlmock.New()
	assert.NoError(rs.T(), err)
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 rs.conn,
		PreferSimpleProtocol: true,
	})
	rs.DB, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(rs.T(), err)

	rs.repository = &categoryRepository{
		DB: rs.DB,
	}
	assert.IsType(rs.T(), &categoryRepository{}, rs.repository)
	rs.entitie = &entitie.Category{
		Name: "Test",
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rs.ctx = context.Background()
}

func (rs *RepositorySuite) AfterTest(_, _ string) {
	assert.NoError(rs.T(), rs.mock.ExpectationsWereMet())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (rs *RepositorySuite) TestE2ECategoryRepository_Store() {

	rs.mock.ExpectBegin()
	rs.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "categories" ("name","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(
			rs.entitie.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	rs.mock.ExpectCommit()

	err := rs.repository.Store(rs.ctx, rs.entitie)

	assert.NoError(rs.T(), err) // valida se houve algum erro

}
