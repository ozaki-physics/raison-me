package share

import (
	"context"
	"fmt"
	"strconv"

	secretManager "cloud.google.com/go/secretmanager/apiv1"
	secretManagerPB "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const (
	GCPProjectID = "raison-me"
)

func GcpSecretValue(name string, version int) (string, error) {
	ctx := context.Background()
	client, err := secretManager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	fullName := "projects/" + GCPProjectID + "/secrets/" + name + "/versions/" + strconv.Itoa(version)
	req := &secretManagerPB.AccessSecretVersionRequest{
		Name: fullName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
