package database

import (
	"database/sql"
	"testing"

	"github.com.br/lucas-negrello/fc-ms-wallet/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDB *ClientDB
}

func (suite *ClientDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.Nil(err)
	suite.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	suite.clientDB = NewClientDB(db)
}

func (suite *ClientDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (suite *ClientDBTestSuite) TestSave() {
	client := &entity.Client{
		ID:    "id",
		Name:  "name",
		Email: "email",
	}
	err := suite.clientDB.Save(client)
	suite.Nil(err)
}

func (suite *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "j@j")
	suite.clientDB.Save(client)

	clientDB, err := suite.clientDB.Get(client.ID)
	suite.Nil(err)
	suite.NotNil(clientDB)
	suite.Equal(client.ID, clientDB.ID)
	suite.Equal(client.Name, clientDB.Name)
	suite.Equal(client.Email, clientDB.Email)
}
