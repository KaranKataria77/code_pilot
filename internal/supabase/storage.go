package supabase

import (
	"bytes"
	"code_pilot/internal/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	storage_go "github.com/supabase-community/storage-go"
)

func FileUpload(c *gin.Context) {
	config.LoadEnv()
	content := []byte("Name=Karan bhanushali")
	file := bytes.NewReader(content)

	reference_id := config.GetEnv("SUPABASE_PROJECT_ID", "")
	api_key := config.GetEnv("SUPABASE_PROJECT_SERVICE_API_KEY", "")
	s := storage_go.NewClient("https://"+reference_id+".supabase.co/storage/v1", api_key, nil)
	resp, err := s.UploadFile("code_pilot", "test.txt", file)
	log.Println("File writing ", resp, err)

}

func FileDownload(c *gin.Context) {
	config.LoadEnv()

	reference_id := config.GetEnv("SUPABASE_PROJECT_ID", "")
	api_key := config.GetEnv("SUPABASE_PROJECT_SERVICE_API_KEY", "")
	s := storage_go.NewClient("https://"+reference_id+".supabase.co/storage/v1", api_key, nil)
	resp, err := s.DownloadFile("code_pilot", "/base_nextjs_page_router_js")
	if err != nil {
		log.Println("Error while downloading file ")
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error while downloading file " + err.Error()})
		return
	}
	log.Println("Response ", resp)
}

func FileLists(fileName string, s *storage_go.Client) ([]storage_go.FileObject, error) {
	resp, err := s.ListFiles("code_pilot", fileName, storage_go.FileSearchOptions{
		Limit:  100,
		Offset: 0,
		SortByOptions: storage_go.SortBy{
			Column: "",
			Order:  "",
		},
	})
	if err != nil {
		log.Println("Error while reading list in files " + err.Error())
		return nil, err
	}
	return resp, nil
}

func GetFiles(c *gin.Context) {
	config.LoadEnv()

	reference_id := config.GetEnv("SUPABASE_PROJECT_ID", "")
	api_key := config.GetEnv("SUPABASE_PROJECT_SERVICE_API_KEY", "")
	s := storage_go.NewClient("https://"+reference_id+".supabase.co/storage/v1", api_key, nil)
	fileName := "base_nextjs_page_router_js/src"
	resp, err := s.ListFiles("code_pilot", fileName, storage_go.FileSearchOptions{
		Limit:  1000,
		Offset: 0,
		SortByOptions: storage_go.SortBy{
			Column: "",
			Order:  "",
		},
	})
	if err != nil {
		log.Println("Error while listing files/folder " + err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"response": resp})
}

func GetFolderStructure(c *gin.Context) {
	config.LoadEnv()
	reference_id := config.GetEnv("SUPABASE_PROJECT_ID", "")
	api_key := config.GetEnv("SUPABASE_PROJECT_SERVICE_API_KEY", "")
	s := storage_go.NewClient("https://"+reference_id+".supabase.co/storage/v1", api_key, nil)
	parentFolder := "base_nextjs_page_router_js"
	q := []string{parentFolder}
	// response := make(map[string]map[string]interface{})
	children := []string{}
	for len(q) > 0 {
		folderName := q[0]
		q = q[1:]
		resp, err := FileLists(folderName, s)
		if err != nil {
			log.Println("Error while reading Folder/file " + err.Error())
			break
		}
		// response[folderName] = map[string]interface{}{
		// 	"child":[]string{},
		// }

		for i := 0; i < len(resp); i++ {
			q = append(q, folderName+"/"+resp[i].Name)
			// children := response[folderName]["child"].([]string)
			children = append(children, folderName+"/"+resp[i].Name)
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": children})
}
