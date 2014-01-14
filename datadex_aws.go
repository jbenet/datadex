package main

import (
	"encoding/json"
	"fmt"
	"github.com/jbenet/data"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const AwsConfigFile = ".awsconfig"

type AwsStsResponse struct {
	Credentials *data.AwsCredentials
}

func init() {
	// ensure AwsConfigFile exists
	if _, err := os.Stat(AwsConfigFile); os.IsNotExist(err) {
		pErr("%s file missing. See README.md\n", AwsConfigFile)
		os.Exit(-1)
	}

	// ensure AwsConfigFile has been modified
	f, err := os.Open(AwsConfigFile)
	if err != nil {
		pErr("%v\n", err)
		pErr("Error opening %s file.\n", AwsConfigFile)
		os.Exit(-1)
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		pErr("%v\n", err)
		pErr("Error reading %s file.\n", AwsConfigFile)
		os.Exit(-1)
	}

	s := string(buf[:])
	if strings.Contains(s, "aws_access_key_id = YOUR_ACCESS_KEY") ||
		strings.Contains(s, "aws_secret_access_key = YOUR_SECRET_KEY") {
		pErr(AwsConfigFileNotModifiedErr, AwsConfigFile)
		os.Exit(-1)
	}
}

func getAwsFederationCredentials(user string) (*data.AwsCredentials, error) {

	// ensure AwsConfigFile exists
	if _, err := os.Stat(AwsConfigFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s file missing.", AwsConfigFile)
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

const AwsConfigFileNotModifiedErr = `
Looks like you have not modified the %s file.
Datadex requires an S3 bucket, and credentials to access it.
Please update the file, replacing 'YOUR_ACCESS_KEY' and
'YOUR_SECRET_KEY' with valid aws keys for the bucket user.

Note that end users of data must also be able to upload to this
bucket; datadex federates access to it (using session tokens
generated with 'aws sts get-session-token').
`
