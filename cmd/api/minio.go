package main

import (
	"log"

	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/samuelmolero26/go-backend-course/internal/env"
	"fmt"
)


func newMinioClient() (*minio.Client, error) {
	endpoint := env.GetString("MINIO_ENDPOINT", "localhost:9000")
	accessKey := env.GetString("MINIO_ACCESS_KEY", "minioadmin")
	secretKey := env.GetString("MINIO_SECRET_KEY", "minioadmin")


	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,	
	})

	if err != nil {
		return nil, err
	}

	log.Printf("MinIO client connected to %s", endpoint)
	return client, nil

}

func setPublicReadPolicy(ctx context.Context, client *minio.Client, bucket string) error {
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucket)

	return client.SetBucketPolicy(ctx, bucket, policy)
}

func ensureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Bucket %s already exists", bucket)
		return nil
	}
	if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		return err
	}

	log.Printf("Bucket %s created", bucket)
	return nil
}