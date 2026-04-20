package main

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	app        *firebase.App
	firestore  *firestore.Client
	auth       *auth.Client
	projectID  string
	databaseID string
}

type AppResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type FirestoreDoc struct {
	ID   string                 `json:"id"`
	Path string                 `json:"path"`
	Data map[string]interface{} `json:"data"`
}

func NewClient(ctx context.Context, projectID string, databaseID string, serviceAccountPath string, isEmulator bool, firestoreHost, authHost string) (*Client, error) {
	var opts []option.ClientOption

	if isEmulator {
		fHost := strings.TrimPrefix(strings.TrimPrefix(firestoreHost, "http://"), "https://")
		aHost := strings.TrimPrefix(strings.TrimPrefix(authHost, "http://"), "https://")
		os.Setenv("FIRESTORE_EMULATOR_HOST", fHost)
		os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", aHost)
		opts = append(opts, option.WithoutAuthentication())
	} else {
		os.Unsetenv("FIRESTORE_EMULATOR_HOST")
		os.Unsetenv("FIREBASE_AUTH_EMULATOR_HOST")
		if serviceAccountPath != "" {
			opts = append(opts, option.WithCredentialsFile(serviceAccountPath))
		}
	}

	config := &firebase.Config{ProjectID: projectID}
	if !isEmulator && serviceAccountPath != "" && projectID == "" {
		// If projectID is empty and we have a service account,
		// firebase.NewApp will try to discover it from the credentials file.
		config.ProjectID = ""
	}

	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, err
	}

	if databaseID == "" {
		databaseID = "(default)"
	}

	actualProjectID := projectID
	if actualProjectID == "" && !isEmulator {
		// Try to detect it from the environment/credentials if not provided
		// Note: firebase.App internal projectID is not exported, but we can use firestore.DetectProjectID
		actualProjectID = firestore.DetectProjectID
	}

	firestoreClient, err := firestore.NewClientWithDatabase(ctx, actualProjectID, databaseID, opts...)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		app:        app,
		firestore:  firestoreClient,
		auth:       authClient,
		projectID:  projectID,
		databaseID: databaseID,
	}, nil
}

func (c *Client) Close() {
	if c.firestore != nil {
		c.firestore.Close()
	}
}

func (c *Client) GetCollections(ctx context.Context, parentPath string) ([]string, error) {
	var collections []string
	if parentPath == "" {
		it := c.firestore.Collections(ctx)
		for {
			coll, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			collections = append(collections, coll.ID)
		}
	} else {
		it := c.firestore.Doc(parentPath).Collections(ctx)
		for {
			coll, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			collections = append(collections, coll.ID)
		}
	}
	return collections, nil
}

type QueryCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // "==", ">", "<", ">=", "<=", "array-contains", "in"
	Value    interface{} `json:"value"`
}

type QueryOptions struct {
	Conditions []QueryCondition `json:"conditions"`
	OrderBy    string           `json:"orderBy"`
	Descending bool             `json:"descending"`
	Limit      int              `json:"limit"`
}

func (c *Client) QueryDocuments(ctx context.Context, path string, opts QueryOptions) ([]FirestoreDoc, error) {
	col := c.firestore.Collection(path)
	var query firestore.Query = col.Query

	for _, cond := range opts.Conditions {
		field := cond.Field
		val := cond.Value
		if field == "id" {
			field = firestore.DocumentID
			// In Go Admin SDK, when filtering by DocumentID, the value must be a *DocumentRef
			if idStr, ok := val.(string); ok {
				val = col.Doc(idStr)
			}
		}
		query = query.Where(field, cond.Operator, val)
	}

	if opts.OrderBy != "" {
		field := opts.OrderBy
		if field == "id" {
			field = firestore.DocumentID
		}
		direction := firestore.Asc
		if opts.Descending {
			direction = firestore.Desc
		}
		query = query.OrderBy(field, direction)
	}

	// Always apply limit to prevent massive downloads if options are empty
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	} else {
		query = query.Limit(100)
	}

	it := query.Documents(ctx)
	docs := make([]FirestoreDoc, 0) // Explicitly initialize empty slice instead of nil
	for {
		doc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		docs = append(docs, FirestoreDoc{
			ID:   doc.Ref.ID,
			Path: path + "/" + doc.Ref.ID,
			Data: doc.Data(),
		})
	}
	return docs, nil
}

func (c *Client) UpdateDocument(ctx context.Context, path string, data map[string]interface{}) error {
	_, err := c.firestore.Doc(path).Set(ctx, data, firestore.MergeAll)
	return err
}

func (c *Client) DeleteDocument(ctx context.Context, path string) error {
	_, err := c.firestore.Doc(path).Delete(ctx)
	return err
}

type AuthUser struct {
	UID         string `json:"uid"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoURL"`
	Disabled    bool   `json:"disabled"`
	Created     int64  `json:"created"`
	LastSignIn  int64  `json:"lastSignIn"`
}

func (c *Client) GetUsers(ctx context.Context, max int) ([]AuthUser, error) {
	users := make([]AuthUser, 0)
	it := c.auth.Users(ctx, "")
	for {
		user, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		users = append(users, AuthUser{
			UID:         user.UID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			PhotoURL:    user.PhotoURL,
			Disabled:    user.Disabled,
			Created:     user.UserMetadata.CreationTimestamp,
			LastSignIn:  user.UserMetadata.LastRefreshTimestamp,
		})

		if len(users) >= max {
			break
		}
	}
	return users, nil
}

func (c *Client) DeleteUser(ctx context.Context, uid string) error {
	return c.auth.DeleteUser(ctx, uid)
}
