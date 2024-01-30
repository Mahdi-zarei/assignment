package test

import (
	"assignment/config"
	"assignment/db/giftcard_repo"
	"assignment/db/users_repo"
	"assignment/pkg/common"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var service *ServiceTest

type ServiceTest struct {
	dbConn       *pgxpool.Pool
	giftCardRepo giftcard_repo.GiftCardRepo
	usersRepo    users_repo.UsersRepo
}

func NewServiceTest() *ServiceTest {
	if service != nil {
		service.cleanDBs()
		return service
	}

	conf := config.GetConfig()

	dbConn := common.MustGetVal(pgxpool.New(context.Background(), conf.GiftShopPGXConf.GenerateConnectURL()))

	service = &ServiceTest{
		dbConn:       dbConn,
		giftCardRepo: giftcard_repo.NewGiftCardRepo(context.Background(), dbConn),
		usersRepo:    nil,
	}
	service.cleanDBs() // make sure tables are clean
	return service
}

func (s *ServiceTest) cleanDBs() {
	tableNames := []string{"gift_card", "user"}
	for _, tableName := range tableNames {
		common.Must2(s.dbConn.Exec(context.Background(), "TRUNCATE TABLE $1", tableName))
	}
}
