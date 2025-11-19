package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/b0pof/ppo/internal/usecase/cart"
)

type cartCommand struct {
	cart cart.ICartUsecase
}

func NewCartCommand(c cart.ICartUsecase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cart",
		Short: "Операции с корзиной",
	}

	cartCmd := cartCommand{
		cart: c,
	}

	cmd.AddCommand(
		cartCmd.clearCartCmd(),
		cartCmd.getCartItemsAmountCmd(),
		cartCmd.getItemsByCartCmd(),
		cartCmd.addItemCmd(),
		cartCmd.deleteItemCmd(),
	)

	return cmd
}

func (c *cartCommand) getCartItemsAmountCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "content",
		Short: "Количество товаров в корзине",
		Run: func(cmd *cobra.Command, args []string) {
			count, err := c.cart.GetCartItemsAmount(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка подсчета товаров: %v", err)
			}
			fmt.Printf("Товаров в корзине: %d\n", count)
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *cartCommand) getItemsByCartCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "cartitems",
		Short: "Посмотреть товары в корзине",
		Run: func(cmd *cobra.Command, args []string) {
			content, err := c.cart.GetCartContent(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}
			fmt.Println("Товары в корзине:")
			for _, item := range content.Items {
				fmt.Printf("- %+v\n", item)
			}
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *cartCommand) addItemCmd() *cobra.Command {
	var userID, itemID int64

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Добавить товар в корзину",
		Run: func(cmd *cobra.Command, args []string) {
			count, err := c.cart.AddItem(context.Background(), userID, itemID)
			if err != nil {
				log.Fatalf("Ошибка добавления: %v", err)
			}
			fmt.Printf("Успешно, количество обновлено: %d\n", count)
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")
	cmd.Flags().Int64Var(&itemID, "itemID", itemID, "Item ID")

	_ = cmd.MarkFlagRequired("userID")
	_ = cmd.MarkFlagRequired("itemID")

	return cmd
}

func (c *cartCommand) deleteItemCmd() *cobra.Command {
	var userID, itemID int64

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Удалить товар из корзины",
		Run: func(cmd *cobra.Command, args []string) {
			count, err := c.cart.DeleteItem(context.Background(), userID, itemID)
			if err != nil {
				log.Fatalf("Ошибка удаления: %v", err)
			}
			fmt.Printf("Успешно, количество обновлено: %d\n", count)
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")
	cmd.Flags().Int64Var(&itemID, "itemID", itemID, "Item ID")

	_ = cmd.MarkFlagRequired("userID")
	_ = cmd.MarkFlagRequired("itemID")

	return cmd
}

func (c *cartCommand) clearCartCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Очистить корзину",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.cart.Clear(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка очистки: %v", err)
			}
			fmt.Println("Успешно")
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}
