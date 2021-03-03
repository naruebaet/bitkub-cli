package templ

const Main = `
package main
import (
	"os"
	"os/signal"
	"{{.ProjectName}}/pkg/routing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)
func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	initViper()
}

func main() {
	var	newFiber  = routing.InitFiber()
	f, router := newFiber.InitFiberMiddleware(nil)
	_ = router.Group("/v1")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		logrus.Info("Gracefully shutting down...")
		_ = f.Shutdown()
	}()
	if err := f.Listen(":" + viper.GetString("app.port")); err != nil {
		logrus.Fatalf("shutting down the server : %s", err)
	}
}

func initViper() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("cannot read in viper config:%s", err)
	}
	viper.AutomaticEnv()
}
`
