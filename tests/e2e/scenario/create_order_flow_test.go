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

type CreateOrderFlow struct {
	suite.Suite
}

const baseURL = "http://localhost:8080/api/1"

func (o *CreateOrderFlow) TestOrderCreationFlow(t provider.T) {
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
		fmt.Println("RESP:", resp)

		var lr loginResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&lr))
		resp.Body.Close()
		ctxA.Assert().NotEmpty(lr.SessionID)
		ctxA.Assert().Equal(int64(1), lr.UserID)

		addItem := map[string]int64{"itemId": 3} // Lego set
		body, _ = json.Marshal(addItem)
		resp, err = client.Post(fmt.Sprintf("%s/users/%d/cart/items", baseURL, lr.UserID),
			"application/json", bytes.NewReader(body))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		resp, err = client.Get(fmt.Sprintf("%s/users/%d/cart/items", baseURL, lr.UserID))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var cart cartItemsResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&cart))
		resp.Body.Close()
		ctxA.Assert().GreaterOrEqual(cart.TotalCount, 1)
		ctxA.Assert().NotEmpty(cart.Items)
		ctxA.Assert().Equal(int64(3), cart.Items[0].ID)

		orderReq := map[string]string{"status": "created"}
		body, _ = json.Marshal(orderReq)
		resp, err = client.Post(fmt.Sprintf("%s/users/%d/orders", baseURL, lr.UserID),
			"application/json", bytes.NewReader(body))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var orderResp createOrderResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&orderResp))
		resp.Body.Close()
		ctxA.Assert().NotZero(orderResp.OrderID)

		resp, err = client.Get(fmt.Sprintf("%s/users/%d/cart/items", baseURL, lr.UserID))
		ctxA.Assert().NoError(err)
		ctxA.Assert().Equal(http.StatusOK, resp.StatusCode)

		var emptyCart cartItemsResponse
		ctxA.Assert().NoError(json.NewDecoder(resp.Body).Decode(&emptyCart))
		resp.Body.Close()
		ctxA.Assert().Equal(0, emptyCart.TotalCount)

		resp, err = client.Get(fmt.Sprintf("%s/users/%d/orders/%d", baseURL, lr.UserID, orderResp.OrderID))
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
