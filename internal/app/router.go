package app

import (
	"hrm-app/config"
	"hrm-app/internal/domain/auth"
	"hrm-app/internal/domain/boards"
	"hrm-app/internal/domain/contact"
	"hrm-app/internal/domain/labels"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskCardComment"
	"hrm-app/internal/domain/taskTab"
	"hrm-app/internal/domain/user"
	"hrm-app/internal/domain/workspaces"
	"hrm-app/internal/middleware"
	"hrm-app/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Use Prometheus middleware
	r.Use(middleware.PrometheusMiddleware())

	// Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// User routes
		userRepo := user.NewRepository()
		contactRepo := contact.NewRepository()
		userUseCase := user.NewUseCase(userRepo, contactRepo)
		userHandler := user.NewHandler(userUseCase)

		// Workspace routes
		workspaceRepo := workspaces.NewRepository()
		workspaceUseCase := workspaces.NewUseCase(workspaceRepo)
		workspaceHandler := workspaces.NewHandler(workspaceUseCase)

		// Initialize Repositories
		boardsRepo := boards.NewRepository()
		taskTabRepo := taskTab.NewRepository()
		taskCardRepo := taskCard.NewRepository()
		labelsRepo := labels.NewRepository()

		// Boards routes
		boardsUseCase := boards.NewUseCase(boardsRepo, taskTabRepo, taskCardRepo)
		boardsHandler := boards.NewHandler(boardsUseCase)

		// TaskTab routes
		taskTabUseCase := taskTab.NewUseCase(taskTabRepo)
		taskTabHandler := taskTab.NewHandler(taskTabUseCase)

		// TaskCard routes
		taskCardUseCase := taskCard.NewUseCase(taskCardRepo, labelsRepo)
		taskCardHandler := taskCard.NewHandler(taskCardUseCase)

		// Labels routes
		labelsUseCase := labels.NewUseCase(labelsRepo)
		labelsHandler := labels.NewHandler(labelsUseCase)

		// TaskCardComment routes
		taskCardCommentRepo := taskCardComment.NewRepository()
		taskCardCommentUseCase := taskCardComment.NewUseCase(taskCardCommentRepo)
		taskCardCommentHandler := taskCardComment.NewHandler(taskCardCommentUseCase)

		// WebSocket handler
		wsHandler := websocket.NewHandler(hub, taskCardUseCase, taskTabUseCase, taskCardCommentUseCase)

		// auth handler needs repo + cfg
		authHandler := auth.NewHandler(userRepo, cfg)

		api.POST("/login", authHandler.Login)
		api.POST("/logout", authHandler.Logout)
		api.POST("/refresh-token", authHandler.RefreshToken)

		user := api.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.GET("/", userHandler.GetAll)
			user.GET("/:id", userHandler.GetByID)
			user.DELETE("/:id", userHandler.Delete)
		}

		workspace := api.Group("/workspaces")
		{
			// workspace.GET("/", workspaceHandler.GetAll)

			protected := workspace.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", workspaceHandler.Create)
				protected.GET("/", workspaceHandler.GetByUserID)
				protected.GET("/:id", workspaceHandler.GetByID)
				protected.DELETE("/:id", workspaceHandler.Delete)
				protected.PUT("/:id", workspaceHandler.Update)
			}
		}

		boards := api.Group("/boards")
		{
			protected := boards.Group("/")
			protected.Use(middleware.AuthMiddleware(cfg))
			{
				protected.POST("/", boardsHandler.CreateBoard)
				protected.GET("/", boardsHandler.GetAllBoard)
				protected.GET("/:id", boardsHandler.GetBoardByID)
				protected.GET("/workspace/:workspace_id", boardsHandler.GetByWorkspaceID)
				protected.DELETE("/:id", boardsHandler.DeleteBoard)
				protected.PUT("/:id", boardsHandler.UpdateBoard)
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

		// WebSocket routes
		ws := api.Group("/ws")
		{
			ws.GET("/task-cards", wsHandler.HandleWebSocket)
			ws.PUT("/task-cards", wsHandler.HandleWebSocket)
			ws.GET("/clients", wsHandler.GetConnectedClients)
		}
	}

	return r
}
