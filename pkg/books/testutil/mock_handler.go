package testutil

import (
	"testing"

	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/nkitlabs/go-http-gorm-example/pkg/books/service"
)

type TestSuite struct {
	Handler    *service.Handler
	Service    *service.Service
	Repository *MockDataProvider
	Logger     *zap.Logger
}

func NewTestSuite(t *testing.T) TestSuite {
	ctrl := gomock.NewController(t)

	repo := NewMockDataProvider(ctrl)
	log := zap.NewNop()
	serv := service.NewService(repo, log)
	handler := service.NewHandler(&serv, log)

	return TestSuite{
		Handler:    &handler,
		Service:    &serv,
		Repository: repo,
		Logger:     log,
	}
}
