package scheduler

import (
	"context"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/vote/internal/usecases"
	"time"
)

type ElectionSyncScheduler struct {
	electionTimeUc usecases.ElectionTimeUseCases
	ticker         *time.Ticker
	done           chan bool
}

func NewElectionSyncScheduler(electionTimeUc usecases.ElectionTimeUseCases) *ElectionSyncScheduler {
	return &ElectionSyncScheduler{
		electionTimeUc: electionTimeUc,
		done:           make(chan bool),
	}
}

func (s *ElectionSyncScheduler) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)

	log.WithFields(log.Fields{
		"interval": interval.String(),
	}).Info("[ElectionTimeSyncScheduler] Starting election time sync scheduler")

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.syncElections()
			case <-s.done:
				log.Info("[ElectionTimeSyncScheduler] Stopping election time sync scheduler")
				return
			}
		}
	}()
}

func (s *ElectionSyncScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
}

func (s *ElectionSyncScheduler) syncElections() {
	ctx := context.Background()

	log.Debug("[ElectionTimeSyncScheduler] Syncing elections")

	err := s.electionTimeUc.SyncElectionStatuses(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("[ElectionTimeSyncScheduler] Failed to sync elections")
		return
	}

	log.Debug("[ElectionTimeSyncScheduler] Successfully synced elections")
}

func (s *ElectionSyncScheduler) SyncOnce() error {
	ctx := context.Background()

	log.Info("[ElectionTimeSyncScheduler] Syncing elections once")

	err := s.electionTimeUc.SyncElectionStatuses(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("[ElectionTimeSyncScheduler] Failed to sync elections")
		return err
	}

	log.Info("[ElectionTimeSyncScheduler] Successfully synced elections")
	return nil
}
