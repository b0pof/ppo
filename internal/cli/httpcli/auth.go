package httpcli

import (
	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/spf13/cobra"
)

var (
	authLoginFlag    string
	authPasswordFlag string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Команды авторизации",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Войти (POST /api/1/auth)",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.LoginRequest{
			Login:    authLoginFlag,
			Password: authPasswordFlag,
		}
		res, err := DoJSON("POST", "/api/1/auth", req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

func init() {
	// register auth group
	rootCmd.AddCommand(authCmd)

	// login
	authCmd.AddCommand(authLoginCmd)
	authLoginCmd.Flags().StringVar(&authLoginFlag, "login", "", "Логин (email)")
	authLoginCmd.Flags().StringVar(&authPasswordFlag, "password", "", "Пароль")
	authLoginCmd.MarkFlagRequired("login")
	authLoginCmd.MarkFlagRequired("password")
}
