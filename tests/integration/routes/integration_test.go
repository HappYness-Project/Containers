package tests

import (
	"fmt"
	"log"
	"net/http/httptest"

	"github.com/happYness-Project/taskManagementGolang/cmd/api"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

func runTestServer() *httptest.Server {
	env := configs.InitConfig("integration_test.env")
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	logger := loggers.Setup(env)
	logger.Info().Msg(connStr)
	database, err := dbs.ConnectToDb(connStr)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewApiServer(fmt.Sprintf(":%d", 8000), env.AccessTokenSecret, database, logger)
	return httptest.NewServer(server.Setup())
}
