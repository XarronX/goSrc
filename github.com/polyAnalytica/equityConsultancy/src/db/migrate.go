package db

func Migrate(migrations interface{}) error {
	return connection.AutoMigrate(migrations)
}
