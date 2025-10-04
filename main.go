package main

import (
	"log"
	"github.com/vishaltalsaniya-7/voip-api/config"
	"github.com/vishaltalsaniya-7/voip-api/database"
	"github.com/vishaltalsaniya-7/voip-api/manager"
	"github.com/vishaltalsaniya-7/voip-api/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	eslMgr := manager.NewESLManager(cfg.FreeSWITCH)

	go eslMgr.ListenEvents()

	callController := controller.NewCallController(eslMgr)
	cdrController := controller.NewCDRController(db)

	r := gin.Default()

	r.POST("/call", callController.InitiateCall)
	r.GET("/call/status/:uuid/", callController.GetCallStatus) 
	r.GET("/cdrs", cdrController.GetCDRs)

	// Start server
	log.Printf("Starting server on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
