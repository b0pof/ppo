////go:build e2e

package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/b0pof/ppo/tests/controller"
)

type E2ELLMDescriptionFlow struct {
	suite.Suite
}

const (
	baseUURL = "http://localhost:8080/api/1"
)

func (o *E2ELLMDescriptionFlow) TestLLMDescriptionGeneration(t provider.T) {
	_ = controller.NewController(t)

	t.WithNewStep("TestLLMDescriptionGeneration", func(ctxA provider.StepCtx) {
		client := &http.Client{}
		jar, _ := cookiejar.New(nil)
		client.Jar = jar

		// 1. Авторизация
		dollNames := []string{"Стейси розовая", "Майк синяя", "Бетмен черная"}
		testItemName := fmt.Sprintf("Кукла %s", dollNames[rand.Intn(len(dollNames))])

		// 2. Создаем айтем с пустым описанием
		createReq := dto.CreateItemRequest{
			Name:        testItemName,
			Description: "",
			Price:       6666,
			ImgSrc:      "https://example.com/images/doll.jpg",
		}

		body, _ := json.Marshal(createReq)
		resp, err := client.Post(baseUURL+"/items", "application/json", bytes.NewReader(body))
		t.Assert().NoError(err, "Error while create item")

		var createResp dto.CreateItemResponse
		t.Assert().NoError(json.NewDecoder(resp.Body).Decode(&createResp))
		resp.Body.Close()
		t.Assert().NotZero(createResp.ItemId)

		respBodyCreate, _ := io.ReadAll(resp.Body)
		ctxA.Logf("Created item with ID: %d and empty description, resp: %s\n", createResp.ItemId, string(respBodyCreate))

		time.Sleep(10 * time.Second)

		// 3. Получаем айтем по ID и проверяем описание
		getResp, err := client.Get(fmt.Sprintf("%s/items/%d", baseUURL, createResp.ItemId))
		t.Assert().NoError(err)
		t.Assert().Equal(http.StatusOK, getResp.StatusCode)

		var item dto.Item
		t.Assert().NoError(json.NewDecoder(getResp.Body).Decode(&item))
		getResp.Body.Close()

		t.Assert().NotEmpty(item.Description)
		t.Assert().Greater(len(item.Description), 3)

		t.Logf("Item description: %q\n", item.Description)
		t.Logf("Description: %s\n", item.Description)
		t.Logf("Description length: %d characters\n", len(item.Description))
	})
}
