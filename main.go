package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

import _ "github.com/muxinyu1/kernel-api/docs"

// @title Kernel Service API
// @version 1.0
// @description API for managing the kernel service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1

func main() {
	router := gin.Default()

	// Swagger 文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由
	apiV1 := router.Group("/v1") 
	{
		apiV1.POST("/start", start)
        apiV1.POST("/stop", stop)
        apiV1.POST("/restart", restart)
        apiV1.GET("/status", status)
        apiV1.POST("/set-domain", setDomain)
        apiV1.POST("/set-ip", setIp)
	}
	// 启动服务
	router.Run(":8080")
}

// @Summary Start the kernel service
// @Description Starts the kernel service
// @Tags Kernel
// @Produce json
// @Success 200 {object} map[string]string
// @Router /start [post]
func start(c *gin.Context) {
	cmd := exec.Command("bash", "-c", "kernel", "start")
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kernel service started",
	})
}

// @Summary Stop the kernel service
// @Description Stops the kernel service
// @Tags Kernel
// @Produce json
// @Success 200 {object} map[string]string
// @Router /stop [post]
func stop(c *gin.Context) {
	cmd := exec.Command("bash", "-c", "kernel", "stop")
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kernel service stopped",
	})
}

// @Summary Restart the kernel service
// @Description Restarts the kernel service
// @Tags Kernel
// @Produce json
// @Success 200 {object} map[string]string
// @Router /restart [post]
func restart(c *gin.Context) {
	cmd := exec.Command("bash", "-c", "kernel", "restart")
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kernel service restarted",
	})
}

// @Summary Get kernel service status
// @Description Returns the current status of the kernel service
// @Tags Kernel
// @Produce json
// @Success 200 {object} map[string]string
// @Router /status [get]
func status(c *gin.Context) {
	cmd := exec.Command("bash", "-c", "kernel", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": string(output) + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": string(output),
	})
}

// @Summary Set domain for the kernel service
// @Description Sets the domain for the kernel service
// @Tags Configuration
// @Accept json
// @Produce json
// @Param domain body string true "Domain to set"
// @Success 200 {object} map[string]string
// @Router /set-domain [post]
func setDomain(c *gin.Context) {
	var req struct {
		Domain string `json:"domain"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}
	if req.Domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Domain is required",
		})
		return
	}
	cmd := exec.Command("bash", "-c", "kernel", "--set-domain", req.Domain)
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Domain set to %s", req.Domain),
	})
}

// @Summary Set IP plan for the kernel service
// @Description Sets the IP plan for the kernel service
// @Tags Configuration
// @Accept json
// @Produce json
// @Param ip_plan body string true "IP plan to set"
// @Success 200 {object} map[string]string
// @Router /set-ip [post]
func setIp(c *gin.Context) {
	var req struct {
		IPPlan string `json:"ipv4"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}
	if req.IPPlan == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "IP plan is required",
		})
		return
	}
	cmd := exec.Command("bash", "-c", "kernel", "--set-ip", req.IPPlan)
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("IP plan set to %s", req.IPPlan),
	})
}