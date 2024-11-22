package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	response_models "caapp-server/src/models/responce_models"

	"github.com/gin-gonic/gin"
)

func GetLanguageDataList(c *gin.Context) {
	filePath := "src/constants/language/language.json"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "loi mo file language data"})
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "loi doc file language data"})
		return
	}

	var languages []response_models.GetLanguageDataListResponseItem
	if err := json.Unmarshal(content, &languages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "loi paste file language data"})
		return
	}

	res := response_models.GetLanguageDataListResponse{
		Languages: languages,
	}

	c.JSON(http.StatusOK, res)
}
