package cmd

import (
	"net/http"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/databases/mysql"
	httphandler "github.com/CelticAlreadyUse/Article-accountservices/internal/dellivery/http"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/repository"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)
var startServe = &cobra.Command{
	Use:   "serve",
	Short: "serve is a command to run http serve",
	Run: func(cmd *cobra.Command, args []string) {
		dbConn := mysql.InitDBConn()
		e := echo.New()
		accountRepository := repository.InitAccountRepository(dbConn)
		accountUSecase := usecase.NewAccountUsecase(accountRepository)
		accountHandler := httphandler.InitAccountHandler(accountUSecase)
		e.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong!")
		})
		accountHandler.RegisterAccountHandler(e)
		e.Start(":" + config.PORT())
	},
}
func init() {
	rootCmd.AddCommand(startServe)
}
