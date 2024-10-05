package database

func NewMigration(models ...interface{}) error {

	err := DB.AutoMigrate(models...)

	return err
}
