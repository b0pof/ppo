package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/usecase/auth"
)

type authCommand struct {
	auth auth.IAuthUsecase
}

func NewAuthCommand(a auth.IAuthUsecase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Авторизация",
	}

	authCmd := authCommand{
		auth: a,
	}

	cmd.AddCommand(
		authCmd.isLoggedInCmd(),
		authCmd.loginCmd(),
		authCmd.logoutCmd(),
		authCmd.signupCmd(),
		authCmd.getUserIDBySessionCmd(),
	)

	return cmd
}

func (c *authCommand) loginCmd() *cobra.Command {
	var login, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Войти в аккаунт",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := c.auth.Login(context.Background(), login, password)
			if err != nil {
				log.Fatalf("Ошибка входа: %v", err)
			}
			fmt.Printf("Session ID: %s\n", token)
		},
	}

	cmd.Flags().StringVar(&login, "login", login, "Логин")
	cmd.Flags().StringVar(&password, "password", password, "Пароль")

	_ = cmd.MarkFlagRequired("login")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}

func (c *authCommand) signupCmd() *cobra.Command {
	var login, password, role string

	cmd := &cobra.Command{
		Use:   "signup",
		Short: "Регистрация",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := c.auth.Signup(context.Background(), login, password, role)
			if err != nil {
				log.Fatalf("Ошибка регистрации: %v", err)
			}
			fmt.Printf("Session ID: %s\n", token)
		},
	}

	cmd.Flags().StringVar(&login, "login", login, "Логин")
	cmd.Flags().StringVar(&password, "password", password, "Пароль для нового аккаунта")
	cmd.Flags().StringVar(&role, "role", role, "Роль нового аккаунта (buyer, seller)")

	_ = cmd.MarkFlagRequired("login")
	_ = cmd.MarkFlagRequired("password")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

func (c *authCommand) logoutCmd() *cobra.Command {
	var sessionID string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Выйти из аккаунта",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.auth.Logout(sessionID)
			if err != nil {
				log.Fatalf("Ошибка выхода: %v", err)
			}
			fmt.Println("Успешно!")
		},
	}

	cmd.Flags().StringVar(&sessionID, "session", sessionID, "Session ID")

	_ = cmd.MarkFlagRequired("session")

	return cmd
}

func (c *authCommand) isLoggedInCmd() *cobra.Command {
	var sessionID string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Проверить, выполнен ли вход",
		Run: func(cmd *cobra.Command, args []string) {
			valid := c.auth.IsLoggedIn(sessionID)
			if valid {
				fmt.Println("ОК")
			} else {
				fmt.Println("Выход НЕ выполнен")
			}
		},
	}

	cmd.Flags().StringVar(&sessionID, "session", sessionID, "Session ID")

	_ = cmd.MarkFlagRequired("session")

	return cmd
}

func (c *authCommand) getUserIDBySessionCmd() *cobra.Command {
	var sessionID string

	cmd := &cobra.Command{
		Use:   "getsessionuser",
		Short: "Get user ID by session ID",
		Run: func(cmd *cobra.Command, args []string) {
			userID, err := c.auth.GetUserIDBySessionID(sessionID)
			if err != nil {
				log.Fatalf("Ошибка получения ID: %v", err)
			}
			fmt.Printf("User ID: %d\n", userID)
		},
	}

	cmd.Flags().StringVar(&sessionID, "session", sessionID, "Session ID")

	_ = cmd.MarkFlagRequired("session")

	return cmd
}
