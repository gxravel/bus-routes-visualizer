package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

type RouteStore interface {
	WithTx(*Tx) RouteStore
	GetByFilter(ctx context.Context, filter *RouteFilter) (*model.RouteJoined, error)
	GetListByFilter(ctx context.Context, filter *RouteFilter) ([]*model.RouteJoined, error)
	Add(ctx context.Context, routes ...*model.Route) error
	Update(ctx context.Context, route *model.Route) error
	Delete(ctx context.Context, filter *RouteFilter) error
}

type RouteFilter struct {
	Buses     []string
	Cities    []string
	Addresses []string
}

func NewRouteFilter() *RouteFilter {
	return &RouteFilter{}
}

// ByBuses filters by route.bus.
func (f *RouteFilter) ByBuses(buses ...string) *RouteFilter {
	f.Buses = buses
	return f
}

// ByCities filters by route.city.
func (f *RouteFilter) ByCities(cities ...string) *RouteFilter {
	f.Cities = cities
	return f
}

// ByAddress filters by route.address.
func (f *RouteFilter) ByAddresses(addresses ...string) *RouteFilter {
	f.Addresses = addresses
	return f
}
