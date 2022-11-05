package api

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"github.com/martikan/carrental_auth-api/config"
	dbConn "github.com/martikan/carrental_auth-api/db/sqlc"
	"github.com/martikan/carrental_common/middleware"
	"github.com/martikan/carrental_common/util"
)

type Api struct {
	router     *gin.Engine
	db         *dbConn.Queries
	sqlDb      *sql.DB
	config     *config.Config
	tokenMaker util.Maker
}

func InitApi() *Api {

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v\n", err)
	}

	tokenMaker, err := util.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}

	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.PostgreUser, config.PostgrePassword,
		config.PostgreHost, config.PostgrePort, config.PostgreDb, config.SSLMode)
	log.Println(dbUrl)

	conn, err := sql.Open(config.PostgreDriver, dbUrl)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("could not reach the database: %v", err)
	}

	log.Println("successfully connected to database")

	api := &Api{
		db:         dbConn.New(conn),
		sqlDb:      conn,
		config:     &config,
		tokenMaker: tokenMaker,
	}
	api.setupRouter()

	return api
}

func (a *Api) setupRouter() {
	router := gin.Default()

	// Open routes

	// FIXME: Implement it for k8s
	// Health
	// router.GET("/health/live", a.live)
	// router.GET("/health/ready", a.ready)

	router.POST("/api/v1/auth/signin", a.signIn)
	router.POST("/api/v1/auth/signup", a.signUp)

	// Authenticated routes

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(a.tokenMaker))

	authRoutes.GET("/api/v1/auth/current_user", a.currentUser)

	a.router = router
}

func (a *Api) Start() {
	log.Fatal(a.router.Run(":" + a.config.Port))
}
