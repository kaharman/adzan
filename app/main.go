package main

import (
	_DeliveryHTTPAdzan "adzan/adzan/delivery/http"
	"adzan/config"

	"flag"
	"net"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
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

	// Set log level
	switch viper.GetString("server.log_level") {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// // Initialize database
	// connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", viper.GetString("mysql.user"), viper.GetString("mysql.password"), viper.GetString("mysql.host"), viper.GetString("mysql.port"), viper.GetString("mysql.database"))
	// val := url.Values{}
	// val.Add("parseTime", "1")
	// val.Add("loc", "Asia/Jakarta")
	// dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	// dbConn, err := sql.Open(`mysql`, dsn)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = dbConn.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer func() {
	// 	err := dbConn.Close()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// // Migrate database if any new schema
	// driver, err := mysql.WithInstance(dbConn, &mysql.Config{})
	// if err == nil {
	// 	mig, err := migrate.NewWithDatabaseInstance(viper.GetString("mysql.path_migrate"), viper.GetString("mysql.database"), driver)
	// 	if err == nil {
	// 		err = mig.Up()
	// 		if err != nil && err != migrate.ErrNoChange {
	// 			version, _, _ := mig.Version()
	// 			ver := int(version) - 1
	// 			mig.Force(ver)
	// 		} else {
	// 			log.Info("Migrate database success")
	// 		}
	// 	} else {
	// 		log.Warn(err)
	// 	}
	// } else {
	// 	log.Warn(err)
	// }

	// // Initialize Redis
	// dbRedis := redis.NewClient(&redis.Options{
	// 	Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
	// 	Username: viper.GetString("redis.username"),
	// 	Password: viper.GetString("redis.password"),
	// 	DB:       viper.GetInt("redis.database"),
	// 	PoolSize: viper.GetInt("redis.max_connection"),
	// })

	// _, err = dbRedis.Ping(context.Background()).Result()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Info("Redis connection established")

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
			log.Fatalf("[ERROR] Failed to listen tcp: %v", err)
		}

		grpcServer := grpc.NewServer()
		// TODO: reegister services

		log.Println("gRPC server is running")
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Initialize HTTP web framework
	app := fiber.New(fiber.Config{
		Prefork:       viper.GetBool("server.prefork"),
		StrictRouting: viper.GetBool("server.strict_routing"),
		CaseSensitive: viper.GetBool("server.case_sensitive"),
		BodyLimit:     viper.GetInt("server.body_limit"),
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))

	// HTTP routing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/web")
	})
	// TODO: routing here
	_DeliveryHTTPAdzan.RouterWeb(app)
	_DeliveryHTTPAdzan.RouterOpenAPI(app)
	_DeliveryHTTPAdzan.RouterGRPC(app)

	if err := app.Listen(":" + viper.GetString("server.port")); err != nil {
		log.Fatal(err)
	}
}
