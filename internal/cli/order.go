package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	usecase "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/order"
)

type orderCommand struct {
	order usecase.IOrderUsecase
}

func NewOrderCommand(o usecase.IOrderUsecase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order",
		Short: "Операции с заказами",
	}

	orderCmd := orderCommand{
		order: o,
	}

	cmd.AddCommand(
		orderCmd.getAllOrdersCmd(),
		orderCmd.createOrderCmd(),
		orderCmd.getOrderByIDCmd(),
		orderCmd.cancelOrderCmd(),
		orderCmd.updateOrderStatusCmd(),
	)

	return cmd
}

func (c *orderCommand) createOrderCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Создать новый заказ",
		Run: func(cmd *cobra.Command, args []string) {
			orderID, err := c.order.Create(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка создания: %v", err)
			}
			fmt.Printf("Заказ создан. ID: %d\n", orderID)
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *orderCommand) getOrderByIDCmd() *cobra.Command {
	var orderID int64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Посмотреть заказ по ID",
		Run: func(cmd *cobra.Command, args []string) {
			order, err := c.order.GetByID(context.Background(), orderID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}
			fmt.Printf("Заказ:\n%#v\n", order)
		},
	}

	cmd.Flags().Int64Var(&orderID, "orderID", orderID, "Order ID")

	_ = cmd.MarkFlagRequired("orderID")

	return cmd
}

func (c *orderCommand) getAllOrdersCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Получить список заказов пользователя",
		Run: func(cmd *cobra.Command, args []string) {
			orders, err := c.order.GetAllOrders(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}

			fmt.Println("Заказы:")
			for _, order := range orders {
				fmt.Printf("%s\n", formatOrder(order))
			}
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *orderCommand) updateOrderStatusCmd() *cobra.Command {
	var orderID int64
	var newStatus string

	cmd := &cobra.Command{
		Use:   "updatestatus",
		Short: "Обновить статус заказа",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.order.UpdateStatus(context.Background(), orderID, newStatus)
			if err != nil {
				log.Fatalf("Ошибка обновления: %v", err)
			}
			fmt.Println("Статус обновлен")
		},
	}

	cmd.Flags().Int64Var(&orderID, "orderID", orderID, "Order ID")
	cmd.Flags().StringVar(&newStatus, "status", newStatus, "Новое значение")

	_ = cmd.MarkFlagRequired("orderID")
	_ = cmd.MarkFlagRequired("status")

	return cmd
}

func (c *orderCommand) cancelOrderCmd() *cobra.Command {
	var orderID int64

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Отменить заказ",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.order.Cancel(context.Background(), orderID)
			if err != nil {
				log.Fatalf("Ошибка отмены: %v", err)
			}
			fmt.Println("Заказ успешно отменен")
		},
	}

	cmd.Flags().Int64Var(&orderID, "orderID", orderID, "Order ID")

	_ = cmd.MarkFlagRequired("orderID")

	return cmd
}
