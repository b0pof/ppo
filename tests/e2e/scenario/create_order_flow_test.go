//go:build e2e

package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/b0pof/ppo/tests/controller"
)

type E2ECreateOrderFlow struct {
	suite.Suite
}

const (
	baseURL = "http://localhost:8080/api/1"
	userID  = 1
)

func (o *E2ECreateOrderFlow) TestOrderCreationFlow(t provider.T) {
	_ = controller.NewController(t)

	t.WithNewStep("TestOrderCreationFlow", func(ctxA provider.StepCtx) {
		client := &http.Client{}
		jar, _ := cookiejar.New(nil)
		client.Jar = jar

		loginPayload := dto.LoginRequest{
			Login:    "user1",
			Password: "testtest",
		}
		body, _ := json.Marshal(loginPayload)
		resp, err := client.Post(baseURL+"/auth", "application/json", bytes.NewReader(body))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var sessionID string
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&sessionID))
		resp.Body.Close()
		ctxA.Assert().NotEmpty(sessionID)

		addItem := map[string]int64{"itemId": 3} // Lego set
		body, _ = json.Marshal(addItem)
		resp, err = client.Post(fmt.Sprintf("%s/users/%d/cart/items", baseURL, userID),
			"application/json", bytes.NewReader(body))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		resp, err = client.Get(fmt.Sprintf("%s/users/%d/cart/items", baseURL, userID))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var cart cartItemsResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&cart))
		resp.Body.Close()
		ctxA.Assert().Equal(cart.TotalCount, 5)
		ctxA.Assert().Equal(cart.TotalCount, 5)
		ctxA.Assert().Equal(cart.TotalPrice, 6199)
		ctxA.Assert().NotEmpty(cart.Items)
		ctxA.Assert().Contains(cart.Items, cartItem{
			ID:    1,
			Name:  "Doll",
			Price: 1000,
			Count: 3,
		})
		ctxA.Assert().Contains(cart.Items, cartItem{
			ID:    2,
			Name:  "Plastic car",
			Price: 1200,
			Count: 1,
		})
		ctxA.Assert().Contains(cart.Items, cartItem{
			ID:    3,
			Name:  "Lego set",
			Price: 1999,
			Count: 1,
		})

		orderReq := map[string]string{"status": "created"}
		body, _ = json.Marshal(orderReq)
		resp, err = client.Post(fmt.Sprintf("%s/users/%d/orders", baseURL, userID),
			"application/json", bytes.NewReader(body))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var orderResp createOrderResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&orderResp))
		resp.Body.Close()
		ctxA.Assert().NotZero(orderResp.OrderID)

		resp, err = client.Get(fmt.Sprintf("%s/users/%d/cart/items", baseURL, userID))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var emptyCart cartItemsResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&emptyCart))
		resp.Body.Close()
		ctxA.Assert().Equal(0, emptyCart.TotalCount)

		resp, err = client.Get(fmt.Sprintf("%s/users/%d/orders/%d", baseURL, userID, orderResp.OrderID))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var order orderResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&order))
		resp.Body.Close()

		ctxA.Assert().Equal(orderResp.OrderID, order.ID)
		ctxA.Assert().Equal("created", order.Status)
		ctxA.Assert().Greater(order.ItemsCount, 0)
		ctxA.Assert().NotEmpty(order.Items)
	})
}
