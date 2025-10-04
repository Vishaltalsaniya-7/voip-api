package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vishaltalsaniya-7/voip-api/manager"
	"github.com/vishaltalsaniya-7/voip-api/request"
)

type CallController struct {
	eslMgr *manager.ESLManager
}

func NewCallController(eslMgr *manager.ESLManager) *CallController {
	return &CallController{
		eslMgr: eslMgr,
	}
}

func (cc *CallController) InitiateCall(c *gin.Context) {
	var req request.CallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	callID, err := cc.eslMgr.OriginateCall(req.Caller, req.Callee)
	if err != nil {
		log.Printf("Failed to originate call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"call_id": callID,
		"status":  "Call initiated",
	})
}


func (cc *CallController) GetCallStatus(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	

	status, err := cc.eslMgr.GetCallStatus(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}