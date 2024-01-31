package test

import (
	"assignment/config"
	"assignment/db/giftcard_repo"
	"assignment/db/users_repo"
	"assignment/pkg/common"
	"context"
	"fmt"
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
		usersRepo:    users_repo.NewUsersRepo(context.Background(), dbConn),
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
