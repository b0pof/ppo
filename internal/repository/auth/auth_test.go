package auth_test

//
//import (
//	"testing"
//	"time"
//
//	"github.com/alicebob/miniredis/v2"
//	"github.com/go-redis/redis"
//	"github.com/ozontech/allure-go/pkg/framework/provider"
//	"github.com/ozontech/allure-go/pkg/framework/suite"
//
//	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
//	. "git.iu7.bmstu.ru/kia22u475/ppo/internal/repository/auth"
//)
//
//type AuthRepositorySuite struct {
//	suite.Suite
//	redisServer *miniredis.Miniredis
//	redisClient *redis.Client
//}
//
//func TestAuthRepository(t *testing.T) {
//	suite.RunSuite(t, new(AuthRepositorySuite))
//}
//
//func (s *AuthRepositorySuite) BeforeEach(t provider.T) {
//	redisServer, err := miniredis.Run()
//	if err != nil {
//		t.Fatal("Failed to start miniredis:", err)
//	}
//	s.redisServer = redisServer
//
//	s.redisClient = redis.NewClient(&redis.Options{
//		Addr: s.redisServer.Addr(),
//	})
//}
//
//func (s *AuthRepositorySuite) AfterEach(t provider.T) {
//	if s.redisClient != nil {
//		s.redisClient.Close()
//	}
//	if s.redisServer != nil {
//		s.redisServer.Close()
//	}
//}
//
//func (s *AuthRepositorySuite) TestRepository_GetUserIDBySessionID(t provider.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name         string
//		prepare      func() string
//		expectations func(a provider.Asserts, got int64, err error)
//	}{
//		{
//			name: "success: session exists",
//			prepare: func() string {
//				sessionID := "session123"
//				s.redisServer.Set(sessionID, "42")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got int64, err error) {
//				a.NoError(err)
//				a.Equal(int64(42), got)
//			},
//		},
//		{
//			name: "error: session does not exist",
//			prepare: func() string {
//				return "non-existent-session"
//			},
//			expectations: func(a provider.Asserts, got int64, err error) {
//				a.Error(err)
//				a.ErrorIs(err, model.ErrNotFound)
//				a.Equal(int64(0), got)
//			},
//		},
//		{
//			name: "error: invalid session data",
//			prepare: func() string {
//				sessionID := "invalid-session"
//				s.redisServer.Set(sessionID, "not-a-number")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got int64, err error) {
//				a.Error(err)
//				a.ErrorIs(err, model.ErrNotFound)
//				a.Equal(int64(0), got)
//			},
//		},
//		{
//			name: "success: session with zero user ID",
//			prepare: func() string {
//				sessionID := "zero-session"
//				s.redisServer.Set(sessionID, "0")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got int64, err error) {
//				a.NoError(err)
//				a.Equal(int64(0), got)
//			},
//		},
//		{
//			name: "success: session with large user ID",
//			prepare: func() string {
//				sessionID := "large-session"
//				s.redisServer.Set(sessionID, "9999999999")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got int64, err error) {
//				a.NoError(err)
//				a.Equal(int64(9999999999), got)
//			},
//		},
//		{
//			name: "error: empty session ID",
//			prepare: func() string {
//				return ""
//			},
//			expectations: func(a provider.Asserts, got int64, err error) {
//				a.Error(err)
//				a.ErrorIs(err, model.ErrNotFound)
//				a.Equal(int64(0), got)
//			},
//		},
//	}
//
//	for _, tc := range tests {
//		tc := tc
//		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
//			repo := New(s.redisClient)
//
//			sessionID := tc.prepare()
//
//			got, err := repo.GetUserIDBySessionID(sessionID)
//
//			tc.expectations(ctx.Assert(), got, err)
//		})
//	}
//}
//
//func (s *AuthRepositorySuite) TestRepository_SessionExists(t provider.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name         string
//		prepare      func() string
//		expectations func(a provider.Asserts, got bool)
//	}{
//		{
//			name: "success: session exists with valid user ID",
//			prepare: func() string {
//				sessionID := "existing-session"
//				s.redisServer.Set(sessionID, "123")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got bool) {
//				a.True(got)
//			},
//		},
//		{
//			name: "false: session does not exist",
//			prepare: func() string {
//				return "non-existent-session"
//			},
//			expectations: func(a provider.Asserts, got bool) {
//				a.False(got)
//			},
//		},
//		{
//			name: "false: session exists but user ID is zero",
//			prepare: func() string {
//				sessionID := "zero-user-session"
//				s.redisServer.Set(sessionID, "0")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got bool) {
//				a.False(got)
//			},
//		},
//		{
//			name: "false: session exists but invalid data",
//			prepare: func() string {
//				sessionID := "invalid-data-session"
//				s.redisServer.Set(sessionID, "not-a-number")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got bool) {
//				a.False(got)
//			},
//		},
//		{
//			name: "false: empty session ID",
//			prepare: func() string {
//				return ""
//			},
//			expectations: func(a provider.Asserts, got bool) {
//				a.False(got)
//			},
//		},
//		{
//			name: "success: session exists with negative user ID",
//			prepare: func() string {
//				sessionID := "negative-session"
//				s.redisServer.Set(sessionID, "-1")
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got bool) {
//				a.True(got)
//			},
//		},
//		{
//			name: "false: session expired",
//			prepare: func() string {
//				sessionID := "expired-session"
//				s.redisServer.Set(sessionID, "123")
//				s.redisServer.Del(sessionID)
//				return sessionID
//			},
//			expectations: func(a provider.Asserts, got bool) {
//				a.False(got)
//			},
//		},
//	}
//
//	for _, tc := range tests {
//		tc := tc
//		t.WithNewStep(tc.name, func(ctx provider.StepCtx) {
//			repo := New(s.redisClient)
//
//			sessionID := tc.prepare()
//
//			got := repo.SessionExists(sessionID)
//
//			tc.expectations(ctx.Assert(), got)
//		})
//	}
//}

//func (s *AuthRepositorySuite) TestRepository_EdgeCases(t provider.T) {
//	t.Parallel()
//
//	t.WithNewStep("multiple sessions for different users", func(ctx provider.StepCtx) {
//		repo := New(s.redisClient)
//
//		sessions := map[string]string{
//			"sess1": "1",
//			"sess2": "2",
//			"sess3": "3",
//		}
//
//		for sessionID, userID := range sessions {
//			s.redisServer.Set(sessionID, userID)
//		}
//
//		for sessionID, expectedUserID := range sessions {
//			userID, err := repo.GetUserIDBySessionID(sessionID)
//			ctx.Assert().NoError(err)
//			ctx.Assert().Equal(expectedUserID, string(rune(userID)))
//
//			exists := repo.SessionExists(sessionID)
//			ctx.Assert().True(exists)
//		}
//
//		exists := repo.SessionExists("non-existent")
//		ctx.Assert().False(exists)
//	})
//
//	t.WithNewStep("special characters in session ID", func(ctx provider.StepCtx) {
//		repo := New(s.redisClient)
//
//		specialSessionID := "session-123@#$%"
//		s.redisServer.Set(specialSessionID, "99")
//
//		userID, err := repo.GetUserIDBySessionID(specialSessionID)
//		ctx.Assert().NoError(err)
//		ctx.Assert().Equal(int64(99), userID)
//
//		exists := repo.SessionExists(specialSessionID)
//		ctx.Assert().True(exists)
//	})
//
//	t.WithNewStep("concurrent access", func(ctx provider.StepCtx) {
//		repo := New(s.redisClient)
//
//		sessionID := "concurrent-session"
//		s.redisServer.Set(sessionID, "100")
//
//		done := make(chan bool, 2)
//
//		go func() {
//			userID, err := repo.GetUserIDBySessionID(sessionID)
//			ctx.Assert().NoError(err)
//			ctx.Assert().Equal(int64(100), userID)
//			done <- true
//		}()
//
//		go func() {
//			exists := repo.SessionExists(sessionID)
//			ctx.Assert().True(exists)
//			done <- true
//		}()
//
//		<-done
//		<-done
//	})
//}
//
//func (s *AuthRepositorySuite) TestRepository_RedisOperations(t provider.T) {
//	t.Parallel()
//
//	t.WithNewStep("set and get operations", func(ctx provider.StepCtx) {
//		repo := New(s.redisClient)
//
//		sessionID := "test-session"
//		userID := int64(555)
//
//		err := s.redisClient.Set(sessionID, userID, time.Hour).Err()
//		ctx.Assert().NoError(err)
//
//		retrievedID, err := repo.GetUserIDBySessionID(sessionID)
//		ctx.Assert().NoError(err)
//		ctx.Assert().Equal(userID, retrievedID)
//
//		exists := repo.SessionExists(sessionID)
//		ctx.Assert().True(exists)
//	})
//
//	t.WithNewStep("delete operation", func(ctx provider.StepCtx) {
//		repo := New(s.redisClient)
//
//		sessionID := "to-delete-session"
//		s.redisServer.Set(sessionID, "777")
//
//		exists := repo.SessionExists(sessionID)
//		ctx.Assert().True(exists)
//
//		s.redisServer.Del(sessionID)
//
//		exists = repo.SessionExists(sessionID)
//		ctx.Assert().False(exists)
//
//		_, err := repo.GetUserIDBySessionID(sessionID)
//		ctx.Assert().Error(err)
//		ctx.Assert().ErrorIs(err, model.ErrNotFound)
//	})
//}
