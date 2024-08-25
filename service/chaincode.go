package service

import (
	"crypto/x509"
	"fmt"
	"os"
	"strconv"
	"time"

	"example.com/m/model"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ChaincodeService struct {
	client   *grpc.ClientConn
	gw       *client.Gateway
	network  *client.Network
	contract *client.Contract
}

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/Users/anle/Documents/Project/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentityFromPEM creates a client identity using the uploaded X.509 certificate.
func newIdentityFromPEM(certificatePEM []byte) *identity.X509Identity {
	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSignFromPEM creates a function that generates a digital signature from a message digest using the uploaded private key.
func newSignFromPEM(privateKeyPEM []byte) identity.Sign {
	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func NewChaincodeService(cert []byte, privateKey []byte) (*ChaincodeService, error) {
	clientConnection := newGrpcConnection()
	// defer clientConnection.Close()

	id := newIdentityFromPEM(cert)
	sign := newSignFromPEM(privateKey)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, err
	}
	// defer gw.Close()

	// Use environment variables for chaincode and channel names
	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	return &ChaincodeService{client: clientConnection, gw: gw, network: network, contract: contract}, nil
}

// func initLedger(contract *client.Contract) {
// 	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

// 	_, err := contract.SubmitTransaction("InitLedger")
// 	if err != nil {
// 		panic(fmt.Errorf("failed to submit transaction: %w", err))
// 	}

// 	fmt.Printf("*** Transaction committed successfully\n")
// }

// Evaluate a transaction to query ledger state.
func (c *ChaincodeService) GetAllAssets() ([]byte, error) {
	fmt.Println("\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

	evaluateResult, err := c.contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return evaluateResult, nil
}

// // Submit a transaction synchronously, blocking until it has been committed to the ledger.
func (c *ChaincodeService) CreateAsset(asset model.Asset) error {
	fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, Color, Size, Owner and AppraisedValue arguments \n")

	// _, err := c.contract.SubmitTransaction("CreateAsset", assetId, "yellow", "5", "Tom", "1300")
	_, err := c.contract.SubmitTransaction("CreateAsset", asset.ID, asset.Color, strconv.Itoa(asset.Size), asset.Owner, strconv.Itoa(asset.AppraisedValue))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	return nil
}

// Evaluate a transaction by assetID to query ledger state.
func (c *ChaincodeService) ReadAssetByID(assetId string) ([]byte, error) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := c.contract.EvaluateTransaction("ReadAsset", assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return evaluateResult, nil
}

// // Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// // this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
// func transferAssetAsync(contract *client.Contract) {
// 	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, updates existing asset owner")

// 	submitResult, commit, err := contract.SubmitAsync("TransferAsset", client.WithArguments(assetId, "Mark"))
// 	if err != nil {
// 		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
// 	}

// 	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Mark. \n", string(submitResult))
// 	fmt.Println("*** Waiting for transaction commit.")

// 	if commitStatus, err := commit.Status(); err != nil {
// 		panic(fmt.Errorf("failed to get commit status: %w", err))
// 	} else if !commitStatus.Successful {
// 		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
// 	}

// 	fmt.Printf("*** Transaction committed successfully\n")
// }
