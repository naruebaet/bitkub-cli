package templ

const Mongodb = `
package drivers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConnection driver
func MongoConnection() (*mongo.Client, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uri string
	var err error
	var client *mongo.Client

	if viper.GetString("app.env") == "localhost" {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+
			viper.GetString("mongo.host")+":"+viper.GetString("mongo.port"),
		).
			SetAuth(options.Credential{
				AuthSource: viper.GetString("mongo.authSource"),
				Username:   viper.GetString("mongo.username"),
				Password:   viper.GetString("mongo.password"),
			}))
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s/?replicaSet=%s&authSource=%s&mechanism=%s&tls=true",
			viper.GetString("mongo.username"),
			viper.GetString("mongo.password"),
			viper.GetString("mongo.host"),
			viper.GetString("mongo.replicaSet"),
			viper.GetString("mongo.authSource"),
			viper.GetString("mongo.mechanism"),
		)
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	}
	logrus.Info(uri)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	logrus.Info("Connect mongodb success")

	return client, ctx, err
}

`
