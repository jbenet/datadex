package main

import (
	"encoding/json"
	"fmt"
	"github.com/jbenet/data"
	"os"
	"os/exec"
	"strings"
)

const AwsConfigFile = ".awsconfig"

type AwsStsResponse struct {
	Credentials *data.AwsCredentials
}

func getAwsFederationCredentials(user string) (*data.AwsCredentials, error) {

	// ensure AwsConfigFile exists
	if _, err := os.Stat(AwsConfigFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s config file missing.", AwsConfigFile)
	}

	// Federated user acl debugging is annoying. Use get-session-token for now.
	// args := "sts get-federation-token --name " + user
	args := "sts get-session-token --duration-seconds 3600"

	cmd := exec.Command("aws", strings.Split(args, " ")...)

	// make the config file point to the local file.
	cmd.Env = []string{fmt.Sprintf("AWS_CONFIG_FILE=%s", AwsConfigFile)}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("aws cli error: %v", err)
	}

	res := &AwsStsResponse{}
	err = json.Unmarshal(out, res)
	if err != nil {
		return nil, fmt.Errorf("error decoding aws response: %v", err)
	}

	return res.Credentials, nil
}
