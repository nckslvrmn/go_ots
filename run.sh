#!/bin/bash

GCP_PROJECT_ID=flawless-snow-457215-k2 \
  FIRESTORE_DATABASE=nckslvrmn-goots \
  GCS_BUCKET=nckslvrmn-goots \
  GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/goots.json \
  go run main.go
