package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/databases/mysql"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/databases/redis"
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
		redisConn := redis.InitRedisClient()
		ctx := context.Background()
		accountRepository := repository.InitAccountRepository(dbConn)
		accountUSecase := usecase.NewAccountUsecase(accountRepository)
		otpRepository := repository.NewOTPRepository(redisConn, ctx)
		otpUsecase := usecase.InitUsecaseOTP(otpRepository, accountRepository)
		accountHandler := httphandler.InitAccountHandler(accountUSecase, otpUsecase)
		accountHandler.RegisterAccountHandler(e)
		startHTTPServer(e)
	},
}

func startHTTPServer(e *echo.Echo) {

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong!")
	})

	if err := e.Start(":" + config.PORT()); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}

}

func init() {
	rootCmd.AddCommand(startServe)
}
