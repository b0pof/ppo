package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/cohere-ai/cohere-go/v2/client"

	dto "github.com/b0pof/ppo/internal/generated/clients"
)

type Usecase struct {
	client *client.Client
}

func New(c *client.Client) *Usecase {
	return &Usecase{
		client: c,
	}
}

func (u *Usecase) Generate(productName string) (string, error) {
	requestMessage := fmt.Sprintf("Придумай и напиши описание для товара с названием \"%s\" одним предложением", productName)
	resp, err := u.client.Chat(
		context.TODO(),
		&cohere.ChatRequest{
			Message: requestMessage,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to send chat request: %w", err)
	}

	rawResp, _ := json.Marshal(resp)

	var respDTO dto.ChatResponse
	err = json.Unmarshal(rawResp, &respDTO)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal chat response: %w", err)
	}

	rawDescription := strings.Trim(respDTO.Text, "\"")

	return rawDescription, nil
}
