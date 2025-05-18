package dao

import (
	"errors"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/vote/internal/domain/model"
	"github.com/nocturna-ta/vote/internal/domain/repository"
)

var (
	ErrOptimisticLockFailed = errors.New("optimistic lock failed: the record was updated by another transaction")
)

type ShardedVoteRepository struct {
	shards     []*VoteRepository
	shardCount int
	client     ethereum.Client
}

func NewShardedVoteRepository(opts *OptsVoteRepository, shardCount int, shardDBs []*sql.Store) repository.VoteRepository {
	if shardCount <= 0 {
		shardCount = 1
	}

	if len(shardDBs) == 0 || len(shardDBs) != shardCount {
		log.Fatal("[ShardedVoteRepository] Invalid shard configuration")
	}

	shards := make([]*VoteRepository, shardCount)

	for i := 0; i < shardCount; i++ {
		shardOpts := &OptsVoteRepository{
			Client:          opts.Client,
			ContractAddress: opts.ContractAddress,
			Contract:        opts.Contract,
			DB:              shardDBs[i],
		}

		shard, err := NewVoteRepository(shardOpts)
		if err != nil {
			log.Fatal("[ShardedVoteRepository] Failed to create shard repository", err)
		}

		shards[i] = shard.(*VoteRepository)
	}

	return &ShardedVoteRepository{
		shards:     shards,
		shardCount: shardCount,
		client:     opts.Client,
	}
}

func (s *ShardedVoteRepository) InsertVote(ctx context.Context, vote *model.Vote) error {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) GetVoteByID(ctx context.Context, id uuid.UUID) (*model.Vote, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) GetVoteByElectionPairID(ctx context.Context, electionPairID uuid.UUID) (*model.Vote, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) UpdateVoteStatus(ctx context.Context, id uuid.UUID, status model.VoteStatus) error {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) UpdateVote(ctx context.Context, vote *model.Vote) error {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) GetVotesByVoterID(ctx context.Context, voterID uuid.UUID) ([]*model.Vote, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) GetLatestVoteByVoterID(ctx context.Context, voterID uuid.UUID) (*model.Vote, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) HasVoted(ctx context.Context, voterID uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) GetVoterByKTP(ctx context.Context, ktp string) (*model.Voter, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardedVoteRepository) GetVoterByID(ctx context.Context, id uuid.UUID) (*model.Voter, error) {
	//TODO implement me
	panic("implement me")
}
