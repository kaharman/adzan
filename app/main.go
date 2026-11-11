package main

import (
	_DeliveryHTTPAdzan "adzan/adzan/delivery/http"
	_RepositoryCSVAdzan "adzan/adzan/repository/csv"
	_RepositoryLibAdzan "adzan/adzan/repository/lib"
	_RepositoryRedisAdzan "adzan/adzan/repository/redis"
	_UsecaseAdzan "adzan/adzan/usecase"
	"adzan/config"
	"adzan/domain"
	"adzan/helper"
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// CLI options parse
	configFile := flag.String("c", "config.yaml", "Config file")
	flag.Parse()

	// Config file
	config.ReadConfig(*configFile)

	// Set log level with standard library slog
	var logLevel slog.Level
	switch strings.ToLower(viper.GetString("server.log_level")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	logType := strings.ToLower(viper.GetString("server.log_type"))
	var logHandler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		Level: logLevel,
	}
	if logType == "text" {
		logHandler = slog.NewTextHandler(os.Stdout, handlerOpts)
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	}
	slog.SetDefault(slog.New(logHandler))
	slog.Info("Initialized slog logger", "level", logLevel.String(), "type", viper.GetString("server.log_type"))

	// Initialize Data Source (Redis or CSV)
	var prayTimeRepo domain.PrayTimeRepository
	var cityRepo domain.CityRepository

	dataSource := viper.GetString("data_source")
	slog.Info("Configured data source", "data_source", dataSource)
	slog.Debug("Resolved server configuration",
		"port", viper.GetString("server.port"),
		"grpc_port", viper.GetString("server.grpc_port"),
		"base_path", viper.GetString("server.base_path"),
		"prefork", viper.GetBool("server.prefork"),
	)

	if dataSource == "csv" {
		csvPath := viper.GetString("csv.path")
		if csvPath == "" {
			csvPath = "./data"
		}
		slog.Info("Initializing CSV Repository", "path", csvPath)
		slog.Debug("Loading prayer times and cities from CSV directory", "csv_path", csvPath)
		prayTimeRepo = _RepositoryCSVAdzan.NewCSVPrayTimeRepository(csvPath)
		cityRepo = _RepositoryCSVAdzan.NewCSVCityRepository(csvPath)
	} else {
		// Initialize Redis
		redisAddr := viper.GetString("redis.host") + ":" + viper.GetString("redis.port")
		slog.Info("Initializing Redis connection", "addr", redisAddr, "db", viper.GetInt("redis.database"))
		slog.Debug("Redis client pool options",
			"pool_size", viper.GetInt("redis.max_connection"),
			"username", viper.GetString("redis.username"),
		)
		dbRedis := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Username: viper.GetString("redis.username"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.database"),
			PoolSize: viper.GetInt("redis.max_connection"),
		})

		_, err := dbRedis.Ping(context.Background()).Result()
		if err != nil {
			slog.Error("Failed to connect to Redis", "error", err)
			os.Exit(1)
		}
		slog.Info("Redis connection established successfully")

		prayTimeRepo = _RepositoryRedisAdzan.NewRedisPrayTimeRepository(dbRedis)
		cityRepo = _RepositoryRedisAdzan.NewRedisCityRepository(dbRedis)
	}

	// // Initialize gRPC pool
	// conn := func() (client *grpc.ClientConn, err error) {
	// 	address := fmt.Sprintf("%s:%s", viper.GetString("grpc.host"), viper.GetString("grpc.port"))
	// 	client, err = grpc.Dial(address, grpc.WithInsecure())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	return
	// }

	// domain.GRPCPool, err = grpcpool.New(conn, viper.GetInt("grpc.init"), viper.GetInt("grpc.capacity"), time.Duration(viper.GetInt("grpc.idle_duration"))*time.Second, time.Duration(viper.GetInt("grpc.max_life_duration"))*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Info("gRPC pool created")

	// Register repository & usecase
	// TODO:

	// Initialize gRPC server
	go func() {
		listen, err := net.Listen("tcp", ":"+viper.GetString("server.grpc_port"))
		if err != nil {
			slog.Error("Failed to listen tcp for gRPC", "error", err)
			os.Exit(1)
		}

		grpcServer := grpc.NewServer()
		// TODO: register services

		slog.Info("gRPC server is running", "port", viper.GetString("server.grpc_port"))
		if err := grpcServer.Serve(listen); err != nil {
			slog.Error("Failed to serve gRPC", "error", err)
			os.Exit(1)
		}
	}()

	// Initialize HTTP web framework
	app := fiber.New(fiber.Config{
		Prefork:       viper.GetBool("server.prefork"),
		StrictRouting: viper.GetBool("server.strict_routing"),
		CaseSensitive: viper.GetBool("server.case_sensitive"),
		BodyLimit:     viper.GetInt("server.body_limit"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return helper.HttpError(c, code, err.Error(), nil)
		},
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))

	// Log incoming request and received headers at slog.Info level
	app.Use(func(c *fiber.Ctx) error {
		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(key, val []byte) {
			headers[string(key)] = string(val)
		})
		slog.Info("Incoming request received",
			"method", c.Method(),
			"path", c.OriginalURL(),
			"remote_ip", c.IP(),
			"headers", headers,
		)
		return c.Next()
	})

	// HTTP routing with optional base_path prefix
	basePath := strings.TrimSpace(viper.GetString("server.base_path"))
	if basePath != "" && !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	basePath = strings.TrimRight(basePath, "/")

	var router fiber.Router = app
	if basePath != "" {
		slog.Info("Using server base_path routing prefix", "base_path", basePath)
		router = app.Group(basePath)
		app.Get("/", func(c *fiber.Ctx) error {
			return c.Redirect(basePath + "/web")
		})
		router.Get("/", func(c *fiber.Ctx) error {
			return c.Redirect(basePath + "/web")
		})
	} else {
		app.Get("/", func(c *fiber.Ctx) error {
			return c.Redirect("/web")
		})
	}

	_DeliveryHTTPAdzan.RouterWeb(router)
	_DeliveryHTTPAdzan.RouterOpenAPI(router)
	_DeliveryHTTPAdzan.RouterGRPC(router)
	hijriRepo := _RepositoryLibAdzan.NewLibHijriRepository()
	prayTimeUsecase := _UsecaseAdzan.NewPrayTimeUsecase(cityRepo, prayTimeRepo)
	cityUsecase := _UsecaseAdzan.NewCityUsecase(cityRepo)
	hijriUsecase := _UsecaseAdzan.NewHijriUsecase(hijriRepo)
	_DeliveryHTTPAdzan.RouterAPI(router, prayTimeUsecase, cityUsecase, hijriUsecase)

	port := viper.GetString("server.port")
	slog.Info("Starting HTTP server", "port", port)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		os.Exit(1)
	}
}
