package templ

const Fiber = `
package routing

import (
	"fmt"
	"{{.ProjectName}}/docs"
	models "{{.ProjectName}}/model"
	"strings"
	"time"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	defaultSkipper = func(c *fiber.Ctx) bool {
		return strings.Contains(c.Path(), "swagger")
	}
)

// FiberSkipper skip next
type fiberSkipper func(c *fiber.Ctx) bool

// FiberMiddleware init
type FiberMiddleware struct {
	Skipper fiberSkipper
	Config  interface{}
}

// InitFiber init http fiber
func InitFiber() *FiberMiddleware {
	m := &FiberMiddleware{
		Skipper: defaultSkipper,
		Config: fiber.New(fiber.Config{
			CaseSensitive: false,
			StrictRouting: false,
			ErrorHandler: func(f *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				if e, isOk := err.(*fiber.Error); isOk {
					code = e.Code
				}
				logrus.WithFields(
					logrus.Fields{
						"code": code,
					}).Error(err)
				return fiberError(code, f)
			},
		}),
	}
	return m
}
// InitFiberMiddleware fiber use
func (m *FiberMiddleware) InitFiberMiddleware(cacheStorage fiber.Storage) (*fiber.App, fiber.Router) {
	service := viper.GetString("app.service")
	basePath := "/api/" + service
	f := m.Config.(*fiber.App)
	f.Use(cors.New())
	f.Use(recover.New())
	f.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			if auth := c.Get(fiber.HeaderAuthorization); len(auth) > 0 {
				return true
			}
			if cacheControl := c.Get(fiber.HeaderCacheControl); len(cacheControl) > 0 {
				return true
			}
			return false
		},
		Expiration:   10 * time.Second,
		CacheControl: true,
		Storage:      cacheStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("%s?%s", c.Path(), c.Request().URI().QueryArgs().String())
		},
	}))
	f.Use(m.fiberLogResponse)
	f.Use(m.fiberLogRequest)
	router := f.Group(basePath)
	if viper.GetString("app.env") != "prod" {
		docs.SwaggerInfo.Title = "Swagger " + service + " service"
		docs.SwaggerInfo.Host = viper.GetString("app.host")
		docs.SwaggerInfo.BasePath = basePath
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		router.Get("/swagger/*", swagger.Handler)
	}
    router.Get("/healthcheck", healthcheck)
	return f, router
}

// Healthcheck handler
// ShowAccount godoc
// @Success 200 {object} model.BaseResponse{data=bool} "server is ok"
// @Failure default {object} model.BaseErrorResponse{code=int,error=string} "server is not ok"
// @Router /healthcheck [get]
func healthcheck(c *fiber.Ctx) error {
	return c.JSON(models.NewBaseResponse(0, true))
}

func fiberError(code int, f *fiber.Ctx) error {
	switch code {
	case 9:
		return f.Status(fiber.StatusNotFound).JSON(models.NewBaseErrorResponse(code, "Not found"))
	case 10:
		return f.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(code, "Invalid parameter"))
	default:
		return f.Status(code).JSON(models.NewBaseErrorResponse(90, "Internal server error"))
	}
}

func (m *FiberMiddleware) fiberLogRequest(c *fiber.Ctx) error {
	if m.Skipper(c) {
		return c.Next()
	}
	body := string(c.Request().Body())
	method := string(c.Request().Header.Method())
	logrus.WithFields(
		logrus.Fields{
			"method": method,
			"path":   c.Path(),
			"body":   string(body),
		}).Infof("Request")
	return c.Next()
}

func (m *FiberMiddleware) fiberLogResponse(c *fiber.Ctx) error {
	if m.Skipper(c) {
		return c.Next()
	}
	if err := c.Next(); err != nil {
		return err
	}
	body := string(c.Response().Body())
	method := string(c.Request().Header.Method())
	logrus.WithFields(
		logrus.Fields{
			"method": method,
			"path":   c.Path(),
			"body":   string(body),
		}).Infof("Response")
	return nil
}
`
