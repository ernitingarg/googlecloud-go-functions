package env

import (
	"fmt"
	"os"
	"soteria-functions/entity"
)

type envvars struct {
	GCP *gcpVars
	BTC *btcVars
	ETH *ethVars
}

type gcpVars struct {
	ProjectID     ProjectID
	KeyPath       string
	KmsLocationID LocationID
}

type btcVars struct {
	ChainID entity.BtcChainID
}

type ethVars struct {
	ChainID entity.EthChainID
}

type ProjectID string

func (p ProjectID) String() string {
	return string(p)
}

const (
	DEVELOP    ProjectID = "black-stream-292507"
	PRODUCTION ProjectID = "soteria-production"
)

type LocationID string

func (l LocationID) String() string {
	return string(l)
}

const (
	GLOBAL LocationID = "global"
)

var EnvVars *envvars

func init() {
	newEnvVars()
}

func newEnvVars() {
	// GCP vars
	gcpProject := os.Getenv("GCP_PROJECT")
	keyPath := os.Getenv("KEY_FILE_PATH")
	kmsLocationID := os.Getenv("KMS_LOCATION_ID")

	fmt.Printf("gcpProject: %v\n", gcpProject)
	fmt.Printf("keyPath: %v\n", keyPath)
	fmt.Printf("kmsLocationID: %v\n", kmsLocationID)

	var projectID ProjectID
	switch gcpProject {
	case "black-stream-292507":
		projectID = DEVELOP
	case "soteria-production":
		projectID = PRODUCTION
	default:
		panic("project id is invalid")
	}

	var locationID LocationID
	switch kmsLocationID {
	case "global":
		locationID = GLOBAL
	default:
		panic("kms location id is invalid")
	}

	// BTC vars
	bcid := os.Getenv("BTC_CHAIN_ID")
	var btcChainID entity.BtcChainID
	switch bcid {
	case "main":
		btcChainID = entity.MainBtc
	case "test":
		btcChainID = entity.TestBtc
	default:
		panic("btc chain id is invalid")
	}

	// ETH vars
	ethid := os.Getenv("ETH_CHAIN_ID")
	var ethChainID entity.EthChainID
	switch ethid {
	case "main":
		ethChainID = entity.MainEth
	case "test":
		ethChainID = entity.TestEth
	default:
		panic("eth chain id is invalid")
	}

	gcpVars := &gcpVars{
		ProjectID:     projectID,
		KeyPath:       keyPath,
		KmsLocationID: locationID,
	}
	btcVars := &btcVars{
		ChainID: btcChainID,
	}
	ethVars := &ethVars{
		ChainID: ethChainID,
	}
	EnvVars = &envvars{
		GCP: gcpVars,
		BTC: btcVars,
		ETH: ethVars,
	}
}
