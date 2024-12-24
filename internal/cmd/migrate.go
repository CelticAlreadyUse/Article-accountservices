package cmd

import (
	"database/sql"
	"log"
	"time"

	dbmysql "github.com/CelticAlreadyUse/Article-accountservices/internal/databases/mysql"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)
var (
	direction = "up"
	step      = 1
)
var migrationCmd = &cobra.Command{
	Use: "migrate",
	Short: "migrate is a commend to start migrate a databases",
	Run: func (cmd *cobra.Command ,args []string)  {
		dbConn, err := sql.Open("mysql", dbmysql.InitConnStr())

		if err != nil {
			panic(err)
		}

		dbConn.SetConnMaxLifetime(time.Minute * 3)
		dbConn.SetMaxOpenConns(10)
		dbConn.SetMaxIdleConns(10)

		if err := dbConn.Ping(); err != nil {
			panic(err)
		}
		if len(args) > 0 {
			direction = args[0]
		}

		switch direction {
		case "up", "down":
		default:
			log.Fatal("direction is not valid")
		}

		migrations := &migrate.FileMigrationSource{Dir: "db/migrations"}
		var migrationCount int
		if direction == "down" {
			migrationCount, err = migrate.ExecMax(dbConn, "mysql", migrations, migrate.Down, step)
		} else {
			migrationCount, err = migrate.ExecMax(dbConn, "mysql", migrations, migrate.Up, step)
		}

		if err != nil {
			log.Fatalf("failed to run migration, error: %s", err.Error())
		}

		log.Printf("successfully applied %d migration(s)", migrationCount)
	},

	}

	func init() {
		migrationCmd.Flags().IntVarP(&step, "step", "s", 1, "number of migration step")
		rootCmd.AddCommand(migrationCmd)
	}
	