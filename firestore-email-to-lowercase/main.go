package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize Firestore client
	ctx := context.Background()
	credentials := map[string]string{
		"type":                        "service_account",
		"project_id":                  os.Getenv("GCLOUD_PROJECT_ID"),
		"private_key_id":              os.Getenv("SERVICE_ACCOUNT_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("SERVICE_ACCOUNT_PRIVATE_KEY"),
		"client_email":                os.Getenv("SERVICE_ACCOUNT_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("SERVICE_ACCOUNT_CLIENT_ID"),
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/" + os.Getenv("SERVICE_ACCOUNT_CLIENT_EMAIL"),
	}

	credJSON, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("Failed to marshal credentials: %v", err)
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCLOUD_PROJECT_ID"), option.WithCredentialsJSON(credJSON))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	collectionName := "users" // Change this to your desired collection name
	const batchSize = 500
	totalProcessed := 0
	totalUpdated := 0

	bulkWriter := client.BulkWriter(ctx)
	var lastDocID string

	for {
		log.Printf("Starting batch processing. Last processed document ID: %s\n", lastDocID)

		query := client.Collection(collectionName).OrderBy(firestore.DocumentID, firestore.Asc)
		if lastDocID != "" {
			query = query.StartAfter(lastDocID).Limit(batchSize)
		} else {
			query = query.Limit(batchSize)
		}

		iter := query.Documents(ctx)
		iterationComplete := true

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				log.Println("No more documents in current iteration.")
				break
			}
			if err != nil {
				log.Fatalf("Error iterating documents: %v", err)
			}

			iterationComplete = false

			data := doc.Data()

			email, exists := data["email"]
			if !exists {
				totalProcessed++
				lastDocID = doc.Ref.ID
				continue
			}

			emailStr, ok := email.(string)
			if !ok {
				totalProcessed++
				lastDocID = doc.Ref.ID
				continue
			}

			trimmedEmail := strings.TrimSpace(emailStr)
			lowerEmail := strings.ToLower(trimmedEmail)

			if lowerEmail != emailStr {
				log.Printf("UPDATE NEEDED - Document ID: %s | Original: %s | Normalized: %s\n",
					doc.Ref.ID, emailStr, lowerEmail)

				bulkWriter.Update(doc.Ref, []firestore.Update{
					{Path: "email", Value: lowerEmail},
				})
				totalUpdated++
			}

			totalProcessed++
			lastDocID = doc.Ref.ID
		}

		// Always continue to next iteration if no fatal errors
		if iterationComplete {
			log.Println("All iterations complete. Exiting processing.")
			break
		}

		log.Printf("Batch processing summary - Processed: %d, Updated: %d\n", totalProcessed, totalUpdated)
	}

	bulkWriter.Flush()
	log.Printf("FINAL PROCESSING SUMMARY - Total documents processed: %d, Total documents updated: %d\n",
		totalProcessed, totalUpdated)
}
