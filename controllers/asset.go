package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/m/model"
	"example.com/m/service"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

func GetAllAssets(c *gin.Context) {
	certificatePEM, privateKeyPEM, err := getCertificateAndPrivateKeyFromForm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chaincode, err := service.NewChaincodeService(certificatePEM, privateKeyPEM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create chaincode service: %s", err)})
		return
	}

	result, err := chaincode.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
		return
	}

	formattedResult := utils.FormatJSON(result)
	c.Data(http.StatusOK, "application/json", []byte(formattedResult))
}

func GetAssetByID(c *gin.Context) {
	assetID := c.Param("assetID")
	certificatePEM, privateKeyPEM, err := getCertificateAndPrivateKeyFromForm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	chaincode, err := service.NewChaincodeService(certificatePEM, privateKeyPEM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create chaincode service: %s", err)})
		return
	}

	result, err := chaincode.ReadAssetByID(assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
		return
	}

	fmt.Printf("GetAssetByID took %s\n", time.Since(start))

	formattedResult := utils.FormatJSON(result)
	c.Data(http.StatusOK, "application/json", []byte(formattedResult))
}

func CreateAsset(c *gin.Context) {
	certificatePEM, privateKeyPEM, err := getCertificateAndPrivateKeyFromForm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	chaincode, err := service.NewChaincodeService(certificatePEM, privateKeyPEM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create chaincode service: %s", err)})
		return
	}

	size, _ := strconv.Atoi(c.PostForm("size"))
	appraisedValue, _ := strconv.Atoi(c.PostForm("appraisedValue"))
	asset := model.Asset{
		ID:             c.PostForm("assetID"),
		Color:          c.PostForm("color"),
		Size:           size,
		Owner:          c.PostForm("owner"),
		AppraisedValue: appraisedValue,
	}

	err = chaincode.CreateAsset(asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
		return
	}

	fmt.Printf("CreateAsset took %s\n", time.Since(start))

	c.JSON(http.StatusOK, gin.H{"message": "Asset created successfully"})
}

func getCertificateAndPrivateKeyFromForm(c *gin.Context) ([]byte, []byte, error) {
	certFile, certErr := c.FormFile("cert")
	keyFile, keyErr := c.FormFile("key")
	if certErr != nil || keyErr != nil {
		return nil, nil, fmt.Errorf("cert and key files are required")
	}

	certContent, certErr := certFile.Open()
	if certErr != nil {
		return nil, nil, fmt.Errorf("unable to open cert file")
	}
	defer certContent.Close()

	keyContent, keyErr := keyFile.Open()
	if keyErr != nil {
		return nil, nil, fmt.Errorf("unable to open key file")
	}
	defer keyContent.Close()

	certificatePEM := make([]byte, certFile.Size)
	_, certErr = certContent.Read(certificatePEM)
	if certErr != nil {
		return nil, nil, fmt.Errorf("unable to read cert file")
	}

	privateKeyPEM := make([]byte, keyFile.Size)
	_, keyErr = keyContent.Read(privateKeyPEM)
	if keyErr != nil {
		return nil, nil, fmt.Errorf("unable to read key file")
	}

	return certificatePEM, privateKeyPEM, nil
}
