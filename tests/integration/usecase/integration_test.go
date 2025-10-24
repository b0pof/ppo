//go:build integration

package usecase

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestSuite(t *testing.T) {
	suite.RunSuite(t, new(ItemFlow))
	suite.RunSuite(t, new(OrderFlow))
	suite.RunSuite(t, new(UserTest))
	suite.RunSuite(t, new(CartTest))
	suite.RunSuite(t, new(AuthTest))
}
