package httpcli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Категории и связанные операции",
}

// get category
var categoryGetID int64
var categoryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Получить категорию (GET /api/1/categories/{id})",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/categories/%d", categoryGetID)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// items by category
var categoryItemsID int64
var categoryItemsPage int
var categoryItemsLimit int
var categoryItemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Получить товары категории (GET /api/1/categories/{id}/items)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/categories/%d/items?page=%d&limit=%d", categoryItemsID, categoryItemsPage, categoryItemsLimit)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(categoryCmd)

	categoryCmd.AddCommand(categoryGetCmd)
	categoryGetCmd.Flags().Int64Var(&categoryGetID, "id", 0, "ID категории")
	categoryGetCmd.MarkFlagRequired("id")

	categoryCmd.AddCommand(categoryItemsCmd)
	categoryItemsCmd.Flags().Int64Var(&categoryItemsID, "id", 0, "ID категории")
	categoryItemsCmd.Flags().IntVar(&categoryItemsPage, "page", 1, "Страница")
	categoryItemsCmd.Flags().IntVar(&categoryItemsLimit, "limit", 20, "Лимит")
	categoryItemsCmd.MarkFlagRequired("id")
}
