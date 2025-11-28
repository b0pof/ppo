package repository_test

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestSuite(t *testing.T) {
	suite.RunSuite(t, new(RepoItemFlow))
	suite.RunSuite(t, new(RepoCartFlow))
	suite.RunSuite(t, new(RepoSessionCreateDelete))
	suite.RunSuite(t, new(RepoOrderFlow))
	suite.RunSuite(t, new(RepoUserFlow))
}
