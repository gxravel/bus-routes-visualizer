module github.com/gxravel/bus-routes-visualizer

go 1.16

require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fogleman/gg v1.3.0
	github.com/go-chi/chi v1.5.4
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-swagger/go-swagger v0.27.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golangci/golangci-lint v1.41.1
	github.com/gxravel/bus-routes/pkg/rmq v0.0.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lopezator/migrator v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/spf13/viper v1.8.1
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
)

replace github.com/gxravel/bus-routes/pkg/rmq v0.0.0 => ../bus-routes/pkg/rmq
