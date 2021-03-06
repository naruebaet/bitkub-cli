package templ

const Postgres = `
package drivers

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func connection(dns string, source string, replica string) *gorm.DB {
	conn, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN: dns,
		}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	conn.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{postgres.Open(source)},
		Replicas: []gorm.Dialector{postgres.Open(replica)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		logrus.Fatalf("cannot open postgres connection:%s", err)
	}
	return conn
}

//PostgresConnection connect postgres db
func PostgresConnection() *gorm.DB {
	DSN := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.dbname"))

	source := fmt.Sprintf("user=%s password=%s host=%s",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"))

	replica := fmt.Sprintf("user=%s password=%s host=%s",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host-repli-1"))

	return connection(DSN, source, replica)
}

`
