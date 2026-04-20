package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/api/firestore/v1"
	"google.golang.org/api/option"
)

type App struct {
	ctx            context.Context
	configManager  *Manager
	firebaseClient *Client
}

func NewApp() *App {
	cm, _ := NewManager()
	return &App{
		configManager: cm,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type ServiceAccountInfo struct {
	ProjectID string `json:"project_id"`
}

func (a *App) ParseServiceAccount(path string) AppResponse {
	data, err := os.ReadFile(path)
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}

	var info ServiceAccountInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return AppResponse{Success: false, Error: "Invalid Service Account JSON"}
	}

	return AppResponse{Success: true, Data: info.ProjectID}
}

func (a *App) SelectServiceAccountFile() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Firebase Service Account JSON",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JSON Files (*.json)",
				Pattern:     "*.json",
			},
		},
	})
	if err != nil {
		return ""
	}
	return path
}

func (a *App) GetConfig() Config {
	return a.configManager.GetConfig()
}

func (a *App) AddProject(p Project) AppResponse {
	err := a.configManager.AddProject(p)
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}
	return AppResponse{Success: true, Data: a.configManager.GetConfig().Projects}
}

func (a *App) RemoveProject(id string) AppResponse {
	err := a.configManager.RemoveProject(id)
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}
	return AppResponse{Success: true, Data: a.configManager.GetConfig().Projects}
}

func (a *App) ConnectProject(id string, databaseID string) AppResponse {
	conf := a.configManager.GetConfig()
	var targetProject *Project
	for _, p := range conf.Projects {
		if p.ID == id {
			targetProject = &p
			break
		}
	}

	if targetProject == nil {
		return AppResponse{Success: false, Error: "Project not found"}
	}

	if a.firebaseClient != nil {
		a.firebaseClient.Close()
	}

	var firestoreHost, authHost string
	if targetProject.EmulatorConfig != nil {
		firestoreHost = targetProject.EmulatorConfig.FirestoreHost
		authHost = targetProject.EmulatorConfig.AuthHost
	}

	// Use provided databaseID if available, else project default, else "(default)"
	dbID := databaseID
	if dbID == "" {
		dbID = targetProject.DatabaseID
	}

	client, err := NewClient(
		a.ctx,
		targetProject.FirebaseProjectID,
		dbID,
		targetProject.ServiceAccountPath,
		targetProject.IsEmulator,
		firestoreHost,
		authHost,
	)

	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}

	a.firebaseClient = client
	return AppResponse{Success: true}
}

func (a *App) GetCollections(parentPath string) AppResponse {
	if a.firebaseClient == nil {
		return AppResponse{Success: false, Error: "Not connected to any project"}
	}
	cols, err := a.firebaseClient.GetCollections(a.ctx, parentPath)
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}
	return AppResponse{Success: true, Data: cols}
}

func (a *App) GetDocuments(path string, opts QueryOptions) AppResponse {
	if a.firebaseClient == nil {
		return AppResponse{Success: false, Error: "Not connected to any project"}
	}

	// Add server-side logging for debug
	log.Printf("Querying path: %s", path)
	log.Printf("With options: %+v", opts)

	docs, err := a.firebaseClient.QueryDocuments(a.ctx, path, opts)
	if err != nil {
		log.Printf("Query error: %v", err)
		return AppResponse{Success: false, Error: err.Error()}
	}

	log.Printf("Query returned %d documents", len(docs))
	return AppResponse{Success: true, Data: docs}
}

func (a *App) UpdateDocument(path string, data map[string]interface{}) AppResponse {
	if a.firebaseClient == nil {
		return AppResponse{Success: false, Error: "Not connected to any project"}
	}
	err := a.firebaseClient.UpdateDocument(a.ctx, path, data)
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}
	return AppResponse{Success: true}
}

func (a *App) DeleteDocuments(paths []string) AppResponse {
	if a.firebaseClient == nil {
		return AppResponse{Success: false, Error: "Not connected to any project"}
	}
	for _, p := range paths {
		err := a.firebaseClient.DeleteDocument(a.ctx, p)
		if err != nil {
			return AppResponse{Success: false, Error: err.Error()}
		}
	}
	return AppResponse{Success: true}
}

func (a *App) GetUsers(max int) AppResponse {
	if a.firebaseClient == nil {
		return AppResponse{Success: false, Error: "Not connected to any project"}
	}

	users, err := a.firebaseClient.GetUsers(a.ctx, max)
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}

	// We need to return a simplified version or ensure the data is serializable
	// Sometimes the ExportedUserRecord contains complex types that Wails/JSON might struggle with
	return AppResponse{Success: true, Data: users}
}

func (a *App) DeleteUser(uid string) AppResponse {
	if a.firebaseClient == nil {
		return AppResponse{Success: false, Error: "Not connected to any project"}
	}
	err := a.firebaseClient.DeleteUser(a.ctx, uid)
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}
	return AppResponse{Success: true}
}

func (a *App) ListDatabases(serviceAccountPath string) AppResponse {
	ctx := context.Background()

	svc, err := firestore.NewService(ctx, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return AppResponse{Success: false, Error: err.Error()}
	}

	infoResp := a.ParseServiceAccount(serviceAccountPath)
	if !infoResp.Success {
		return infoResp
	}
	projectID := infoResp.Data.(string)

	parent := fmt.Sprintf("projects/%s", projectID)
	resp, err := svc.Projects.Databases.List(parent).Do()
	if err != nil {
		log.Printf("List databases failed: %v", err)
		return AppResponse{Success: true, Data: []string{"(default)"}}
	}

	var names []string
	for _, db := range resp.Databases {
		parts := strings.Split(db.Name, "/")
		if len(parts) > 0 {
			names = append(names, parts[len(parts)-1])
		}
	}

	if len(names) == 0 {
		names = append(names, "(default)")
	}

	return AppResponse{Success: true, Data: names}
}
