package utils

import (
	"fmt"
	"github.com/oleiade/reflections"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"testing"

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

var configFilePath = "../../kubernetes/helm/functional-testing/values/test.yaml"

type node struct {
	ID   string
	HOST string
}


type testnet struct {
	ContractAddresses struct{
		PaymentObligation string `yaml:"paymentObligation"`
	}`yaml:"contractAddresses"`
}

type config struct {

	Nodes string `yaml:"nodes"`
	Network string `yaml:"network"`

	Namespace string `yaml:"namespace"`

	Rinkeby testnet `yaml:"rinkeby"`
	Kovan testnet `yaml:"kovan"`
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
		Network = "testing"
	}

}

func GetConfig() config{

	var c config
	yamlFile, err := ioutil.ReadFile(configFilePath)
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}

	n, err := reflections.GetField(c, "Rinkeby")

	fmt.Println(n)



	a, ok := n.(testnet)
	if ok {
		fmt.Println("paymentObligation")
		fmt.Println(a.ContractAddresses.PaymentObligation)
	}

	return c
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
