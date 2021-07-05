package model

// Route describes route joined with route_point.
type RouteJoined struct {
	ID   int64  `db:"id"`
	Bus  string `db:"bus"`
	City string `db:"city"`

	Step    int8   `db:"step"`
	Address string `db:"address"`
}

// Route describes route in bus_routes_visualizor.route.
type Route struct {
	ID   int64  `db:"id"`
	Bus  string `db:"bus"`
	City string `db:"city"`

	// implicitly
	Points []*RoutePoint
}

// Route describes route point in bus_routes_visualizor.route_point.
type RoutePoint struct {
	Step    int8   `db:"step"`
	Address string `db:"address"`
	RouteID int64  `db:"route_id"`
}
