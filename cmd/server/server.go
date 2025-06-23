package server

import (
	"context"
	"github.com/nocturna-ta/golib/cache"
	_ "github.com/nocturna-ta/golib/cache/redis"
	"github.com/nocturna-ta/golib/cache/redlock"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/utils/encryption"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/handler/api"
	"github.com/nocturna-ta/vote/internal/infrastructures/ethereum"
	"github.com/nocturna-ta/vote/internal/infrastructures/kafka"
	"github.com/nocturna-ta/vote/internal/infrastructures/sms"
	"github.com/nocturna-ta/vote/internal/scheduler"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	redLock := redlock.New(&redlock.Config{
		ConnectionUrl: cfg.Redis.Connection,
	})

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

	redis, err := cache.New(cfg.Redis.Connection)
	if err != nil {
		log.Fatalf("Failed to instantiate redis client: %v", err)
	}

	smsProvider, err := sms.NewSMSProvider(cfg.SMS)
	if err != nil {
		log.Fatalf("Failed to instantiate SMS provider: %v", err)
	}

	encryptor, err := encryption.NewEncryption(cfg.Encryption.Key)
	if err != nil {
		log.Fatalf("Failed to instantiate encryption service: %v", err)
	}

	appContainer := newContainer(&options{
		Cfg:       cfg,
		DB:        database,
		Client:    client,
		Publisher: publisher,
		RedLock:   redLock,
		Cache:     redis,
		SMS:       smsProvider,
		Encryptor: encryptor,
	})

	electionScheduler := scheduler.NewElectionSyncScheduler(appContainer.ElectionTimeUc)
	electionScheduler.Start(1 * time.Minute)

	log.Info("Performing initial election status sync...")
	err = electionScheduler.SyncOnce()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Initial election sync failed, but continuing...")
	}

	server := api.New(&api.Options{
		Cfg:            appContainer.Cfg,
		VoteUc:         appContainer.VoteUc,
		ElectionTimeUc: appContainer.ElectionTimeUc,
		OtpUc:          appContainer.OtpUc,
	})

	go server.Run()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	log.Info("Server started successfully. Press Ctrl+C to shutdown.")

	select {
	case <-term:
		log.Info("Exiting gracefully...")
		electionScheduler.Stop()
		log.Info("Election sync scheduler stopped.")
	case err := <-server.ListenError():
		log.Error("Error starting web server, exiting gracefully:", err)

		electionScheduler.Stop()
	}

	return nil
}
