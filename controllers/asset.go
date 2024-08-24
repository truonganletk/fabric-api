package controllers

import (
	"fmt"
	"net/http"
	"time"

	"example.com/m/service"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

func GetAllAssets(c *gin.Context) {
	// Parse the certificate and key files from form data
	certFile, certErr := c.FormFile("cert")
	keyFile, keyErr := c.FormFile("key")
	if certErr != nil || keyErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cert and key files are required"})
		return
	}

	// Open the uploaded files
	certContent, certErr := certFile.Open()
	if certErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open cert file"})
		return
	}
	defer certContent.Close()

	keyContent, keyErr := keyFile.Open()
	if keyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open key file"})
		return
	}
	defer keyContent.Close()

	// Read the content of the cert and key files
	certificatePEM := make([]byte, certFile.Size)
	_, certErr = certContent.Read(certificatePEM)
	if certErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read cert file"})
		return
	}

	privateKeyPEM := make([]byte, keyFile.Size)
	_, keyErr = keyContent.Read(privateKeyPEM)
	if keyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read key file"})
		return
	}

	chaincode, err := service.NewChaincodeService(certificatePEM, privateKeyPEM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create chaincode service: %s", err)})
		return
	}

	// Retrieve all assets
	result, err := chaincode.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
		return
	}

	// Format the result as pretty JSON
	formattedResult := utils.FormatJSON(result)
	c.Data(http.StatusOK, "application/json", []byte(formattedResult))
}

func GetAssetByID(c *gin.Context) {
	assetID := c.Param("assetID")
	// Parse the certificate and key files from form data
	certFile, certErr := c.FormFile("cert")
	keyFile, keyErr := c.FormFile("key")
	if certErr != nil || keyErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cert and key files are required"})
		return
	}

	// Open the uploaded files
	certContent, certErr := certFile.Open()
	if certErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open cert file"})
		return
	}
	defer certContent.Close()

	keyContent, keyErr := keyFile.Open()
	if keyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open key file"})
		return
	}
	defer keyContent.Close()

	// Read the content of the cert and key files
	certificatePEM := make([]byte, certFile.Size)
	_, certErr = certContent.Read(certificatePEM)
	if certErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read cert file"})
		return
	}

	privateKeyPEM := make([]byte, keyFile.Size)
	_, keyErr = keyContent.Read(privateKeyPEM)
	if keyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read key file"})
		return
	}

	start := time.Now()
	chaincode, err := service.NewChaincodeService(certificatePEM, privateKeyPEM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create chaincode service: %s", err)})
		return
	}

	// Retrieve all assets
	result, err := chaincode.ReadAssetByID(assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
		return
	}

	fmt.Printf("GetAssetByID took %s\n", time.Since(start))

	// Format the result as pretty JSON
	formattedResult := utils.FormatJSON(result)
	c.Data(http.StatusOK, "application/json", []byte(formattedResult))
}
