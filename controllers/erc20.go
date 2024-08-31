package controllers

import (
	"fmt"
	"net/http"

	"example.com/m/service"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

func GetBalance(c *gin.Context) {
	certificatePEM, privateKeyPEM, err := getCertificateAndPrivateKeyFromForm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ccname := "token_erc20"
	chaincode, err := service.NewChaincodeService(certificatePEM, privateKeyPEM, &ccname, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create chaincode service: %s", err)})
		return
	}

	result, err := chaincode.GetBalance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
		return
	}

	formattedResult := utils.FormatJSON(result)
	c.Data(http.StatusOK, "application/json", []byte(formattedResult))
}

func Transfer(c *gin.Context) {
	certificatePEM, privateKeyPEM, err := getCertificateAndPrivateKeyFromForm(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ccname := "token_erc20"
	chaincode, err := service.NewChaincodeService(certificatePEM, privateKeyPEM, &ccname, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create chaincode service: %s", err)})
		return
	}

	recipientCN := c.PostForm("recipientCN")
	if recipientCN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipientCN is required"})
		return
	}

	amount := c.PostForm("amount")
	if amount == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount is required"})
		return
	}

	err = chaincode.Transfer(recipientCN, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}
