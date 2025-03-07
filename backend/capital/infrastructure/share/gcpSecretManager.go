package share

import (
	"context"
	"fmt"
	"strconv"

	secretManager "cloud.google.com/go/secretmanager/apiv1"
	secretManagerPB "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/ozaki-physics/raison-me/share/config"
)

func GcpSecretValue(name string, version int) (string, error) {
	ctx := context.Background()
	client, err := secretManager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	globalConfig := config.NewConfig()
	fullName := "projects/" + globalConfig.GetGCPProjectID() + "/secrets/" + name + "/versions/" + strconv.Itoa(version)
	req := &secretManagerPB.AccessSecretVersionRequest{
		Name: fullName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
