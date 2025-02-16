package grpcclient

import (
	"context"
	"log"
	"net/http"

	pb "code-sandbox/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func ExecuteCode(c *gin.Context) {
	type ReqBody struct {
		Code     string `json:"code" binding:"required"`
		Language string `json:"language" binding:"required"`
		Input    string `json:"input"`
	}
	log.Println("Executing Code .... ")

	var req ReqBody

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal("Required Fields are missing ")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields are missing."})
		return
	}

	// create grpc connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to connect grpc " + err.Error()})
		c.Abort()
		return
	}

	defer conn.Close()

	// create gRPC client
	client := pb.NewCodeExecutionServiceClient(conn)

	// call grpc method
	resp, err := client.ExecuteCode(context.Background(), &pb.ExecutionRequest{
		Language: req.Language,
		Code:     req.Code,
		Input:    req.Input,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while executing Code " + err.Error()})
		return
	}
	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": string(resp.Output), "error": resp.Error})
}
