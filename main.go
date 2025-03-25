package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

import _ "github.com/muxinyu1/kernel-api/docs"

// TaskStatus 表示任务状态
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusFinished  TaskStatus = "finished"
	StatusFailed    TaskStatus = "failed"
)

// Task 表示一个异步任务
type Task struct {
	ID        string     `json:"id"`
	Status    TaskStatus `json:"status"`
	Command   string     `json:"command"`
	Output    string     `json:"output,omitempty"`
	Error     string     `json:"error,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TaskStore 存储所有任务
type TaskStore struct {
	sync.RWMutex
	tasks map[string]*Task
}

var taskStore = TaskStore{
	tasks: make(map[string]*Task),
}

// @title Kernel Service API
// @version 1.0
// @description API for managing the kernel service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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
		apiV1.GET("/status", status)  // 改为异步版本
		apiV1.POST("/set-domain", setDomain)
		apiV1.POST("/set-ip", setIp)
		apiV1.GET("/task/:id", getTaskStatus)
	}
	// 启动服务
	router.Run(":8080")
}

// executeCommandAsync 异步执行命令
func executeCommandAsync(taskID, command string) {
	taskStore.Lock()
	task := taskStore.tasks[taskID]
	task.Status = StatusRunning
	task.UpdatedAt = time.Now()
	taskStore.Unlock()

	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()

	taskStore.Lock()
	defer taskStore.Unlock()
	task.UpdatedAt = time.Now()
	
	if err != nil {
		task.Status = StatusFailed
		task.Error = err.Error()
		task.Output = string(output)
	} else {
		task.Status = StatusFinished
		task.Output = string(output)
	}
}

// createTask 创建新任务
func createTask(command string) *Task {
	taskID := uuid.New().String()
	task := &Task{
		ID:        taskID,
		Status:    StatusPending,
		Command:   command,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	taskStore.Lock()
	taskStore.tasks[taskID] = task
	taskStore.Unlock()

	return task
}

// @Summary Get task status
// @Description Returns the current status of a task
// @Tags Task
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} Task
// @Router /task/{id} [get]
func getTaskStatus(c *gin.Context) {
	taskID := c.Param("id")

	taskStore.RLock()
	task, exists := taskStore.tasks[taskID]
	taskStore.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary Get kernel service status
// @Description Returns the current status of the kernel service (asynchronous)
// @Tags Kernel
// @Produce json
// @Success 202 {object} map[string]string
// @Router /status [get]
func status(c *gin.Context) {
	task := createTask("kernel status")
	go executeCommandAsync(task.ID, task.Command)
	
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"task_id": task.ID,
		"message": "Status check task created, use task_id to check result",
	})
}

// @Summary Start kernel service
// @Description Starts the kernel service asynchronously
// @Tags Kernel
// @Produce json
// @Success 202 {object} map[string]string
// @Router /start [post]
func start(c *gin.Context) {
	task := createTask("kernel start")
	go executeCommandAsync(task.ID, task.Command)
	
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"task_id": task.ID,
		"message": "Task created, use task_id to check status",
	})
}

// @Summary Stop kernel service
// @Description Stops the kernel service asynchronously
// @Tags Kernel
// @Produce json
// @Success 202 {object} map[string]string
// @Router /stop [post]
func stop(c *gin.Context) {
	task := createTask("kernel stop")
	go executeCommandAsync(task.ID, task.Command)
	
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"task_id": task.ID,
		"message": "Task created, use task_id to check status",
	})
}

// @Summary Restart kernel service
// @Description Restarts the kernel service asynchronously
// @Tags Kernel
// @Produce json
// @Success 202 {object} map[string]string
// @Router /restart [post]
func restart(c *gin.Context) {
	task := createTask("kernel restart")
	go executeCommandAsync(task.ID, task.Command)
	
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"task_id": task.ID,
		"message": "Task created, use task_id to check status",
	})
}

// @Summary Set domain for kernel service
// @Description Sets the domain for the kernel service asynchronously
// @Tags Configuration
// @Accept json
// @Produce json
// @Param domain body object{domain=string} true "Domain to set"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
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
	cmdStr := fmt.Sprintf("kernel --set-domain %s", req.Domain)
	task := createTask(cmdStr)
	go executeCommandAsync(task.ID, task.Command)
	
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"task_id": task.ID,
		"message": "Domain set task created, use task_id to check status",
	})
}

// @Summary Set IP plan for kernel service
// @Description Sets the IP plan for the kernel service asynchronously
// @Tags Configuration
// @Accept json
// @Produce json
// @Param ipv4 body object{ipv4=string} true "IPv4 address to set"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
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
	cmdStr := fmt.Sprintf("kernel --set-ip %s", req.IPPlan)
	task := createTask(cmdStr)
	go executeCommandAsync(task.ID, task.Command)
	
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "accepted",
		"task_id": task.ID,
		"message": "IP plan set task created, use task_id to check status",
	})
}