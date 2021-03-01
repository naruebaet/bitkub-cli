package templ

const Mysql = `
package drivers

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func myConnection(dns string, source string, replica string) *gorm.DB {
	conn, err := gorm.Open(mysql.New(
		mysql.Config{
			DSN: dns,
		}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	conn.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(source)},
		Replicas: []gorm.Dialector{mysql.Open(replica)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		logrus.Fatalf("cannot open mysql connection: %s", err)
	}
	logrus.Info("connect mysql success")
	return conn
}

//MySQLWeb connect mysql dbms
func MySQLWeb() *gorm.DB {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("mysql_web.username"),
		viper.GetString("mysql_web.password"),
		viper.GetString("mysql_web.host"),
		viper.GetString("mysql_web.port"),
		viper.GetString("mysql_web.database"),
	)
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("mysql_web.username"),
		viper.GetString("mysql_web.password"),
		viper.GetString("mysql_web.host"),
		viper.GetString("mysql_web.port"),
		viper.GetString("mysql_web.database"),
	)
	replica := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("mysql_web_read.username"),
		viper.GetString("mysql_web_read.password"),
		viper.GetString("mysql_web_read.host"),
		viper.GetString("mysql_web_read.port"),
		viper.GetString("mysql_web_read.database"),
	)
	return myConnection(DSN, source, replica)
}

`
