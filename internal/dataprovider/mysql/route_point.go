package mysql

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	"github.com/gxravel/bus-routes-visualizer/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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

func routePointCond(f *dataprovider.RoutePointFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	var cond sq.Sqlizer = eq

	if len(f.IDs) > 0 {
		eq["route_point.id"] = f.IDs
	}

	return cond
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

	return execContext(ctx, qb, s.tableName, s.db)
}

// Update updates route_point's step and address.
func (s *RoutePointStore) Update(ctx context.Context, point *model.RoutePoint) error {
	qb := sq.Update(s.tableName).
		SetMap(map[string]interface{}{
			"step":    point.Step,
			"address": point.Address,
		}).
		Where(sq.Eq{"route_id": point.RouteID})

	return execContext(ctx, qb, s.tableName, s.db)
}

// Delete deletes route_point depend on received filter.
func (s *RoutePointStore) Delete(ctx context.Context, filter *dataprovider.RoutePointFilter) error {
	qb := sq.Delete(s.tableName).Where(routePointCond(filter))

	return execContext(ctx, qb, s.tableName, s.db)
}
