package mysql

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	"github.com/gxravel/bus-routes-visualizer/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// PermissionStore is permission mysql store.
type PermissionStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

// NewPermissionStore creates new instance of PermissionStore.
func NewPermissionStore(db sqlx.ExtContext, txer dataprovider.Txer) *PermissionStore {
	return &PermissionStore{
		db:        db,
		txer:      txer,
		tableName: "permission",
	}
}

// WithTx sets transaction as active connection.
func (s *PermissionStore) WithTx(tx *dataprovider.Tx) dataprovider.PermissionStore {
	return &PermissionStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func permissionCond(f *dataprovider.PermissionFilter) sq.Sqlizer {
	and := make(sq.And, 0)

	if len(f.UserIDs) > 0 {
		eq := make(sq.Eq)
		eq["permission.user_id"] = f.UserIDs

		and = append(and, eq)
	}

	if len(f.Actions) > 0 {
		or := make(sq.Or, 0, len(f.Actions))
		for _, action := range f.Actions {
			or = append(or, sq.Expr("JSON_CONTAINS(actions, JSON_QUOTE(?), '$.actions') = 1", action))
		}

		and = append(and, or)
	}

	return and
}

func (s *PermissionStore) columns(filter *dataprovider.PermissionFilter) []string {
	return []string{
		"user_id",
		"actions",
	}
}

// GetByFilter returns permission depend on received filters.
func (s *PermissionStore) GetByFilter(ctx context.Context, filter *dataprovider.PermissionFilter) (*model.Permission, error) {
	permissions, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(permissions) == 0:
		return nil, nil
	case len(permissions) == 1:
		return permissions[0], nil
	default:
		return nil, errors.New("fetched more than 1 permission")
	}
}

// GetListByFilter returns permissions depend on received filters.
func (s *PermissionStore) GetListByFilter(ctx context.Context, filter *dataprovider.PermissionFilter) ([]*model.Permission, error) {
	qb := sq.
		Select(s.columns(filter)...).
		From(s.tableName).
		Where(permissionCond(filter))

	query, args, codewords, err := toSql(ctx, qb, s.tableName)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Permission, 0)

	if err := sqlx.SelectContext(ctx, s.db, &result, query, args...); err != nil {
		return nil, errors.Wrapf(err, "%s by filter with query %s", codewords, query)
	}

	return result, nil
}

// Add creates new permissions.
func (s *PermissionStore) Add(ctx context.Context, permissions ...*model.Permission) error {
	qb := sq.Insert(s.tableName).Columns(s.columns(nil)...)

	for _, p := range permissions {
		actions, err := p.Actions.Value()
		if err != nil {
			return errors.Wrapf(err, "couldn't get driver.Value of actions: %v", actions)
		}

		qb = qb.Values(p.UserID, actions)
	}

	return execContext(ctx, qb, s.tableName, s.db)
}

// Update updates permissions' actions.
// It skips the missing ids and returns affected ones.
func (s *PermissionStore) Update(ctx context.Context, permissions ...*model.Permission) ([]int64, error) {
	qb := sq.Update(s.tableName)

	ids := make([]int64, 0)

	f := func(tx *dataprovider.Tx) error {
		for _, p := range permissions {
			actions, err := p.Actions.Value()
			if err != nil {
				return errors.Wrapf(err, "couldn't get driver.Value of actions: %v", actions)
			}

			ub := qb.Set("actions", actions).
				Where(sq.Eq{"user_id": p.UserID})

			if err := execContext(ctx, ub, s.tableName, tx); err != nil {
				if errors.Is(err, errNoRowsAffected) {
					continue
				}
				return err
			}

			ids = append(ids, p.UserID)
		}

		return nil
	}

	return ids, dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
}

// Delete deletes permission depend on received filter.
func (s *PermissionStore) Delete(ctx context.Context, filter *dataprovider.PermissionFilter) error {
	qb := sq.Delete(s.tableName).Where(permissionCond(filter))

	return execContext(ctx, qb, s.tableName, s.db)
}
