package mysql

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	"github.com/gxravel/bus-routes-visualizer/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// RouteStore is route mysql store.
type RouteStore struct {
	db        sqlx.ExtContext
	txer      dataprovider.Txer
	tableName string
}

// NewRouteStore creates new instance of RouteStore.
func NewRouteStore(db sqlx.ExtContext, txer dataprovider.Txer) *RouteStore {
	return &RouteStore{
		db:        db,
		txer:      txer,
		tableName: "route",
	}
}

// WithTx sets transaction as active connection.
func (s *RouteStore) WithTx(tx *dataprovider.Tx) dataprovider.RouteStore {
	return &RouteStore{
		db:        tx,
		tableName: s.tableName,
	}
}

func routeCond(f *dataprovider.RouteFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	var cond sq.Sqlizer = eq

	if len(f.Buses) > 0 {
		eq["route.bus"] = f.Buses
	}
	if len(f.Cities) > 0 {
		eq["route.city"] = f.Cities
	}
	if len(f.Addresses) > 0 {
		eq["route_point.address"] = f.Addresses
	}

	return cond
}

func (s *RouteStore) columns(filter *dataprovider.RouteFilter) []string {
	if filter == nil {
		return []string{
			"bus",
			"city",
		}
	}
	return []string{
		"bus",
		"city",
		"step",
		"address",
	}
}

func (s *RouteStore) joins(qb sq.SelectBuilder, filter *dataprovider.RouteFilter) sq.SelectBuilder {
	qb = qb.Join("route_point ON route.id = route_point.id")
	return qb
}

func (s *RouteStore) ordersBy(qb sq.SelectBuilder, filter *dataprovider.RouteFilter) sq.SelectBuilder {
	qb = qb.OrderBy("city", "bus", "step")
	return qb
}

// GetByFilter returns route depend on received filters.
func (s *RouteStore) GetByFilter(ctx context.Context, filter *dataprovider.RouteFilter) (*model.RouteJoined, error) {
	routes, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(routes) == 0:
		return nil, nil
	case len(routes) == 1:
		return routes[0], nil
	default:
		return nil, errors.New("fetched more than 1 route")
	}
}

// GetListByFilter returns routes depend on received filters.
func (s *RouteStore) GetListByFilter(ctx context.Context, filter *dataprovider.RouteFilter) ([]*model.RouteJoined, error) {
	qb := sq.
		Select(s.columns(filter)...).
		From(s.tableName).
		Where(routeCond(filter))

	qb = s.joins(qb, filter)
	qb = s.ordersBy(qb, filter)

	return selectContext(ctx, qb, s.tableName, s.db)
}

// Add creates new routes.
func (s *RouteStore) Add(ctx context.Context, routes ...*model.Route) error {
	qb := sq.Insert(s.tableName).Columns(s.columns(nil)...)

	for _, route := range routes {
		values := qb.Values(route.Bus, route.City)

		query, args, codewords, err := toSql(ctx, values, "route")
		if err != nil {
			return err
		}

		f := func(tx *dataprovider.Tx) error {
			result, err := tx.ExecContext(ctx, query, args...)
			if err != nil {
				return errors.Wrapf(err, codewords+" with query %s", query)
			}

			lastID, err := result.LastInsertId()
			if err != nil {
				return errors.Wrap(err, "failed to call LastInsertId")
			}

			for i := range route.Points {
				route.Points[i].RouteID = lastID
			}

			pointStore := NewRoutePointStore(s.db, s.txer)

			err = pointStore.WithTx(tx).Add(ctx, route.Points...)
			if err != nil {
				return err
			}

			return nil
		}

		err = dataprovider.BeginAutoCommitedTx(ctx, s.txer, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update updates route's stop_id.
func (s *RouteStore) Update(ctx context.Context, route *model.Route) error {
	qb := sq.Update(s.tableName).Set("bus", route.Bus).Set("city", route.City).Where(sq.Eq{"id": route.ID})

	return execContext(ctx, qb, s.tableName, s.txer)
}

// Delete deletes route depend on received filter.
func (s *RouteStore) Delete(ctx context.Context, filter *dataprovider.RouteFilter) error {
	qb := sq.Delete(s.tableName).Where(routeCond(filter))

	return execContext(ctx, qb, s.tableName, s.txer)
}
