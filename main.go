// package main

// import (
// 	"bytes"
// 	"crypto/x509"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hyperledger/fabric-gateway/pkg/client"
// 	"github.com/hyperledger/fabric-gateway/pkg/identity"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// )

// const (
// 	mspID        = "Org1MSP"
// 	cryptoPath   = "../test-network/organizations/peerOrganizations/org1.example.com"
// 	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
// 	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
// 	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
// 	peerEndpoint = "dns:///localhost:7051"
// 	gatewayPeer  = "peer0.org1.example.com"
// )

// // newGrpcConnection creates a gRPC connection to the Gateway server.
// func newGrpcConnection() *grpc.ClientConn {
// 	certificatePEM, err := os.ReadFile(tlsCertPath)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
// 	}

// 	certificate, err := identity.CertificateFromPEM(certificatePEM)
// 	if err != nil {
// 		panic(err)
// 	}

// 	certPool := x509.NewCertPool()
// 	certPool.AddCert(certificate)
// 	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

// 	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
// 	if err != nil {
// 		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
// 	}

// 	return connection
// }

// // func readFirstFile(dirPath string) ([]byte, error) {
// // 	dir, err := os.Open(dirPath)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	fileNames, err := dir.Readdirnames(1)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	return os.ReadFile(path.Join(dirPath, fileNames[0]))
// // }

// // func newIdentity() *identity.X509Identity {
// // 	certificatePEM, err := readFirstFile(certPath)
// // 	if err != nil {
// // 		panic(fmt.Errorf("failed to read certificate file: %w", err))
// // 	}

// // 	certificate, err := identity.CertificateFromPEM(certificatePEM)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	id, err := identity.NewX509Identity(mspID, certificate)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	return id
// // }

// func formatJSON(data []byte) string {
// 	var prettyJSON bytes.Buffer
// 	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
// 		panic(fmt.Errorf("failed to parse JSON: %w", err))
// 	}
// 	return prettyJSON.String()
// }

// // newSign creates a function that generates a digital signature from a message digest using a private key.
// // func newSign() identity.Sign {
// // 	privateKeyPEM, err := readFirstFile(keyPath)
// // 	if err != nil {
// // 		panic(fmt.Errorf("failed to read private key file: %w", err))
// // 	}

// // 	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	sign, err := identity.NewPrivateKeySign(privateKey)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	return sign
// // }

// // newIdentityFromPEM creates a client identity using the uploaded X.509 certificate.
// func newIdentityFromPEM(certificatePEM []byte) *identity.X509Identity {
// 	certificate, err := identity.CertificateFromPEM(certificatePEM)
// 	if err != nil {
// 		panic(err)
// 	}

// 	id, err := identity.NewX509Identity(mspID, certificate)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return id
// }

// // newSignFromPEM creates a function that generates a digital signature from a message digest using the uploaded private key.
// func newSignFromPEM(privateKeyPEM []byte) identity.Sign {
// 	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
// 	if err != nil {
// 		panic(err)
// 	}

// 	sign, err := identity.NewPrivateKeySign(privateKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return sign
// }

// func main() {
// 	router := gin.Default()

// 	// Define the API route to get all assets
// 	router.GET("/assets", func(c *gin.Context) {
// 		// Parse the certificate and key files from form data
// 		certFile, certErr := c.FormFile("cert")
// 		keyFile, keyErr := c.FormFile("key")
// 		if certErr != nil || keyErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Cert and key files are required"})
// 			return
// 		}

// 		// Open the uploaded files
// 		certContent, certErr := certFile.Open()
// 		if certErr != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open cert file"})
// 			return
// 		}
// 		defer certContent.Close()

// 		keyContent, keyErr := keyFile.Open()
// 		if keyErr != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open key file"})
// 			return
// 		}
// 		defer keyContent.Close()

// 		// Read the content of the cert and key files
// 		certificatePEM := make([]byte, certFile.Size)
// 		_, certErr = certContent.Read(certificatePEM)
// 		if certErr != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read cert file"})
// 			return
// 		}

// 		privateKeyPEM := make([]byte, keyFile.Size)
// 		_, keyErr = keyContent.Read(privateKeyPEM)
// 		if keyErr != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read key file"})
// 			return
// 		}

// 		start := time.Now()
// 		// Connect to the gateway
// 		clientConnection := newGrpcConnection()
// 		defer clientConnection.Close()

// 		id := newIdentityFromPEM(certificatePEM)
// 		sign := newSignFromPEM(privateKeyPEM)

// 		gw, err := client.Connect(
// 			id,
// 			client.WithSign(sign),
// 			client.WithClientConnection(clientConnection),
// 			client.WithEvaluateTimeout(5*time.Second),
// 			client.WithEndorseTimeout(15*time.Second),
// 			client.WithSubmitTimeout(5*time.Second),
// 			client.WithCommitStatusTimeout(1*time.Minute),
// 		)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		defer gw.Close()

// 		// Use environment variables for chaincode and channel names
// 		chaincodeName := "basic"
// 		if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
// 			chaincodeName = ccname
// 		}

// 		channelName := "mychannel"
// 		if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
// 			channelName = cname
// 		}

// 		network := gw.GetNetwork(channelName)
// 		contract := network.GetContract(chaincodeName)

// 		// Retrieve all assets
// 		result, err := contract.EvaluateTransaction("ReadAsset", "asset1")
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %s", err)})
// 			return
// 		}

// 		fmt.Printf("Time to get by ID: %v\n", time.Since(start))

// 		// Format the result as pretty JSON
// 		formattedResult := formatJSON(result)
// 		c.Data(http.StatusOK, "application/json", []byte(formattedResult))
// 	})

// 	// Start the Gin server
// 	router.Run(":8080")
// }

package main

import (
	"example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/assets", controllers.GetAllAssets)
	router.GET("/assets/:assetID", controllers.GetAssetByID)
	router.POST("/assets", controllers.CreateAsset)

	// Start the Gin server
	router.Run(":8080")
}
