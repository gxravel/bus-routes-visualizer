package mysql

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	"github.com/gxravel/bus-routes-visualizer/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// RoutePointStore is route_point mysql store.
type RoutePointStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

// NewRoutePointStore creates new instance of RoutePointStore.
func NewRoutePointStore(db sqlx.ExtContext, txer dataprovider.Txer) *RoutePointStore {
	return &RoutePointStore{
		db:        db,
		txer:      txer,
		tableName: "route_point",
	}
}

// WithTx sets transaction as active connection.
func (s *RoutePointStore) WithTx(tx *dataprovider.Tx) dataprovider.RoutePointStore {
	return &RoutePointStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func (s *RoutePointStore) columns() []string {
	return []string{
		"route_id",
		"step",
		"address",
	}
}

// Add creates new routes.
func (s *RoutePointStore) Add(ctx context.Context, points ...*model.RoutePoint) error {
	qb := sq.Insert(s.tableName).Columns(s.columns()...)

	for _, point := range points {
		qb = qb.Values(point.RouteID, point.Step, point.Address)
	}

	query, args, codewords, err := toSql(ctx, qb, s.tableName)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, codewords+" with query %s", query)
	}
	return nil
}

// Update updates route_point's stop_id.
func (s *RoutePointStore) Update(ctx context.Context, point *model.RoutePoint) error {
	qb := sq.Update(s.tableName).Set("step", point.Step).Set("address", point.Address).Where(sq.Eq{"route_id": point.RouteID})

	return execContext(ctx, qb, s.tableName, s.txer)
}

// Delete deletes route_point depend on received filter.
func (s *RoutePointStore) Delete(ctx context.Context, id int64) error {
	qb := sq.Delete(s.tableName).Where(sq.Eq{"id": id})

	return execContext(ctx, qb, s.tableName, s.txer)
}
