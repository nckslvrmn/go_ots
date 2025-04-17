package ots_firestore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/nckslvrmn/go_ots/pkg/simple_crypt"
	"github.com/nckslvrmn/go_ots/pkg/utils"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreStore struct {
	client *firestore.Client
}

func NewFirestoreStore() *FirestoreStore {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, utils.FirestoreProjectID, option.WithScopes("https://www.googleapis.com/auth/cloud-platform"))
	if err != nil {
		// In a real application, we would handle this error more gracefully
		panic(fmt.Sprintf("failed to create Firestore client: %v", err))
	}
	return &FirestoreStore{
		client: client,
	}
}

func (f *FirestoreStore) StoreSecret(s *simple_crypt.Secret) error {
	ctx := context.Background()

	// Create document data
	secretData := map[string]interface{}{
		"view_count": s.ViewCount,
		"data":       utils.B64E(s.Data),
		"is_file":    s.IsFile,
		"nonce":      utils.B64E(s.Nonce),
		"salt":       utils.B64E(s.Salt),
		"header":     utils.B64E(s.Header),
		"ttl":        time.Now().AddDate(0, 0, utils.TTLDays).Unix(),
	}

	// Store in Firestore using secret_id as document ID
	_, err := f.client.Collection(utils.FirestoreCollection).Doc(s.SecretId).Set(ctx, secretData)
	if err != nil {
		return fmt.Errorf("failed to store secret in Firestore: %w", err)
	}

	return nil
}

func (f *FirestoreStore) GetSecret(secretId string) (*simple_crypt.Secret, error) {
	ctx := context.Background()

	// Get document from Firestore
	doc, err := f.client.Collection(utils.FirestoreCollection).Doc(secretId).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("secret not found")
		}
		return nil, fmt.Errorf("failed to get secret from Firestore: %w", err)
	}

	// Create secret object
	secret := &simple_crypt.Secret{
		SecretId: secretId,
	}

	// Extract data from document
	data := doc.Data()

	if viewCount, ok := data["view_count"].(int64); ok {
		secret.ViewCount = int(viewCount)
	}

	if isFile, ok := data["is_file"].(bool); ok {
		secret.IsFile = isFile
	}

	if encData, ok := data["data"].(string); ok {
		secret.Data, err = utils.B64D(encData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode data: %w", err)
		}
	}

	if nonce, ok := data["nonce"].(string); ok {
		secret.Nonce, err = utils.B64D(nonce)
		if err != nil {
			return nil, fmt.Errorf("failed to decode nonce: %w", err)
		}
	}

	if salt, ok := data["salt"].(string); ok {
		secret.Salt, err = utils.B64D(salt)
		if err != nil {
			return nil, fmt.Errorf("failed to decode salt: %w", err)
		}
	}

	if header, ok := data["header"].(string); ok {
		secret.Header, err = utils.B64D(header)
		if err != nil {
			return nil, fmt.Errorf("failed to decode header: %w", err)
		}
	}

	return secret, nil
}

func (f *FirestoreStore) DeleteSecret(secretId string) error {
	ctx := context.Background()

	// Delete document from Firestore
	_, err := f.client.Collection(utils.FirestoreCollection).Doc(secretId).Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil // Document already deleted or doesn't exist
		}
		return fmt.Errorf("failed to delete secret from Firestore: %w", err)
	}

	return nil
}

func (f *FirestoreStore) UpdateSecret(s *simple_crypt.Secret) error {
	ctx := context.Background()

	// Update only the view_count field
	_, err := f.client.Collection(utils.FirestoreCollection).Doc(s.SecretId).Update(ctx, []firestore.Update{
		{
			Path:  "view_count",
			Value: s.ViewCount,
		},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("secret not found")
		}
		return fmt.Errorf("failed to update view count for secret: %w", err)
	}

	return nil
}
