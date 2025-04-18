package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.robinthrift.com/conveyor/internal/domain"
	"go.robinthrift.com/conveyor/internal/storage/database"
	"go.robinthrift.com/conveyor/internal/storage/database/sqlite/sqlc"
	"go.robinthrift.com/conveyor/internal/storage/database/sqlite/types"
)

type AccountRepo struct {
	db database.Database
}

func NewAccountRepo(db database.Database) *AccountRepo {
	return &AccountRepo{db}
}

func (r *AccountRepo) CountAccounts(ctx context.Context) (int64, error) {
	count, err := queries.CountAccounts(ctx, r.db.Conn(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("error counting accounts: %w", err)
	}

	return count, nil
}

func (r *AccountRepo) Get(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	account, err := queries.GetAccount(ctx, r.db.Conn(ctx), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: id %d", domain.ErrAccountNotFound, id)
		}

		return nil, fmt.Errorf("error getting account by id: id %d: %w", id, err)
	}

	return &domain.Account{
		ID:       account.ID,
		Username: account.Username,
		Password: domain.AccountPassword{
			Algorithm:      account.Algorithm,
			Params:         account.Params,
			Salt:           account.Salt,
			Password:       account.Password,
			RequiresChange: account.RequiresPasswordChange,
		},
		CreatedAt: account.CreatedAt.Time,
		UpdatedAt: account.UpdatedAt.Time,
	}, nil
}

func (r *AccountRepo) GetByUsername(ctx context.Context, username string) (*domain.Account, error) {
	account, err := queries.GetAccountByUsername(ctx, r.db.Conn(ctx), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", domain.ErrAccountNotFound, username)
		}

		return nil, fmt.Errorf("error getting account by username: %s: %w", username, err)
	}

	return &domain.Account{
		ID:       account.ID,
		Username: account.Username,
		Password: domain.AccountPassword{
			Algorithm:      account.Algorithm,
			Params:         account.Params,
			Salt:           account.Salt,
			Password:       account.Password,
			RequiresChange: account.RequiresPasswordChange,
		},
		CreatedAt: account.CreatedAt.Time,
		UpdatedAt: account.UpdatedAt.Time,
	}, nil
}

func (r *AccountRepo) Create(ctx context.Context, toCreate *domain.Account) error {
	return queries.CreateAccount(ctx, r.db.Conn(ctx), sqlc.CreateAccountParams{
		Username:  toCreate.Username,
		Algorithm: toCreate.Password.Algorithm,
		Params:    toCreate.Password.Params,
		Salt:      toCreate.Password.Salt,
		Password:  toCreate.Password.Password,
	})
}

func (r *AccountRepo) Update(ctx context.Context, toUpdate *domain.Account) error {
	return queries.UpdateAccount(ctx, r.db.Conn(ctx), sqlc.UpdateAccountParams{
		ID:                     toUpdate.ID,
		Algorithm:              toUpdate.Password.Algorithm,
		Params:                 toUpdate.Password.Params,
		Salt:                   toUpdate.Password.Salt,
		Password:               toUpdate.Password.Password,
		RequiresPasswordChange: toUpdate.Password.RequiresChange,
		UpdatedAt:              types.NewSQLiteDatetime(time.Now().UTC()),
	})
}

func (r *AccountRepo) CreateAccountKey(ctx context.Context, toCreate *domain.AccountKey) error {
	return queries.CreateAccountKey(ctx, r.db.Conn(ctx), sqlc.CreateAccountKeyParams{
		AccountID: int64(toCreate.AccountID),
		Name:      toCreate.Name,
		Type:      toCreate.Type,
		Data:      toCreate.Data,
	})
}

func (r *AccountRepo) GetAccountKeyByName(ctx context.Context, accountID domain.AccountID, name string) (*domain.AccountKey, error) {
	row, err := queries.GetAccountKeyByName(ctx, r.db.Conn(ctx), sqlc.GetAccountKeyByNameParams{
		AccountID: int64(accountID),
		Name:      name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", domain.ErrAccountKeyNotFound, name)
		}

		return nil, err
	}

	return &domain.AccountKey{
		ID:        domain.AccountKeyID(row.ID),
		AccountID: domain.AccountID(row.AccountID),
		Name:      row.Name,
		Type:      row.Type,
		Data:      row.Data,
	}, nil
}
