//go:build e2e

package fa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strings"
	"time"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/tests/controller"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

const (
	baseURL      = "http://localhost:8080/api/1"
	testPassword = "12345678"
	newPassword  = "123456789"
)

type E2E2FAFlow struct {
	suite.Suite
}

type GmailReader struct {
	email    string
	password string
}

func NewGmailReader(email, password string) *GmailReader {
	return &GmailReader{
		email:    email,
		password: password,
	}
}

func (g *GmailReader) extractVerificationCode(t provider.T, toEmail string) string {
	t.Logf("Connecting to Gmail IMAP...")

	// Подключаемся к Gmail IMAP
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		t.Fatalf("Failed to connect to Gmail: %v", err)
	}
	defer c.Logout()

	// Логинимся с App Password
	if err := c.Login(g.email, g.password); err != nil {
		t.Fatalf("Failed to login to Gmail: %v", err)
	}
	t.Log("Connected to Gmail")

	// Выбираем INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		t.Fatalf("Failed to select INBOX: %v", err)
	}
	t.Logf("INBOX selected, total messages: %d", mbox.Messages)

	// Ждем письмо
	maxAttempts := 2
	for i := 0; i < maxAttempts; i++ {
		time.Sleep(2 * time.Second)
		t.Logf("Checking for verification email... attempt %d/%d", i+1, maxAttempts)

		if mbox.Messages < 1 {
			continue
		}

		start := mbox.Messages - 1

		seqset := new(imap.SeqSet)
		seqset.AddRange(start, mbox.Messages)

		// СОЗДАЕМ СЕКЦИЮ ДЛЯ ТЕЛА ПИСЬМА
		// Для plain text писем используем "BODY[]"
		bodySection := &imap.BodySectionName{
			BodyPartName: imap.BodyPartName{
				Specifier: imap.TextSpecifier, // или можно использовать ""
			},
			Peek: true, // Не помечать как прочитанное
		}

		messages := make(chan *imap.Message, 5)
		done := make(chan error, 1)

		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{
				imap.FetchEnvelope,
				imap.FetchBodyStructure,
				bodySection.FetchItem(), // ЗАПРАШИВАЕМ ТЕЛО ЯВНО
			}, messages)
		}()

		if err := <-done; err != nil {
			t.Logf("Fetch error: %v", err)
			continue
		}

		for msg := range messages {
			if msg == nil || msg.Envelope == nil {
				continue
			}

			// Проверяем получателя
			for _, addr := range msg.Envelope.To {
				if addr.Address() == toEmail &&
					strings.Contains(msg.Envelope.Subject, "Ваш код подтверждения") {
					fmt.Println("Нашли на этот адрес:", addr.Address())

					t.Logf("Found verification email for %s", toEmail)

					// Получаем тело из правильной секции
					for section, literal := range msg.Body {
						t.Logf("Section: %v, size: %d bytes", section, literal.Len())

						// Читаем содержимое
						buf := new(bytes.Buffer)
						_, err := buf.ReadFrom(literal)
						if err != nil {
							t.Logf("Failed to read section: %v", err)
							continue
						}

						body := buf.String()
						t.Logf("Body content: %s", body)

						// Ищем код
						re := regexp.MustCompile(`Ваш код:\s*(\d{6})`)
						matches := re.FindStringSubmatch(body)
						if len(matches) > 1 {
							t.Logf("Verification code extracted: %s", matches[1])
							return matches[1]
						}
					}

					// Альтернативный способ: запрашиваем BODY[] без указания specifier
					t.Log("Trying alternative body section...")

					altSeqset := new(imap.SeqSet)
					altSeqset.AddNum(msg.SeqNum)

					altSection := &imap.BodySectionName{
						Peek: true,
					}

					altMessages := make(chan *imap.Message, 1)
					altDone := make(chan error, 1)

					go func() {
						altDone <- c.Fetch(altSeqset, []imap.FetchItem{
							altSection.FetchItem(),
						}, altMessages)
					}()

					if err := <-altDone; err != nil {
						t.Logf("Alt fetch error: %v", err)
						continue
					}

					altMsg := <-altMessages
					if altMsg != nil {
						for _, literal := range altMsg.Body {
							buf := new(bytes.Buffer)
							buf.ReadFrom(literal)
							body := buf.String()
							t.Logf("Alt body: %s", body)

							re := regexp.MustCompile(`Ваш код:\s*(\d{6})`)
							matches := re.FindStringSubmatch(body)
							if len(matches) > 1 {
								t.Logf("Verification code extracted (alt): %s", matches[1])
								return matches[1]
							}
						}
					}
				}
			}
		}

		// Обновляем количество сообщений в INBOX
		mbox, err = c.Select("INBOX", false)
		if err != nil {
			t.Logf("Failed to reselect INBOX: %v", err)
		}
	}

	t.Fatalf("Failed to extract verification code from email after %d attempts", maxAttempts)
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (f *E2E2FAFlow) getTestEmail() string {
	baseEmail := os.Getenv("TEST_GMAIL")
	if baseEmail == "" {
		baseEmail = "ilya.kovalev002@gmail.com"
	}

	username := baseEmail[:len(baseEmail)-10]
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s+e2e-%d@gmail.com", username, timestamp)
}

func (f *E2E2FAFlow) createTestUser(t provider.T) {
	client := &http.Client{}
	testEmail := f.getTestEmail()

	t.Logf("Creating test user with email: %s", testEmail)

	registerPayload := map[string]interface{}{
		"email":    testEmail,
		"password": testPassword,
		"name":     "e2e",
	}

	body, _ := json.Marshal(registerPayload)
	resp, err := client.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Logf("Register error: %v, trying login...", err)
		// Пробуем залогиниться - пользователь может уже существовать
		loginPayload := dto.LoginRequest{
			Login:    testEmail,
			Password: testPassword,
		}
		loginBody, _ := json.Marshal(loginPayload)
		loginResp, _ := client.Post(baseURL+"/auth", "application/json", bytes.NewReader(loginBody))
		if loginResp != nil {
			loginResp.Body.Close()
		}
	}
	if resp != nil {
		resp.Body.Close()
	}
}

func (f *E2E2FAFlow) Test2FALoginFlowWithGmail(t provider.T) {
	t.Title("2FA Login Flow with Real Gmail")
	t.Description("Test complete 2FA authentication flow using real Gmail IMAP")

	ctrl := controller.NewController(t)

	// Проверяем наличие Gmail credentials
	gmailAddress := os.Getenv("GMAIL_ADDRESS")
	if gmailAddress == "" {
		gmailAddress = "ilya.kovalev002@gmail.com"
	}
	gmailAppPassword := os.Getenv("GMAIL_APP_PASSWORD")

	if gmailAddress == "" || gmailAppPassword == "" {
		t.Skip("GMAIL_ADDRESS or GMAIL_APP_PASSWORD not set, skipping real Gmail test")
	}

	// Инициализируем Gmail reader
	gmailReader := NewGmailReader(gmailAddress, gmailAppPassword)

	t.WithNewStep("Test 2FA login flow with Gmail", func(ctx provider.StepCtx) {
		client := &http.Client{}
		jar, _ := cookiejar.New(nil)
		client.Jar = jar

		// Создаем уникальный email для этого теста
		testEmail := f.getTestEmail()
		ctx.Logf("Testing with email: %s", testEmail)

		// Step 1: Регистрируем пользователя с уникальным email
		ctx.Log("Step 1: Creating test user")
		registerPayload := dto.SignupRequest{
			Login:    testEmail,
			Password: testPassword,
			Role:     "buyer",
		}
		body, _ := json.Marshal(registerPayload)
		resp, err := client.Post(baseURL+"/users", "application/json", bytes.NewReader(body))
		ctx.Assert().NoError(err)
		ctx.Assert().Equal(http.StatusOK, resp.StatusCode)

		var signupResp dto.SignupResponse
		err = json.NewDecoder(resp.Body).Decode(&signupResp)
		ctx.Assert().NoError(err)

		newUserID := signupResp.Id

		resp.Body.Close()

		// Step 2: Получаем код подтверждения из Gmail
		ctx.Log("Step 3: Extracting verification code from Gmail")
		code := gmailReader.extractVerificationCode(t, testEmail)
		fmt.Println("CODE:", code)
		ctx.Assert().NotEmpty(code, "Verification code should not be empty")
		ctx.WithNewStep(fmt.Sprintf("Extracted code: %s", code), func(stepCtx provider.StepCtx) {})

		time.Sleep(time.Second)

		// Step 3: Verify 2FA code
		ctx.Log("Step 4: Verifying 2FA code")
		verifyPayload := dto.Verify2FARequest{
			Email: testEmail,
			Code:  code,
		}
		verifyBody, _ := json.Marshal(verifyPayload)

		resp, err = client.Post("http://localhost:8080/api/auth/verify-2fa", "application/json", bytes.NewReader(verifyBody))
		ctx.Assert().NoError(err)
		ctx.Assert().Equal(http.StatusOK, resp.StatusCode)

		var verifyResp dto.Verify2FAResponse

		err = json.NewDecoder(resp.Body).Decode(&verifyResp)
		resp.Body.Close()
		ctx.Assert().NoError(err)
		ctx.Assert().True(*verifyResp.Success, "2FA verification should succeed")
		ctx.Assert().NotEmpty(*verifyResp.Session, "Should receive session ID")

		time.Sleep(time.Second)

		// Step X: Change password
		changePasswordPayload := dto.UpdatePasswordRequest{
			Password:    testPassword,
			NewPassword: newPassword,
		}
		changeBody, _ := json.Marshal(changePasswordPayload)

		req, err := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("%s/users/%d/password", baseURL, newUserID),
			bytes.NewReader(changeBody),
		)
		ctx.Assert().NoError(err)

		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		ctx.Assert().NoError(err)
		defer resp.Body.Close()
		ctx.Assert().Equal(http.StatusOK, resp.StatusCode)

		// Step 4: login - should NOT require 2FA
		ctx.Log("Step 5: Testing re-login within grace period")
		time.Sleep(2 * time.Second)

		loginPayload := dto.LoginRequest{
			Login:    testEmail,
			Password: newPassword,
		}
		loginBody, _ := json.Marshal(loginPayload)

		resp, err = client.Post(baseURL+"/auth", "application/json", bytes.NewReader(loginBody))
		ctx.Assert().NoError(err)
		ctx.Assert().Equal(http.StatusOK, resp.StatusCode)

		var reloginResp string
		err = json.NewDecoder(resp.Body).Decode(&reloginResp)
		resp.Body.Close()
		ctx.Assert().NoError(err)
		ctx.Assert().NotEmpty(reloginResp, "2FA should not be empty")

		time.Sleep(time.Second)

		// Step 5: Delete user
		_, err = ctrl.GetDB().Exec(`delete from "user" where id = $1;`, newUserID)
		ctx.Assert().NoError(err)
	})
}
