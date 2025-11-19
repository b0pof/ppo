package httpcli

import (
	"fmt"

	dto "github.com/b0pof/ppo/internal/generated"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Работа с пользователями",
}

// signup
var (
	suLogin    string
	suPassword string
	suRole     string
)
var userSignupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Регистрация (POST /api/1/users)",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.SignupRequest{
			Login:    suLogin,
			Password: suPassword,
			Role:     dto.SignupRequestRole(suRole),
		}
		res, err := DoJSON("POST", "/api/1/users", req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// get
var userGetID int64
var userGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Получить профиль пользователя (GET /api/1/users/{id})",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/users/%d", userGetID)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// update
var (
	userUpdateID    int64
	userUpdateName  string
	userUpdateLogin string
	userUpdatePhone string
)
var userUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Обновить профиль (PUT /api/1/users/{id})",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.UpdateUserRequest{
			Name:  userUpdateName,
			Login: userUpdateLogin,
			Phone: userUpdatePhone,
		}
		path := fmt.Sprintf("/api/1/users/%d", userUpdateID)
		res, err := DoJSON("PUT", path, req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// password
var (
	userPasswordID  int64
	userPasswordOld string
	userPasswordNew string
)
var userPasswordCmd = &cobra.Command{
	Use:   "password",
	Short: "Обновить пароль (PATCH /api/1/users/{id}/password)",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := dto.UpdatePasswordRequest{
			Password:    userPasswordOld,
			NewPassword: userPasswordNew,
		}
		path := fmt.Sprintf("/api/1/users/%d/password", userPasswordID)
		res, err := DoJSON("PATCH", path, req)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

// meta (user meta)
var userMetaID int64
var userMetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Получить метаданные пользователя (GET /api/1/users/{id}/meta)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := fmt.Sprintf("/api/1/users/%d/meta", userMetaID)
		res, err := DoJSON("GET", path, nil)
		if err != nil {
			return err
		}
		printJSON(res)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// signup
	userCmd.AddCommand(userSignupCmd)
	userSignupCmd.Flags().StringVar(&suLogin, "login", "", "Логин")
	userSignupCmd.Flags().StringVar(&suPassword, "password", "", "Пароль")
	userSignupCmd.Flags().StringVar(&suRole, "role", "buyer", "Роль: buyer|seller|admin")
	userSignupCmd.MarkFlagRequired("login")
	userSignupCmd.MarkFlagRequired("password")
	userSignupCmd.MarkFlagRequired("role")

	// get
	userCmd.AddCommand(userGetCmd)
	userGetCmd.Flags().Int64Var(&userGetID, "id", 0, "ID пользователя")
	userGetCmd.MarkFlagRequired("id")

	// update
	userCmd.AddCommand(userUpdateCmd)
	userUpdateCmd.Flags().Int64Var(&userUpdateID, "id", 0, "ID пользователя")
	userUpdateCmd.Flags().StringVar(&userUpdateName, "name", "", "Имя")
	userUpdateCmd.Flags().StringVar(&userUpdateLogin, "login", "", "Логин")
	userUpdateCmd.Flags().StringVar(&userUpdatePhone, "phone", "", "Телефон")
	userUpdateCmd.MarkFlagRequired("id")

	// password
	userCmd.AddCommand(userPasswordCmd)
	userPasswordCmd.Flags().Int64Var(&userPasswordID, "id", 0, "ID пользователя")
	userPasswordCmd.Flags().StringVar(&userPasswordOld, "old", "", "Старый пароль")
	userPasswordCmd.Flags().StringVar(&userPasswordNew, "new", "", "Новый пароль")
	userPasswordCmd.MarkFlagRequired("id")
	userPasswordCmd.MarkFlagRequired("old")
	userPasswordCmd.MarkFlagRequired("new")

	// meta
	userCmd.AddCommand(userMetaCmd)
	userMetaCmd.Flags().Int64Var(&userMetaID, "id", 0, "ID пользователя")
	userMetaCmd.MarkFlagRequired("id")
}
