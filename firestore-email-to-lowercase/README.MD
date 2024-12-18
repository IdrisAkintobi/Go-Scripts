# Firestore Bulk Email Normalization Script

A Go script to normalize `email` fields in a Firestore collection by trimming whitespace and converting to lowercase.

## Key Features

- Targets the `users` collection by default (modifiable via `collectionName`).
- Processes documents in batches using Firestore's BulkWriter.
- Updates only when normalization is required.

## Setup

1. Add the following to a `.env` file:
   ```env
   GCLOUD_PROJECT_ID=your-gcp-project-id
   SERVICE_ACCOUNT_PRIVATE_KEY_ID=your-private-key-id
   SERVICE_ACCOUNT_PRIVATE_KEY="your-private-key"
   SERVICE_ACCOUNT_CLIENT_EMAIL=your-service-account-email
   SERVICE_ACCOUNT_CLIENT_ID=your-service-account-id
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the script:
   ```bash
   go run main.go
   ```

## Notes

- Change `collectionName` to target a different Firestore collection.
- Requires Go 1.19+, Firestore enabled in your Google Cloud project, and valid service account credentials.

## License

Available under the MIT License.
