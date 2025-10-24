package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/model"
	usecase "git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/user"
)

type userCommand struct {
	user usecase.IUserUsecase
}

func NewUserCommand(u usecase.IUserUsecase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Операции с пользователем",
	}

	userCmd := userCommand{
		user: u,
	}

	cmd.AddCommand(
		userCmd.getUserRoleByIDCmd(),
		userCmd.getUserMetaInfoByUserIDCmd(),
		userCmd.getUserByIDCmd(),
		userCmd.updatePasswordCmd(),
		userCmd.updateUserByIDCmd(),
	)

	return cmd
}

func (c *userCommand) getUserByIDCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "getbyid",
		Short: "Профиль пользователя по ID",
		Run: func(cmd *cobra.Command, args []string) {
			user, err := c.user.GetByID(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}
			fmt.Printf("Пользователь:\n%s\n", formatUser(user))
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *userCommand) getUserRoleByIDCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "getrole",
		Short: "Посмотреть роль пользователя",
		Run: func(cmd *cobra.Command, args []string) {
			role, err := c.user.GetRoleByID(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}
			fmt.Printf("Роль: %s\n", role)
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *userCommand) getUserMetaInfoByUserIDCmd() *cobra.Command {
	var userID int64

	cmd := &cobra.Command{
		Use:   "getmetainfo",
		Short: "Получить мета-информацию профиля",
		Run: func(cmd *cobra.Command, args []string) {
			metaInfo, err := c.user.GetMetaInfoByUserID(context.Background(), userID)
			if err != nil {
				log.Fatalf("Ошибка поиска: %v", err)
			}
			fmt.Printf("Мета-информация: %+v\n", metaInfo)
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *userCommand) updateUserByIDCmd() *cobra.Command {
	var userID int64
	var name, login, phone string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Обновить профиль по ID",
		Run: func(cmd *cobra.Command, args []string) {
			user := model.User{
				Name:  name,
				Login: login,
				Phone: phone,
			}
			err := c.user.UpdateByID(context.Background(), userID, user)
			if err != nil {
				log.Fatalf("Ошибка обновления: %v", err)
			}
			fmt.Println("Успешно обновлен")
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")
	cmd.Flags().StringVar(&name, "name", name, "Имя")
	cmd.Flags().StringVar(&login, "login", login, "Логин")
	cmd.Flags().StringVar(&phone, "phone", phone, "Номер телефона")

	_ = cmd.MarkFlagRequired("userID")

	return cmd
}

func (c *userCommand) updatePasswordCmd() *cobra.Command {
	var userID int64
	var oldPassword, newPassword string

	cmd := &cobra.Command{
		Use:   "updatepassword",
		Short: "Сменить пароль",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.user.UpdatePassword(context.Background(), userID, oldPassword, newPassword)
			if err != nil {
				log.Fatalf("Ошибка обновления: %v", err)
			}
			fmt.Println("Успешно обновлен")
		},
	}

	cmd.Flags().Int64Var(&userID, "userID", userID, "User ID")
	cmd.Flags().StringVar(&oldPassword, "old", oldPassword, "Старый пароль")
	cmd.Flags().StringVar(&newPassword, "new", newPassword, "Новый пароль")

	_ = cmd.MarkFlagRequired("userID")
	_ = cmd.MarkFlagRequired("old")
	_ = cmd.MarkFlagRequired("new")

	return cmd
}
