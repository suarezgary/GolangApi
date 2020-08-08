package cloudfiles

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/suarezgary/GolangApi/config"
	"google.golang.org/api/option"
)

//UploadFile uploadFile uploads an object.
func UploadFile(file []byte, object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(config.Cfg().StorageKeyLocation))
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(config.Cfg().StorageBucket).Object(object).NewWriter(ctx)

	if _, err = wc.Write(file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	url := "https://storage.googleapis.com/" + config.Cfg().StorageBucket + "/" + object
	return url, nil
}

// DeleteFile removes specified object.
func DeleteFile(object string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithAPIKey(config.Cfg().StorageKeyLocation))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(config.Cfg().StorageBucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	return nil
}
