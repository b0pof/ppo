package httpcli

import (
	"fmt"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/spf13/cobra"
)

var itemCmd = &cobra.Command{
	Use:   "item",
	Short: "Работа с товарами",
}

// list
var itemListPage int
var itemListLimit int
var itemListCmd = &cobra.Command{
	Use:   "list",
	Short: "Список товаров (GET /api/1/items)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/items?page=%d&limit=%d", itemListPage, itemListLimit)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// get
var itemGetID int64
var itemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Получить товар по ID (GET /api/1/items/{id})",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/items/%d", itemGetID)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// create
var (
	itemCreateName        string
	itemCreateDescription string
	itemCreatePrice       int
	itemCreateImgSrc      string
)
var itemCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать товар (POST /api/1/items)",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.CreateItemRequest{
			Name:        itemCreateName,
			Description: itemCreateDescription,
			Price:       itemCreatePrice,
			ImgSrc:      itemCreateImgSrc,
		}
		res, err := DoJSON("POST", "/api/1/items", req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// update
var (
	itemUpdateID          int64
	itemUpdateName        string
	itemUpdateDescription string
	itemUpdatePrice       int
	itemUpdateImgSrc      string
)
var itemUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Обновить товар (PUT /api/1/items/{id})",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.UpdateItemRequest{
			Id:          itemUpdateID,
			Name:        itemUpdateName,
			Description: itemUpdateDescription,
			Price:       itemUpdatePrice,
			ImgSrc:      itemUpdateImgSrc,
		}
		path := fmt.Sprintf("/api/1/items/%d", itemUpdateID)
		res, err := DoJSON("PUT", path, req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// delete
var itemDeleteID int64
var itemDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Удалить товар (DELETE /api/1/items/{id})",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/items/%d", itemDeleteID)
		res, err := DoJSON("DELETE", path, nil)
		if err != nil {
			return err
		}
		// API returns 204 maybe with empty body — print status-friendly message
		if len(res) == 0 {
			fmt.Println("deleted")
		} else {
			printJSON(res)
		}
		return nil
	},
}

// reviews list & create
var (
	itemReviewsID int64
	reviewRating  int
	reviewAdv     string
	reviewDis     string
	reviewComm    string
)
var itemReviewsListCmd = &cobra.Command{
	Use:   "reviews",
	Short: "Список отзывов / Создать отзыв (subcommands: list/create)",
}
var itemReviewsListSub = &cobra.Command{
	Use:   "list",
	Short: "Получить отзывы (GET /api/1/items/{id}/reviews)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/items/%d/reviews", itemReviewsID)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}
var itemReviewsCreateSub = &cobra.Command{
	Use:   "create",
	Short: "Создать отзыв (POST /api/1/items/{id}/reviews)",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.CreateReviewRequest{
			ItemId:        itemReviewsID,
			Rating:        reviewRating,
			Advantages:    reviewAdv,
			Disadvantages: reviewDis,
			Comment:       reviewComm,
		}
		path := fmt.Sprintf("/api/1/items/%d/reviews", itemReviewsID)
		res, err := DoJSON("POST", path, req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(itemCmd)

	// list
	itemCmd.AddCommand(itemListCmd)
	itemListCmd.Flags().IntVar(&itemListPage, "page", 1, "Номер страницы")
	itemListCmd.Flags().IntVar(&itemListLimit, "limit", 20, "Размер страницы")

	// get
	itemCmd.AddCommand(itemGetCmd)
	itemGetCmd.Flags().Int64Var(&itemGetID, "id", 0, "ID товара")
	itemGetCmd.MarkFlagRequired("id")

	// create
	itemCmd.AddCommand(itemCreateCmd)
	itemCreateCmd.Flags().StringVar(&itemCreateName, "name", "", "Название")
	itemCreateCmd.Flags().StringVar(&itemCreateDescription, "description", "", "Описание")
	itemCreateCmd.Flags().IntVar(&itemCreatePrice, "price", 0, "Цена (в копейках/центах)")
	itemCreateCmd.Flags().StringVar(&itemCreateImgSrc, "img", "", "Путь к изображению")
	itemCreateCmd.MarkFlagRequired("name")
	itemCreateCmd.MarkFlagRequired("description")
	itemCreateCmd.MarkFlagRequired("price")

	// update
	itemCmd.AddCommand(itemUpdateCmd)
	itemUpdateCmd.Flags().Int64Var(&itemUpdateID, "id", 0, "ID товара")
	itemUpdateCmd.Flags().StringVar(&itemUpdateName, "name", "", "Название")
	itemUpdateCmd.Flags().StringVar(&itemUpdateDescription, "description", "", "Описание")
	itemUpdateCmd.Flags().IntVar(&itemUpdatePrice, "price", 0, "Цена")
	itemUpdateCmd.Flags().StringVar(&itemUpdateImgSrc, "img", "", "Путь к изображению")
	itemUpdateCmd.MarkFlagRequired("id")

	// delete
	itemCmd.AddCommand(itemDeleteCmd)
	itemDeleteCmd.Flags().Int64Var(&itemDeleteID, "id", 0, "ID товара")
	itemDeleteCmd.MarkFlagRequired("id")

	// reviews
	itemCmd.AddCommand(itemReviewsListCmd)
	itemReviewsListCmd.AddCommand(itemReviewsListSub)
	itemReviewsListCmd.AddCommand(itemReviewsCreateSub)
	itemReviewsListSub.Flags().Int64Var(&itemReviewsID, "id", 0, "ID товара")
	itemReviewsListSub.MarkFlagRequired("id")

	itemReviewsCreateSub.Flags().Int64Var(&itemReviewsID, "id", 0, "ID товара")
	itemReviewsCreateSub.Flags().IntVar(&reviewRating, "rating", 5, "Оценка 1..5")
	itemReviewsCreateSub.Flags().StringVar(&reviewAdv, "advantages", "", "Достоинства")
	itemReviewsCreateSub.Flags().StringVar(&reviewDis, "disadvantages", "", "Недостатки")
	itemReviewsCreateSub.Flags().StringVar(&reviewComm, "comment", "", "Комментарий")
	itemReviewsCreateSub.MarkFlagRequired("id")
	itemReviewsCreateSub.MarkFlagRequired("rating")
}
