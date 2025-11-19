package httpcli

import (
	"fmt"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/spf13/cobra"
)

var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Работа с корзиной пользователя",
}

// get
var cartGetUser int64
var cartGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Получить корзину пользователя (GET /api/1/users/{id}/cart/items)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/users/%d/cart/items", cartGetUser)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// add
var cartAddUser int64
var cartAddItem int64
var cartAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавить товар в корзину (POST /api/1/users/{id}/cart/items)",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.AddItemRequest{
			ItemId: cartAddItem,
		}
		path := fmt.Sprintf("/api/1/users/%d/cart/items", cartAddUser)
		res, err := DoJSON("POST", path, req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// clear
var cartClearUser int64
var cartClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Очистить корзину (DELETE /api/1/users/{id}/cart/items)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/users/%d/cart/items", cartClearUser)
		res, err := DoJSON("DELETE", path, nil)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			fmt.Println("cart cleared")
		} else {
			printJSON(res)
		}
		return nil
	},
}

// remove single item
var cartRemoveUser int64
var cartRemoveItem int64
var cartRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Удалить позицию из корзины (DELETE /api/1/users/{userId}/cart/items/{itemId})",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/users/%d/cart/items/%d", cartRemoveUser, cartRemoveItem)
		res, err := DoJSON("DELETE", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cartCmd)

	cartCmd.AddCommand(cartGetCmd)
	cartGetCmd.Flags().Int64Var(&cartGetUser, "user", 0, "ID пользователя")
	cartGetCmd.MarkFlagRequired("user")

	cartCmd.AddCommand(cartAddCmd)
	cartAddCmd.Flags().Int64Var(&cartAddUser, "user", 0, "ID пользователя")
	cartAddCmd.Flags().Int64Var(&cartAddItem, "item", 0, "ID товара")
	cartAddCmd.MarkFlagRequired("user")
	cartAddCmd.MarkFlagRequired("item")

	cartCmd.AddCommand(cartClearCmd)
	cartClearCmd.Flags().Int64Var(&cartClearUser, "user", 0, "ID пользователя")
	cartClearCmd.MarkFlagRequired("user")

	cartCmd.AddCommand(cartRemoveCmd)
	cartRemoveCmd.Flags().Int64Var(&cartRemoveUser, "user", 0, "ID пользователя")
	cartRemoveCmd.Flags().Int64Var(&cartRemoveItem, "item", 0, "ID товара")
	cartRemoveCmd.MarkFlagRequired("user")
	cartRemoveCmd.MarkFlagRequired("item")
}
