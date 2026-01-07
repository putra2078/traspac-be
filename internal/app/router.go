package app

import (
	"hrm-app/config"
	"hrm-app/internal/domain/auth"
	"hrm-app/internal/domain/boards"
	"hrm-app/internal/domain/boardsUsers"
	"hrm-app/internal/domain/contact"
	"hrm-app/internal/domain/labels"
	room_chats "hrm-app/internal/domain/roomChats"
	room_messages "hrm-app/internal/domain/roomMessages"
	"hrm-app/internal/domain/roomUsers"
	"hrm-app/internal/domain/storage"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskCardComment"
	"hrm-app/internal/domain/taskCardUsers"
	"hrm-app/internal/domain/taskTab"
	"hrm-app/internal/domain/user"
	"hrm-app/internal/domain/workspaces"
	"hrm-app/internal/domain/workspacesUsers"
	"hrm-app/internal/infrastructure/storage/supabase"
	"hrm-app/internal/middleware"
	"hrm-app/internal/pkg/database"
	rmqManager "hrm-app/internal/pkg/rabbitmq/manager"
	"hrm-app/internal/websocket"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, channelManager *rmqManager.ChannelManager, rateLimiter *middleware.RateLimiter) (*gin.Engine, *websocket.Hub) {
	r := gin.Default()

	// CORS middleware (if needed)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Use Prometheus middleware
	r.Use(middleware.PrometheusMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// COMMENTED: Initialize Kafka - Replaced with RabbitMQ
	// kafkautil.InitKafka(cfg.Kafka.Brokers)

	// Initialize WebSocket Hub with RabbitMQ
	rabbitmqURL := "amqp://appuser:strongpassword@localhost:5672/" // TODO: Move to config
	hub := websocket.NewHub(database.RDB, rabbitmqURL, channelManager, rateLimiter)
	go hub.Run()

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"rabbitmq": gin.H{
					"activeChannels": channelManager.GetActiveCount(),
				},
			})
		})

		// Monitoring endpoints
		monitoring := api.Group("/monitoring")
		{
			monitoring.GET("/active-users", func(c *gin.Context) {
				stats := channelManager.GetAllStats()
				c.JSON(200, gin.H{
					"totalUsers": len(stats),
					"users":      stats,
				})
			})

			monitoring.GET("/user/:userId", func(c *gin.Context) {
				userID := c.Param("userId")
				uc, exists := channelManager.GetUserChannel(userID)
				if !exists {
					c.JSON(404, gin.H{"error": "user not found"})
					return
				}

				// Use the thread-safe GetStats method
				stats := uc.GetStats()
				c.JSON(200, stats)
			})
		}

		// Initialize Repositories
		userRepo := user.NewRepository()
		contactRepo := contact.NewRepository()
		workspaceRepo := workspaces.NewRepository()
		boardsRepo := boards.NewRepository()
		taskTabRepo := taskTab.NewRepository()
		taskCardRepo := taskCard.NewRepository()
		labelsRepo := labels.NewRepository()
		boardsUsersRepo := boardsUsers.NewRepository()
		taskCardCommentRepo := taskCardComment.NewRepository()
		taskCardUsersRepo := taskCardUsers.NewRepository()
		workspacesUsersRepo := workspacesUsers.NewRepository()
		roomChatRepo := room_chats.NewRepository()
		roomUserRepo := roomUsers.NewRepository()
		roomMessageRepo := room_messages.NewRepository()

		// Initialize UseCases
		// Initialize UseCases
		storageRepo, err := supabase.NewSupabaseStorageRepository(cfg)
		if err != nil {
			panic("Failed to initialize Supabase storage repository: " + err.Error())
		}

		uploadService := storage.NewService(storageRepo)

		userUseCase := user.NewUseCase(userRepo, contactRepo, uploadService)
		workspaceUseCase := workspaces.NewUseCase(workspaceRepo, workspacesUsersRepo, cfg)
		boardsUseCase := boards.NewUseCase(boardsRepo, taskTabRepo, taskCardRepo, boardsUsersRepo, labelsRepo, taskCardUsersRepo)
		taskTabUseCase := taskTab.NewUseCase(taskTabRepo)
		taskCardUseCase := taskCard.NewUseCase(taskCardRepo)
		labelsUseCase := labels.NewUseCase(labelsRepo)
		taskCardCommentUseCase := taskCardComment.NewUseCase(taskCardCommentRepo)
		taskCardUsersUseCase := taskCardUsers.NewUseCase(taskCardUsersRepo)
		workspaceRepoAdapter := workspaces.NewRepositoryAdapter(workspaceRepo)
		workspacesUsersUseCase := workspacesUsers.NewUseCase(workspacesUsersRepo, workspaceRepoAdapter, cfg)
		boardRepoAdapter := boards.NewRepositoryAdapter(boardsRepo)
		boardWorkspaceRepoAdapter := workspaces.NewBoardWorkspaceRepositoryAdapter(workspaceRepo)
		boardsUsersUseCase := boardsUsers.NewUseCase(boardsUsersRepo, boardRepoAdapter, boardWorkspaceRepoAdapter, cfg)
		roomChatUseCase := room_chats.NewUseCase(roomChatRepo, uploadService, cfg.Supabase.S3.Bucket)
		roomUserUseCase := roomUsers.NewUseCase(roomUserRepo)
		roomMessageUseCase := room_messages.NewUseCase(roomMessageRepo)

		// Initialize Handlers
		userHandler := user.NewHandler(userUseCase)
		workspaceHandler := workspaces.NewHandler(workspaceUseCase)
		boardsHandler := boards.NewHandler(boardsUseCase)
		taskTabHandler := taskTab.NewHandler(taskTabUseCase)
		taskCardHandler := taskCard.NewHandler(taskCardUseCase)
		labelsHandler := labels.NewHandler(labelsUseCase)
		taskCardCommentHandler := taskCardComment.NewHandler(taskCardCommentUseCase)
		taskCardUsersHandler := taskCardUsers.NewHandler(taskCardUsersUseCase)
		workspacesUsersHandler := workspacesUsers.NewHandler(workspacesUsersUseCase)
		boardsUsersHandler := boardsUsers.NewHandler(boardsUsersUseCase)
		roomChatHandler := room_chats.NewHandler(roomChatUseCase)
		roomUserHandler := roomUsers.NewHandler(roomUserUseCase)

		// Contact UseCase and Handler
		contactUseCase := contact.NewUseCase(contactRepo, storageRepo)
		contactHandler := contact.NewHandler(contactUseCase, cfg.Supabase.S3.Bucket)

		// WebSocket handler
		wsHandler := websocket.NewHandler(hub, taskCardUseCase, taskTabUseCase, taskCardCommentUseCase, labelsUseCase, taskCardUsersUseCase, boardsUsersUseCase, workspacesUsersUseCase, boardsUseCase, roomMessageUseCase, roomChatUseCase, roomUserUseCase, contactUseCase, userUseCase)

		// auth handler needs repo + cfg
		authHandler := auth.NewHandler(userRepo, cfg)

		// Note: NewSupabaseStorageRepository creates its own client internally in current implementation
		// Ideally we should inject the client if following the user's manual wiring request exactly,
		// but complying with the existing codebase structure.

		storageHandler := storage.NewHandler(storageRepo, cfg.Supabase.S3.Bucket)

		api.POST("/login", authHandler.Login)

		// Upload route
		api.POST("/upload", middleware.AuthMiddleware(cfg), storageHandler.UploadFile)
		api.POST("/logout", authHandler.Logout)
		api.POST("/refresh-token", authHandler.RefreshToken)

		user := api.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.GET("/", userHandler.GetAll)
			user.GET("/:id", userHandler.GetByID)
			user.PUT("/:id", userHandler.Update)
			user.DELETE("/:id", userHandler.Delete)
		}

		contacts := api.Group("/contacts")
		{
			protected := contacts.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.GET("/me", contactHandler.GetMyContact)
				protected.PUT("/me", contactHandler.UpdateMyContact)
			}
		}

		workspace := api.Group("/workspaces")
		{
			// workspace.GET("/", workspaceHandler.GetAll)

			protected := workspace.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", workspaceHandler.Create)
				protected.GET("/", workspaceHandler.GetByUserID)
				protected.GET("/guest", workspaceHandler.GetGuestWorkspaces)
				protected.GET("/:id", workspaceHandler.GetByID)
				protected.DELETE("/:id", workspaceHandler.Delete)
				protected.PUT("/:id", workspaceHandler.Update)
				protected.POST("/join", workspacesUsersHandler.Join)
				protected.GET("/:id/join-token", workspacesUsersHandler.GenerateJoinToken)
			}
		}

		boards := api.Group("/boards")
		{
			protected := boards.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", boardsHandler.CreateBoard)
				protected.GET("/", boardsHandler.GetByUserID)
				protected.GET("/:id", boardsHandler.GetBoardByID)
				protected.GET("/workspace/:workspace_id", boardsHandler.GetByWorkspaceID)
				protected.DELETE("/:id", boardsHandler.DeleteBoard)
				protected.PUT("/:id", boardsHandler.UpdateBoard)
				protected.POST("/join", boardsUsersHandler.Join)
				protected.GET("/:id/join-token", boardsUsersHandler.GenerateJoinToken)
				protected.GET("/:id/tabs", boardsHandler.GetBoardTabs)
				protected.GET("/tabs/:tab_id/cards", boardsHandler.GetTabCards)
			}
		}

		taskTab := api.Group("/task-tabs")
		{
			protected := taskTab.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", taskTabHandler.Create)
				protected.GET("/", taskTabHandler.GetAll)
				protected.GET("/:id", taskTabHandler.GetByID)
				protected.DELETE("/:id", taskTabHandler.Delete)
				protected.PUT("/:id", taskTabHandler.Update)
			}
		}

		taskCard := api.Group("/task-cards")
		{
			protected := taskCard.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", taskCardHandler.Create)
				protected.GET("/", taskCardHandler.GetAll)
				protected.GET("/:id", taskCardHandler.GetByID)
				protected.GET("/task-tab/:task_tab_id", taskCardHandler.GetByTaskTabID)
				protected.DELETE("/:id", taskCardHandler.Delete)
				protected.PUT("/:id", taskCardHandler.Update)
			}
		}

		labels := api.Group("/labels")
		{
			protected := labels.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", labelsHandler.Create)
				protected.GET("/", labelsHandler.GetAll)
				protected.GET("/:id", labelsHandler.GetByID)
				protected.GET("/task-card/:task_card_id", labelsHandler.GetByTaskCardID)
				protected.DELETE("/:id", labelsHandler.Delete)
				protected.PUT("/:id", labelsHandler.Update)
			}
		}

		taskCardComment := api.Group("/task-card-comments")
		{
			protected := taskCardComment.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", taskCardCommentHandler.CreateTaskCardComment)
				protected.GET("/", taskCardCommentHandler.GetAllTaskCardComment)
				protected.GET("/:id", taskCardCommentHandler.GetTaskCardCommentByID)
				protected.GET("/task-card/:task_card_id", taskCardCommentHandler.GetTaskCardCommentByTaskCardID)
				protected.DELETE("/:id", taskCardCommentHandler.DeleteTaskCardComment)
				protected.PUT("/:id", taskCardCommentHandler.UpdateTaskCardComment)
			}
		}

		taskCardUsers := api.Group("/task-card-users")
		{
			protected := taskCardUsers.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", taskCardUsersHandler.CreateTaskCardUser)
				protected.GET("/task-card/:task_card_id", taskCardUsersHandler.GetTaskCardUserByTaskCardID)
				protected.PUT("/:id", taskCardUsersHandler.Update)
				protected.DELETE("/:id", taskCardUsersHandler.Delete)
			}
		}

		workspacesUsers := api.Group("/workspaces-users")
		{
			protected := workspacesUsers.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", workspacesUsersHandler.Create)
				protected.GET("/workspace/:workspace_id", workspacesUsersHandler.GetByWorkspaceID)
				protected.GET("/user", workspacesUsersHandler.GetByUserID)
				protected.GET("/:id", workspacesUsersHandler.GetByID)
				protected.PUT("/:id", workspacesUsersHandler.Update)
				protected.DELETE("/:id", workspacesUsersHandler.Delete)
			}
		}

		roomChats := api.Group("/room-chats")
		{
			protected := roomChats.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", roomChatHandler.Create)
				protected.POST("/upload", roomChatHandler.UploadAttachment)
				protected.GET("/workspace/:workspace_id", roomChatHandler.GetByWorkspaceID)
				protected.GET("/:id", roomChatHandler.GetByID)
				protected.PUT("/:id", roomChatHandler.Update)
				protected.DELETE("/:id", roomChatHandler.Delete)
			}
		}

		roomUsers := api.Group("/room-users")
		{
			protected := roomUsers.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/join", roomUserHandler.Join)
				protected.GET("/room/:room_id", roomUserHandler.GetUsersByRoom)
			}
		}

		boardsUsers := api.Group("/boards-users")
		{
			protected := boardsUsers.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", boardsUsersHandler.Create)
				protected.GET("/board/:board_id", boardsUsersHandler.GetByBoardID)
				protected.GET("/user", boardsUsersHandler.GetByUserID)
				protected.GET("/:id", boardsUsersHandler.GetByID)
				protected.PUT("/:id", boardsUsersHandler.Update)
				protected.DELETE("/:id", boardsUsersHandler.Delete)
			}
		}

		// WebSocket routes - use WebSocket-specific auth middleware
		ws := api.Group("/ws")
		{
			ws.Use(middleware.AuthMiddlewareWS(cfg))
			ws.GET("/task-cards", wsHandler.HandleWebSocket)
			ws.PUT("/task-cards", wsHandler.HandleWebSocket)
			ws.GET("/clients", wsHandler.GetConnectedClients)
		}
	}

	return r, hub
}
