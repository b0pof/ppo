//go:build e2e

package scenario

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestSuite(t *testing.T) {
	suite.RunSuite(t, new(CreateOrderFlow))
}
