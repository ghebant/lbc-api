package constants

// Handler path
const (
	HealthPath   = "/"
	SearchPath   = "/search"
	AdPath       = "/ad"
	AdWithIdPath = "/ad/:id"
)

// Mysql
const (
	AdPrimaryKey = "ad_id"
	MysqlDriver  = "mysql"
)

var Vehicles map[string]string
