package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	event2 "github.com/nocturna-ta/common-model/models/event"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/domain/repository"
	"github.com/nocturna-ta/vote/internal/infrastructures/cache"
	"github.com/nocturna-ta/votechain-contract/binding/electionManager"
	"github.com/nocturna-ta/votechain-contract/interfaces"
	"time"
)

type VoteProcessor struct {
	voteRepo      repository.VoteRepository
	contract      interfaces.ElectionManagerInterface
	redisCache    cache.VoteCache
	redisLock     cache.DistributedLock
	consumerGroup string
	client        ethereum.Client
}

type VoteProcessorOptions struct {
	VoteRepo      repository.VoteRepository
	ContractAddr  common.Address
	Contract      interfaces.ElectionManagerInterface
	RedisCache    cache.VoteCache
	RedisLock     cache.DistributedLock
	Client        ethereum.Client
	ConsumerGroup string
}

func NewVoteProcessor(opts VoteProcessorOptions) *VoteProcessor {
	var contractInterface interfaces.ElectionManagerInterface
	contract, err := electionManager.NewElectionManager(opts.ContractAddr, opts.Client.GetEthClient())
	if err != nil {
		return nil
	}
	contractInterface = contract
	return &VoteProcessor{
		voteRepo:      opts.VoteRepo,
		contract:      contractInterface,
		redisCache:    opts.RedisCache,
		redisLock:     opts.RedisLock,
		client:        opts.Client,
		consumerGroup: opts.ConsumerGroup,
	}
}

func (p *VoteProcessor) Start(ctx context.Context, cfg *event.ConsumerConfig) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteProcessor.Start")
	defer span.End()

	consumer, err := event.NewConsumer(ctx, cfg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[VoteProcessor.Start] Failed to create consumer")
		return err
	}

	err = consumer.Subscribe(ctx, "votes", p.consumerGroup, p.han)
}

func (p *VoteProcessor) handleVoteMessages(ctx context.Context, message *event.EventConsumeMessage) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteProcessor.handleVoteMessages")
	defer span.End()

	var voteEvent event2.VoteMessage
	err := json.Unmarshal(message.Data, &voteEvent)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(message.Data),
		}).ErrorWithCtx(ctx, "[VoteProcessor.handleVoteMessages] Failed to unmarshal message")
		return err
	}

	lockKey := fmt.Sprintf("vote-processor:%s", voteEvent.VoterID)
	err = p.redisLock.AcquireLockWithTTL(ctx, lockKey, 30*time.Second, 3)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"voterID": voteEvent.VoterID,
		}).ErrorWithCtx(ctx, "[VoteProcessor.handleVoteMessages] Failed to acquire lock")
		return err
	}
	defer p.redisLock.ReleaseLock(ctx, lockKey)
}

func (p *VoteProcessor) processVote(ctx context.Context, voteEvent event2.VoteMessage) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteProcessor.processVote")
	defer span.End()

	hasVoted, err := p.voteRepo.
}
