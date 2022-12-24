package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

const (
	SORACOM_CLI_BIN  = "soracom"
	AUTH_KEY_ID_FLAG = "--auth-key-id"
	AUTH_KEY_FLAG    = "--auth-key"
	BODY_FLAG        = "--body"
)

const (
	AUTH_KEY_ID_ENV_KEY         = "SORACOM_AUTH_KEY_ID"
	AUTH_KEY_ENV_KEY            = "SORACOM_AUTH_KEY"
	AUTH_KEY_SECRET_ARN_ENV_KEY = "SORACOM_AUTH_KEY_SECRET_ARN"
)

type SoracomCliEvent struct {
	Command string      `json:"command"`
	Body    interface{} `json:"body"`
}

type SoracomAuthKey struct {
	AuthKeyId string `json:"AUTH_KEY_ID"`
	AuthKey   string `json:"AUTH_KEY"`
}

var (
	secretCache, _ = secretcache.New()
)

func HandleRequest(ctx context.Context, event SoracomCliEvent) (interface{}, error) {
	fmt.Printf("event: %+v\n", event)
	out, err := soracom(event.Command, event.Body)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = json.Unmarshal([]byte(out), &result)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON: %s", out)
	}
	return result, err
}

func soracom(command string, body interface{}) (string, error) {
	arg, err := soracomcli(command, body)
	if err != nil {
		return "", fmt.Errorf("could not build command: %w", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", arg)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if stderr.String() != "" {
		fmt.Printf("stderr: %s", stderr.String())
	}
	if err != nil {
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func authKey() (SoracomAuthKey, error) {
	if os.Getenv(AUTH_KEY_ID_ENV_KEY) != "" {
		return SoracomAuthKey{AuthKeyId: os.Getenv(AUTH_KEY_ID_ENV_KEY), AuthKey: os.Getenv(AUTH_KEY_ENV_KEY)}, nil
	}

	secretstr, err := secretCache.GetSecretString(os.Getenv(AUTH_KEY_SECRET_ARN_ENV_KEY))
	if err != nil {
		return SoracomAuthKey{}, err
	}

	var authKey SoracomAuthKey
	err = json.Unmarshal([]byte(secretstr), &authKey)
	if err != nil {
		return SoracomAuthKey{}, err
	}
	return authKey, nil
}

func soracomcli(command string, body interface{}) (string, error) {
	key, err := authKey()
	if err != nil {
		return "", err
	}
	cmd := SORACOM_CLI_BIN + " " +
		AUTH_KEY_ID_FLAG + " " + key.AuthKeyId + " " +
		AUTH_KEY_FLAG + " " + key.AuthKey + " " +
		command
	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		cmd = cmd + " " + BODY_FLAG + " '" + string(bytes) + "'"
	}
	return cmd, nil
}

func main() {
	lambda.Start(HandleRequest)
}
