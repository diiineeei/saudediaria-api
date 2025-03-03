package utils

import (
	"cloud.google.com/go/compute/metadata"
	"context"
)

func AppEngineGetProjectID() string {
	projectID, err := metadata.ProjectIDWithContext(context.Background())
	if err != nil {
		return "rodinei"
	}
	return projectID
}
