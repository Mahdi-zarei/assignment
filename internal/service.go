package internal

import (
	"assignment/api"
	"assignment/config"
	"assignment/db/giftcard_repo"
	"assignment/db/users_repo"
	"assignment/internal/giftcard"
	"assignment/internal/users"
	"assignment/pkg/common"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Service struct {
	apiHandler *api.ApiHandler

	giftCardRepo giftcard_repo.GiftCardRepo
	usersRepo    users_repo.UsersRepo

	giftCardModule giftcard.GiftCardModule
	usersModule    users.UsersModule

	dbPool *pgxpool.Pool
	logger *logrus.Logger
}

func (s *Service) Start(ctx context.Context) {
	conf := config.GetConfig()
	s.logger = &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		},
		Level: logrus.DebugLevel,
	}

	s.dbPool = common.MustGetVal(pgxpool.New(ctx, conf.GiftShopPGXConf.GenerateConnectURL()))

	s.giftCardRepo = giftcard_repo.NewGiftCardRepo(ctx, s.dbPool)
	s.usersRepo = users_repo.NewUsersRepo(ctx, s.dbPool)

	s.giftCardModule = giftcard.NewGiftCardModule(s.giftCardRepo, s.usersRepo, s.logger)
	s.usersModule = users.NewUsersModule(s.usersRepo, s.logger)

	s.apiHandler = api.NewApiHandler(s.giftCardModule, s.usersModule, s.logger, conf.EndpointConf)
}

func (s *Service) Close() {
	// currently nothing
}
