package test

import (
	"assignment/api"
	"assignment/config"
	"assignment/db/giftcard_repo"
	"assignment/db/users_repo"
	"assignment/internal/giftcard"
	"assignment/internal/users"
	"assignment/pkg/common"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"os"
)

var service *ServiceTest

type ServiceTest struct {
	dbConn       *pgxpool.Pool
	giftCardRepo giftcard_repo.GiftCardRepo
	usersRepo    users_repo.UsersRepo

	giftCardModule giftcard.GiftCardModule
	userModule     users.UsersModule

	apiHandler *api.ApiHandler
}

func NewServiceTest() *ServiceTest {
	if service != nil {
		service.cleanDBs()
		return service
	}

	conf := config.GetConfig()
	logger := &logrus.Logger{
		Out: os.Stdout,
	}

	dbConn := common.MustGetVal(pgxpool.New(context.Background(), conf.GiftShopPGXConf.GenerateConnectURL()))
	giftcardRepo := giftcard_repo.NewGiftCardRepo(context.Background(), dbConn)
	usersRepo := users_repo.NewUsersRepo(context.Background(), dbConn)

	giftCardModule := giftcard.NewGiftCardModule(giftcardRepo, usersRepo, logger)
	userModule := users.NewUsersModule(usersRepo, logger)

	apiHandler := api.NewApiHandler(giftCardModule, userModule, logger, conf.EndpointConf)

	service = &ServiceTest{
		dbConn:         dbConn,
		giftCardRepo:   giftcardRepo,
		usersRepo:      usersRepo,
		giftCardModule: giftCardModule,
		userModule:     userModule,
		apiHandler:     apiHandler,
	}
	service.cleanDBs() // make sure tables are clean
	return service
}

func (s *ServiceTest) cleanDBs() {
	tableNames := []string{"gift_card", "users"} // list of table names to clear
	for _, tableName := range tableNames {
		common.Must2(s.dbConn.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s", tableName)))
	}
}
