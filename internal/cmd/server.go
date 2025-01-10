package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

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
		wg := &sync.WaitGroup{}
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		dbConn := mysql.InitDBConn()
		e := echo.New()

		accountRepository := repository.InitAccountRepository(dbConn)
		accountUSecase := usecase.NewAccountUsecase(accountRepository)
		accountHandler := httphandler.InitAccountHandler(accountUSecase)

		accountHandler.RegisterAccountHandler(e)

		wg.Add(2)
		go startHTTPServer(e, quit, wg)
		

		wg.Wait()
		log.Println("All servers stopped gracefully")
	},
}

func startHTTPServer(e *echo.Echo, quit chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong!")
	})

	if err := e.Start(":" + config.PORT()); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown HTTP server:", err)
	}
}

func init() {
	rootCmd.AddCommand(startServe)
}
