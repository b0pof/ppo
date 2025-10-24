//go:build integration

package repository

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestSuite(t *testing.T) {
	suite.RunSuite(t, new(ItemFlow))
	suite.RunSuite(t, new(CartFlow))
	suite.RunSuite(t, new(SessionCreateDelete))
	suite.RunSuite(t, new(OrderFlow))
	suite.RunSuite(t, new(UserFlow))
}
