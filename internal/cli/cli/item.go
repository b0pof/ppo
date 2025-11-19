package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/b0pof/ppo/internal/model"
	usecase "github.com/b0pof/ppo/internal/usecase/item"
)

type itemCommand struct {
	item usecase.IItemUsecase
}

func NewItemCommand(i usecase.IItemUsecase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item",
		Short: "Операции с товарами",
	}

	itemCmd := itemCommand{
		item: i,
	}

	cmd.AddCommand(
		itemCmd.getItemByIDCmd(),
		itemCmd.getAllItemsCmd(),
		itemCmd.createItemCmd(),
		itemCmd.getItemsBySellerIDCmd(),
		itemCmd.updateItemCmd(),
		itemCmd.deleteItemCmd(),
	)

	return cmd
}

func (c *itemCommand) createItemCmd() *cobra.Command {
	var name, description, imgSrc string
	var price int
	var sellerID int64

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Создать новый товар",
		Run: func(cmd *cobra.Command, args []string) {
			newItem := model.Item{
				Name:        name,
				Description: description,
				Price:       price,
				Seller: model.Seller{
					ID: sellerID,
				},
				ImgSrc: imgSrc,
			}
			id, err := c.item.Create(context.Background(), newItem)
			if err != nil {
				log.Fatalf("Ошибка создания: %v", err)
			}
			fmt.Printf("Успешно создан. ID: %d\n", id)
		},
	}

	cmd.Flags().StringVar(&name, "name", name, "Название товара")
	cmd.Flags().StringVar(&description, "desc", description, "Описание товара")
	cmd.Flags().StringVar(&imgSrc, "imgSrc", imgSrc, "Ссылка на картинку товара")
	cmd.Flags().IntVar(&price, "price", price, "Цена товара")
	cmd.Flags().Int64Var(&sellerID, "sellerID", sellerID, "Seller ID")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("price")
	_ = cmd.MarkFlagRequired("sellerID")

	return cmd
}

func (c *itemCommand) getItemByIDCmd() *cobra.Command {
	var userID, itemID int64

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Получить товар по ID",
		Run: func(cmd *cobra.Command, args []string) {
			item, err := c.item.GetByID(context.Background(), userID, itemID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}
			fmt.Printf("Товар:\n%s\n", formatItemExtended(item))
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")
	cmd.Flags().Int64Var(&itemID, "itemID", itemID, "Item ID")

	_ = cmd.MarkFlagRequired("userID")
	_ = cmd.MarkFlagRequired("itemID")

	return cmd
}

func (c *itemCommand) getAllItemsCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Получить список всех товаров",
		Run: func(cmd *cobra.Command, args []string) {
			items, err := c.item.GetAllItems(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}

			fmt.Print("Товары:\n")
			for _, item := range items {
				fmt.Printf("%s\n", formatItemExtended(item))
			}
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *itemCommand) getItemsBySellerIDCmd() *cobra.Command {
	var sellerID int64

	cmd := &cobra.Command{
		Use:   "getbysellerid",
		Short: "Получить список товаров продавца",
		Run: func(cmd *cobra.Command, args []string) {
			items, err := c.item.GetItemsBySellerID(context.Background(), sellerID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}

			fmt.Printf("Товары селлера с ID = %d:\n", sellerID)
			for _, item := range items {
				fmt.Printf("%s\n", formatItem(item))
			}
		},
	}

	cmd.Flags().Int64Var(&sellerID, "sellerID", sellerID, "Seller ID")

	_ = cmd.MarkFlagRequired("sellerID")

	return cmd
}

func (c *itemCommand) deleteItemCmd() *cobra.Command {
	var itemID int64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Удалить товар по id",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.item.Delete(context.Background(), itemID)
			if err != nil {
				log.Fatalf("Ошибка удаления: %v", err)
			}
			fmt.Println("Товар удален")
		},
	}

	cmd.Flags().Int64Var(&itemID, "itemID", itemID, "Item ID")

	_ = cmd.MarkFlagRequired("itemID")

	return cmd
}

func (c *itemCommand) updateItemCmd() *cobra.Command {
	var id int64
	var name, description, imgSrc string
	var price int

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Обновить информацию о товаре",
		Run: func(cmd *cobra.Command, args []string) {
			updateItem := model.Item{
				ID:          id,
				Name:        name,
				Description: description,
				Price:       price,
				ImgSrc:      imgSrc,
			}
			err := c.item.Update(context.Background(), updateItem)
			if err != nil {
				log.Fatalf("Ошибка обновления: %v", err)
			}
			fmt.Println("Товар обновлен")
		},
	}

	cmd.Flags().Int64Var(&id, "id", id, "Item ID")
	cmd.Flags().StringVar(&name, "name", name, "Название товара")
	cmd.Flags().StringVar(&imgSrc, "imgSrc", imgSrc, "Ссылка на картинку товара")
	cmd.Flags().StringVar(&description, "desc", description, "Описание товара")
	cmd.Flags().IntVar(&price, "price", price, "Цена товара")

	_ = cmd.MarkFlagRequired("id")

	return cmd
}
