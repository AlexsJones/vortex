package secrets

import (
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudkms/v1"
)

func GoogleKMSFetch(project, location, key string) (string, error) {
	client, err := google.DefaultClient(context.Background(), cloudkms.CloudPlatformScope)
	if err != nil {
		return "" ,err
	}
	_ = client

	return "", nil
}