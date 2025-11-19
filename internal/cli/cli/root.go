package cli

import (
	"github.com/spf13/cobra"

	"github.com/b0pof/ppo/internal/usecase/auth"
	"github.com/b0pof/ppo/internal/usecase/cart"
	"github.com/b0pof/ppo/internal/usecase/item"
	"github.com/b0pof/ppo/internal/usecase/order"
	"github.com/b0pof/ppo/internal/usecase/user"
)

type CLI struct {
	rootCmd *cobra.Command
}

func New(u user.IUserUsecase, a auth.IAuthUsecase, c cart.ICartUsecase, o order.IOrderUsecase, i item.IItemUsecase) *CLI {
	cli := &CLI{}

	rootCmd := &cobra.Command{
		Use:   "shop",
		Short: "CLI для функционала маркетплейса",
	}

	rootCmd.AddCommand(
		NewUserCommand(u),
		NewAuthCommand(a),
		NewCartCommand(c),
		NewOrderCommand(o),
		NewItemCommand(i),
	)

	cli.rootCmd = rootCmd

	return cli
}

func (c *CLI) Execute() error {
	return c.rootCmd.Execute()
}
