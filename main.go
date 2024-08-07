package main

import (
	"bot-middleware/config"
	"bot-middleware/internal/application"
	appAccount "bot-middleware/internal/application/account"
	appBot "bot-middleware/internal/application/bot"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/messaging/rabbit"
	"bot-middleware/internal/pkg/repository/postgre"
	"bot-middleware/internal/pkg/util"
	workerTelegram "bot-middleware/internal/worker/telegram"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"bot-middleware/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	webHookTelegram "bot-middleware/internal/webhook/telegram"
	webhookTole "bot-middleware/internal/webhook/tole"
)

// @title Bot Middleware API
// @version 1.0
// @description API documentation for the Bot Middleware service.
// @BasePath /api/v1
func main() {
	// Set Swagger annotations dynamically
	setSwaggerInfo()

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		util.HandleAppError(err, "main", "loadConfig", true)
	}

	// Init DB
	applicationService := initDB()

	// Init RabbitMQ
	rabbitPublisher, rabbitSubscriber := initRabbitMQ(cfg)

	// Init Messaging
	messagingService := messaging.NewMessagingGeneral(rabbitPublisher, rabbitSubscriber)

	// Init Subscriber
	initSubscriber(messagingService, applicationService)

	// Init Router
	router := initRouter(messagingService, applicationService)

	port := util.GodotEnv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(fmt.Sprintf(":%s", port))
}

func setSwaggerInfo() {
	host := util.GodotEnv("HOST")
	if host != "" {
		docs.SwaggerInfo.Host = host
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	}
}

func loadConfig() (config.RabbitMQConfig, error) {
	cfg := config.LoadRabbitMQConfig()
	if cfg.URL == "" {
		return cfg, errors.New("RabbitMQ URL not provided")
	}
	return cfg, nil
}

func initRabbitMQ(cfg config.RabbitMQConfig) (*rabbit.RabbitMQPublisher, *rabbit.RabbitMQSubscriber) {
	rabbitPublisher, err := rabbit.NewRabbitMQPublisher(cfg)
	if err != nil {
		util.HandleAppError(err, "initRabbitMQ", "NewRabbitMQPublisher", true)

	}

	rabbitSubscriber, err := rabbit.NewRabbitMQSubscriber(cfg, false)
	if err != nil {
		util.HandleAppError(err, "initRabbitMQ", "NewRabbitMQSubscriber", true)
	}

	return rabbitPublisher, rabbitSubscriber
}

func initRouter(messagingGeneral messaging.MessagingGeneral, applicationService *application.Services) *gin.Engine {
	router := gin.Default()

	// Add Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize API routes
	routeGroup := router.Group("/api/v1")
	webhookTole.InitRouterTole(messagingGeneral, routeGroup)
	webHookTelegram.InitRouterTelegram(messagingGeneral, routeGroup, applicationService)

	return router
}

func initDB() *application.Services {
	db, err := postgre.GetDB()
	if err != nil {
		util.HandleAppError(err, "main", "initDB", true)
	}

	services := &application.Services{
		AccountService:  appAccount.NewAccountService(db),
		SessinonService: appSession.NewSessionService(db),
		BotService:      appBot.NewBotService(db),
	}
	return services

}

func initSubscriber(messagingGeneral messaging.MessagingGeneral, applicationService *application.Services) {
	workerTelegram.NewTelegramIncoming(messagingGeneral, applicationService, "exchange", "incoming", "onx:onx_dev:telegram:@BaruBelajarGolangBot", false)
	workerTelegram.NewTelegramBotProcess(messagingGeneral, applicationService, "exchange", "bot-process", "onx:onx_dev:telegram:@BaruBelajarGolangBot:bot", false)
}
