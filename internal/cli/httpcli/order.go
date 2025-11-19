package httpcli

import (
	"fmt"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Работа с заказами",
}

// list
var orderListUser int64
var orderListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список заказов пользователя (GET /api/1/users/{id}/orders)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/users/%d/orders", orderListUser)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// create (from cart)
var orderCreateUser int64
var orderCreateStatus string
var orderCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать заказ из корзины (POST /api/1/users/{id}/orders)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// swagger/sdk suggests body type UpdateOrderRequest for POST too; it's optional.
		var body any = nil
		if orderCreateStatus != "" {
			body = dto.UpdateOrderRequest{Status: orderCreateStatus}
		}
		path := fmt.Sprintf("/api/1/users/%d/orders", orderCreateUser)
		res, err := DoJSON("POST", path, body)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// get by id
var orderGetUser int64
var orderGetID int64
var orderGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Получить заказ (GET /api/1/users/{userId}/orders/{orderId})",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/users/%d/orders/%d", orderGetUser, orderGetID)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// update status
var orderUpdateUser int64
var orderUpdateID int64
var orderUpdateStatus string
var orderUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Обновить статус заказа (PATCH /api/1/users/{userId}/orders/{orderId})",
	RunE: func(cmd *cobra.Command, args []string) error {
		body := dto.UpdateOrderRequest{Status: orderUpdateStatus}
		path := fmt.Sprintf("/api/1/users/%d/orders/%d", orderUpdateUser, orderUpdateID)
		res, err := DoJSON("PATCH", path, body)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(orderCmd)

	orderCmd.AddCommand(orderListCmd)
	orderListCmd.Flags().Int64Var(&orderListUser, "user", 0, "ID пользователя")
	orderListCmd.MarkFlagRequired("user")

	orderCmd.AddCommand(orderCreateCmd)
	orderCreateCmd.Flags().Int64Var(&orderCreateUser, "user", 0, "ID пользователя")
	orderCreateCmd.Flags().StringVar(&orderCreateStatus, "status", "", "Статус (опционально)")
	orderCreateCmd.MarkFlagRequired("user")

	orderCmd.AddCommand(orderGetCmd)
	orderGetCmd.Flags().Int64Var(&orderGetUser, "user", 0, "ID пользователя")
	orderGetCmd.Flags().Int64Var(&orderGetID, "order", 0, "ID заказа")
	orderGetCmd.MarkFlagRequired("user")
	orderGetCmd.MarkFlagRequired("order")

	orderCmd.AddCommand(orderUpdateCmd)
	orderUpdateCmd.Flags().Int64Var(&orderUpdateUser, "user", 0, "ID пользователя")
	orderUpdateCmd.Flags().Int64Var(&orderUpdateID, "order", 0, "ID заказа")
	orderUpdateCmd.Flags().StringVar(&orderUpdateStatus, "status", "", "Новый статус")
	orderUpdateCmd.MarkFlagRequired("user")
	orderUpdateCmd.MarkFlagRequired("order")
	orderUpdateCmd.MarkFlagRequired("status")
}
