package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/guish/university-social/internal/config"
	"github.com/guish/university-social/internal/handler"
	"github.com/guish/university-social/internal/middleware"
	"github.com/guish/university-social/internal/repository"
)

func Setup(db *sql.DB, cfg *config.Config) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	postHandler := handler.NewPostHandler(postRepo)
	authHandler := handler.NewAuthHandler(userRepo, cfg)
	friendRepo := repository.NewFriendRepository(db)
	friendHandler := handler.NewFriendHandler(friendRepo)
	circleRepo := repository.NewCircleRepository(db)
	circleHandler := handler.NewCircleHandler(circleRepo)
	adminHandler := handler.NewAdminHandler(db)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.WxLogin)

		auth := api.Group("")
		auth.Use(middleware.AuthRequired(cfg))
		{
			auth.GET("/user/profile", authHandler.GetProfile)
			auth.PUT("/user/profile", authHandler.UpdateProfile)
			auth.POST("/posts", postHandler.CreatePost)
			auth.GET("/posts", postHandler.ListPosts)
			auth.GET("/posts/:id", postHandler.GetPost)
			auth.DELETE("/posts/:id", postHandler.DeletePost)
			auth.POST("/posts/:id/comments", postHandler.CreateComment)
			auth.GET("/posts/:id/comments", postHandler.ListComments)
			auth.POST("/posts/:id/like", postHandler.ToggleLike)
			auth.GET("/topics", postHandler.ListTopics)
			auth.GET("/users/search", friendHandler.SearchUsers)
			auth.GET("/friends", friendHandler.ListFriends)
			auth.GET("/friends/requests", friendHandler.ListRequests)
			auth.POST("/friends/request/:id", friendHandler.SendRequest)
			auth.POST("/friends/accept/:id", friendHandler.AcceptRequest)
			auth.DELETE("/friends/:id", friendHandler.RemoveFriend)
			auth.GET("/notifications", friendHandler.ListNotifications)
			auth.POST("/notifications/read", friendHandler.MarkRead)
			auth.GET("/circles", circleHandler.List)
			auth.POST("/circles", circleHandler.Create)
			auth.GET("/circles/:id", circleHandler.GetByID)
			auth.PUT("/circles/:id", circleHandler.Update)
			auth.POST("/circles/:id/join", circleHandler.Join)
			auth.POST("/circles/:id/leave", circleHandler.Leave)
			auth.GET("/circles/:id/members", circleHandler.ListMembers)
			auth.GET("/circles/:id/posts", circleHandler.ListPosts)
			auth.POST("/circles/:id/posts", circleHandler.CreatePost)
			auth.POST("/circles/:postId/like", circleHandler.TogglePostLike)
			auth.POST("/circles/:postId/comments", circleHandler.CreateComment)
			auth.GET("/circles/:postId/comments", circleHandler.ListComments)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired(cfg))
		admin.Use(middleware.AdminRequired(cfg))
		{
			admin.GET("/dashboard", adminHandler.Dashboard)
			admin.GET("/users", adminHandler.ListUsers)
			admin.PUT("/users/:id/ban", adminHandler.BanUser)
			admin.PUT("/users/:id/unban", adminHandler.UnbanUser)
			admin.GET("/posts", adminHandler.ListPosts)
			admin.DELETE("/posts/:id", adminHandler.DeletePost)
			admin.GET("/circles", adminHandler.ListCircles)
			admin.DELETE("/circles/:id", adminHandler.DeleteCircle)
		}
	}

	return r
}
