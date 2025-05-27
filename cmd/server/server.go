package server

import (
	"context"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/handler/api"
	"github.com/nocturna-ta/vote/internal/infrastructures/ethereum"
	"github.com/nocturna-ta/vote/internal/infrastructures/kafka"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var (
	serverHTTPCmd = &cobra.Command{
		Use:   "server-http",
		Short: "Voting Service HTTP",
		Long:  "Voting Service HTTP",
		RunE:  run,
	}
)

func ServeHttpCmd() *cobra.Command {
	serverHTTPCmd.Flags().StringP("config", "c", "", "Config Path, both relative or absolute. i.e: /usr/local/bin/config/files")
	return serverHTTPCmd
}

func run(cmd *cobra.Command, args []string) error {
	configLocation, _ := cmd.Flags().GetString("config")
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, configLocation)

	database := sql.New(sql.DBConfig{
		SlaveDSN:        cfg.Database.SlaveDSN,
		MasterDSN:       cfg.Database.MasterDSN,
		RetryInterval:   cfg.Database.RetryInterval,
		MaxIdleConn:     cfg.Database.MaxIdleConn,
		MaxConn:         cfg.Database.MaxConn,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}, sql.DriverPostgres)

	client, err := ethereum.GetEthereumClient(&cfg.Blockchain)
	if err != nil {
		return err
	}

	defer client.Close()

	publisher, err := kafka.NewPublisher(context.Background(), cfg.Kafka.Producer)
	if err != nil {
		log.Fatalf("Failed to instantiate kafka publisher: %w", err)
		return err
	}

	appContainer := newContainer(&options{
		Cfg:       cfg,
		Client:    client,
		DB:        database,
		Publisher: publisher,
	})

	server := api.New(&api.Options{
		Cfg:    appContainer.Cfg,
		VoteUc: appContainer.VoteUc,
	})

	go server.Run()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Info("Exiting gracefully...")
	case err := <-server.ListenError():
		log.Error("Error starting web server, exiting gracefully:", err)
	}

	return nil
}
