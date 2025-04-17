package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
)

var S3Bucket string
var TTLDays int
var DynamoTable string
var AWSRegion string
var FirestoreProjectID string
var FirestoreCollection string
var GCSBucket string
var UsesAWS bool
var UsesGCP bool

func LoadEnv() error {
	// required ENV vars
	S3Bucket = os.Getenv("S3_BUCKET")
	DynamoTable = os.Getenv("DYNAMO_TABLE")
	FirestoreProjectID = os.Getenv("FIRESTORE_PROJECT_ID")
	FirestoreCollection = os.Getenv("FIRESTORE_COLLECTION")
	GCSBucket = os.Getenv("GCS_BUCKET")

	// Check if either AWS or Google Cloud configuration is provided
	UsesAWS = S3Bucket != "" && DynamoTable != ""
	UsesGCP = FirestoreProjectID != "" && FirestoreCollection != "" && GCSBucket != ""
	if !UsesAWS && !UsesGCP {
		return fmt.Errorf("missing required ENV vars - must provide either AWS (S3_BUCKET + DYNAMO_TABLE) or Google Cloud (FIRESTORE_PROJECT_ID + FIRESTORE_COLLECTION + GCS_BUCKET) configuration")
	}

	// optional ENV vars
	AWSRegion = "us-east-1"
	if regionValue, regionExists := os.LookupEnv("AWS_REGION"); regionExists {
		AWSRegion = regionValue
	}

	TTLDays = 7
	if ttlValue, ttlExists := os.LookupEnv("TTL_DAYS"); ttlExists {
		TTLDays, _ = strconv.Atoi(ttlValue)
	}
	return nil
}

func SanitizeViewCount(view_count string) int {
	vc, err := strconv.ParseFloat(view_count, 64)
	if err != nil {
		return 1
	}
	vc = math.Abs(vc)

	if vc > 0 && vc < 10 {
		return int(vc)
	}
	return 1
}

func RandString(length int, url_safe bool) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if url_safe == false {
		chars = chars + "!#$%&*+-=?@_~"
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[num.Int64()]
	}
	return string(result)
}

func RandBytes(length int) []byte {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return randomBytes
}

func B64E(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func B64D(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}
