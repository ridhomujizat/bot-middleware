package main

import (
	"bot-middleware/config"
	"bot-middleware/internal/application"
	appAccount "bot-middleware/internal/application/account"
	appBot "bot-middleware/internal/application/bot"
	appSession "bot-middleware/internal/application/session"
	"bot-middleware/internal/pkg/libs"
	"bot-middleware/internal/pkg/messaging"
	"bot-middleware/internal/pkg/messaging/rabbit"
	"bot-middleware/internal/pkg/repository/postgre"
	"bot-middleware/internal/pkg/repository/redis"
	"bot-middleware/internal/pkg/util"
	webhookFacebook "bot-middleware/internal/webhook/facebook"
	webhookLivechat "bot-middleware/internal/webhook/livechat"
	webhookWhatsapp "bot-middleware/internal/webhook/whatsapp"
	workerWhatsapp "bot-middleware/internal/worker/whatsapp"
	"log"

	webhookTole "bot-middleware/internal/webhook/tole"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"bot-middleware/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	webHookTelegram "bot-middleware/internal/webhook/telegram"
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

	// Init Redis
	newRedisClient := redis.NewRedisClient(fmt.Sprintf("%s:%s", util.GodotEnv("REDIS_HOST"), util.GodotEnv("REDIS_PORT")), util.GodotEnv("REDIS_PASSWORD"), 0)
	if newRedisClient == nil {
		log.Fatal("Failed to create Redis client")
	}
	defer newRedisClient.Close()

	// Init DB
	applicationService, libsService := initDB()

	// Init RabbitMQ
	rabbitPublisher, rabbitSubscriber := initRabbitMQ(cfg, newRedisClient)

	// Init Messaging
	messagingService := messaging.NewMessagingGeneral(rabbitPublisher, rabbitSubscriber)

	// Init Subscriber
	initSubscriber(messagingService, applicationService, libsService)

	// Init Router
	router := initRouter(messagingService)

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

func initRabbitMQ(cfg config.RabbitMQConfig, redisClient *redis.RedisClient) (*rabbit.RabbitMQPublisher, *rabbit.RabbitMQSubscriber) {
	rabbitPublisher, err := rabbit.NewRabbitMQPublisher(cfg)
	if err != nil {
		util.HandleAppError(err, "initRabbitMQ", "NewRabbitMQPublisher", true)
	}

	rabbitSubscriber, err := rabbit.NewRabbitMQSubscriber(cfg, rabbitPublisher, redisClient)
	if err != nil {
		util.HandleAppError(err, "initRabbitMQ", "NewRabbitMQSubscriber", true)
	}

	return rabbitPublisher, rabbitSubscriber
}

func initRouter(messagingGeneral messaging.MessagingGeneral) *gin.Engine {
	router := gin.Default()

	// Add Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize API routes
	routeGroup := router.Group("/api/v1")

	webhookTole.InitRouterTole(messagingGeneral, routeGroup)
	webHookTelegram.InitRouterTelegram(messagingGeneral, routeGroup)
	webhookFacebook.InitRouterFacebook(messagingGeneral, routeGroup)
	webhookWhatsapp.InitRouterWhatsapp(messagingGeneral, routeGroup)
	webhookLivechat.InitRouterLivechat(messagingGeneral, routeGroup)

	return router
}

func initDB() (*application.Services, *libs.LibsService) {
	db, err := postgre.GetDB()
	if err != nil {
		util.HandleAppError(err, "main", "initDB", true)
	}

	services := &application.Services{
		AccountService: appAccount.NewAccountService(db),
		SessionService: appSession.NewSessionService(db),
		BotService:     appBot.NewBotService(db),
	}

	libsService := libs.NewLibsService(db)

	return services, libsService
}

func initSubscriber(messagingGeneral messaging.MessagingGeneral, applicationService *application.Services, libsService *libs.LibsService) {
	whatappSubscriber := workerWhatsapp.NewWhatsappService(messagingGeneral, applicationService, libsService)
	whatappSubscriber.Subscribe("exchange", "routingKey", util.GodotEnv("QUEUE_WHATSAPP_INITIATE"), false, whatappSubscriber.Process)
	whatappSubscriber.Subscribe("exchange", "routingKey", util.GodotEnv("QUEUE_WHATSAPP_BOT"), false, whatappSubscriber.ProcessBot)
	whatappSubscriber.Subscribe("exchange", "routingKey", util.GodotEnv("QUEUE_WHATSAPP_OUTGOING"), false, whatappSubscriber.ProcessOutgoing)
	whatappSubscriber.Subscribe("exchange", "routingKey", util.GodotEnv("QUEUE_WHATSAPP_FINISH"), false, whatappSubscriber.ProcessFinish)
}
