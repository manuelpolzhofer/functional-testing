package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/oleiade/reflections"
	"gopkg.in/yaml.v2"

	"github.com/gavv/httpexpect"
)

const (
	NODE1         = "node1"
	NODE2         = "node2"
	INVOICE       = "invoice"
	PURCHASEORDER = "purchaseorder"
)

var Nodes map[string]node
var Network string
var Testnet string

var configFilePath = "../../kubernetes/helm/functional-testing/values/test.yaml"
var testNet *testnet

type node struct {
	ID   string
	HOST string
}

type testnet struct {
	ContractAddresses struct {
		PaymentObligation string `yaml:"paymentObligation"`
	} `yaml:"contractAddresses"`
}

type network struct {
	Nodes    string `yaml:"nodes"`
	Testnets struct {
		Rinkeby testnet `yaml:"rinkeby"`
		Kovan   testnet `yaml:"kovan"`
	} `yaml:"testnets"`
}

type config struct {
	Nodes   string `yaml:"nodes"`
	Network string `yaml:"network"`

	Namespace string `yaml:"namespace"`

	Networks struct {
		Russianhill   network `yaml:"russianhill"`
		Bernalheights network `yaml:"bernalheights"`
	}
}

func SetupEnvironment() {
	nodesEnv := os.Getenv("NODES")
	idsEnv := os.Getenv("IDS")
	nodesSlice := SplitString(nodesEnv)
	idsSlice := SplitString(idsEnv)

	if len(nodesSlice) == 0 {
		nodesSlice = append(nodesSlice, "https://localhost:8082", "https://localhost:8083")
	}

	if len(idsSlice) == 0 {
		idsSlice = append(idsSlice, "0x8c8cfaf732d3", "0x24fe6555beb9")
	}

	Nodes = map[string]node{
		NODE1: {
			idsSlice[0],
			nodesSlice[0],
		},
		NODE2: {
			idsSlice[1],
			nodesSlice[1],
		},
	}

	Network = os.Getenv("NETWORK")
	if Network == "" {
		Network = "russianhill"
	}

	Testnet = os.Getenv("TESTNET")
	if Testnet == "" {
		Testnet = "rinkeby"
	}

	testConfig := readConfig()
	var err error
	testNetwork, err := getNetwork(testConfig, Network)

	if err != nil {
		panic(err)
	}

	testNet, err = getTestNet(testNetwork, Testnet)

}

func readConfig() config {

	var c config
	yamlFile, err := ioutil.ReadFile(configFilePath)
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}

	return c

}

func getTestNet(network *network, testnetName string) (*testnet, error) {

	t, err := reflections.GetField(network.Testnets, strings.Title(testnetName))
	if err != nil {
		return nil, err
	}

	testNet, ok := t.(testnet)
	if ok {
		return &testNet, nil
	}

	return nil, fmt.Errorf("could not parse testnet name")

}

func getNetwork(config config, networkName string) (*network, error) {

	testNetwork, err := reflections.GetField(config.Networks, strings.Title(networkName))
	if err != nil {
		return nil, err
	}

	t, ok := testNetwork.(network)
	if ok {
		return &t, nil
	}

	return nil, fmt.Errorf("could not parse network name")

}

func GetPaymentObigationAddress() string {

	return testNet.ContractAddresses.PaymentObligation
}

func GetInsecureClient(t *testing.T, nodeId string) *httpexpect.Expect {
	SetupEnvironment()
	return CreateInsecureClient(t, Nodes[nodeId].HOST)
}

func SplitString(data string) []string {
	result := strings.Split(data, ",")
	if result[0] == "" {
		return []string{}
	}

	return result
}
