package model

// Permission describes user permissions in bus_routes_visualizor.permission.
type Permission struct {
	UserID  int64 `db:"user_id"`
	Actions JSON  `db:"actions"`
}
