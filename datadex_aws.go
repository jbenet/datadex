package main

import (
	"encoding/json"
	"fmt"
	"github.com/jbenet/data"
	"os"
	"os/exec"
	"strings"
)

type AwsStsResponse struct {
	Credentials *data.AwsCredentials
}

func getAwsFederationCredentials(user string) (*data.AwsCredentials, error) {
	ak := os.Getenv("S3_ACCESS_KEY")
	sk := os.Getenv("S3_SECRET_KEY")
	if len(ak) < 1 || len(sk) < 1 {
		return nil, fmt.Errorf("aws credentials not provided.")
	}

	// Federated user acl debugging is annoying. Use get-session-token for now.
	// args := "sts get-federation-token --name " + user
	args := "sts get-session-token --duration-seconds 3600"

	cmd := exec.Command("aws", strings.Split(args, " ")...)
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
