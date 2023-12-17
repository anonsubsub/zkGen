package zkpolicy

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Protection struct {
	Algorithm string `json:"algorithm"`
	Parameter string `json:"parameter"`
}

type PrivateArgument struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Protection Protection `json:"protection"`
}

type PublicArgument struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type Relation struct {
	PrivateArgument PrivateArgument `json:"private_argument"`
	PublicArgument  PublicArgument  `json:"public_argument"`
}

type Constraint struct {
	String string `json:"string"`
}

type ZkPolicy struct {
	Name        string     `json:"name"`
	Relations   []Relation `json:"relations"`
	Constraints []string   `json:"constraints"`
}

func ParseZkPolicy(policyName string) (ZkPolicy, error) {

	// init zkPolicy struct
	var zkPolicy ZkPolicy

	// read config file
	jsonFile, err := os.Open("zkpolicy/" + policyName + ".json")
	if err != nil {
		log.Println("os.Open() error", err)
		return zkPolicy, err
	}
	defer jsonFile.Close()

	// parse json
	byteValue, _ := io.ReadAll(io.Reader(jsonFile))
	json.Unmarshal(byteValue, &zkPolicy)

	return zkPolicy, nil
}
