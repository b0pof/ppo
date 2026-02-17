//go:build e2e

package fa

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestSuite(t *testing.T) {
	suite.RunSuite(t, new(E2ELLMDescriptionFlow))
}
