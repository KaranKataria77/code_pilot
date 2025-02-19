package grpcclient

import (
	pb "code-execution-sandbox/proto"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func CodeExecution(c *gin.Context) {
	folderName := "base_nextjs_page_router_js"

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while creating gRPC connection " + err.Error()})
		return
	}

	defer conn.Close()

	// create gRPC client
	client := pb.NewFileDownloadServiceClient(conn)

	// call gRPC method
	resp, err := client.DownloadFile(context.Background(), &pb.FileRequest{
		FolderName: folderName,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while creating gRPC client connection " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": resp.FilesDownloaded, "error": resp.Error})
}
