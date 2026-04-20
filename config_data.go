package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Project struct {
	ID                 string          `json:"id"`                // Internal Claret ID
	FirebaseProjectID  string          `json:"firebaseProjectId"` // Actual Firebase Project ID
	Name               string          `json:"name"`
	Type               string          `json:"type"` // "service_account" or "emulator"
	ServiceAccountPath string          `json:"serviceAccountPath,omitempty"`
	IsEmulator         bool            `json:"isEmulator"`
	EmulatorConfig     *EmulatorConfig `json:"emulatorConfig,omitempty"`
	Color              string          `json:"color"`
	DatabaseID         string          `json:"databaseId,omitempty"` // "(default)" or other
}

type EmulatorConfig struct {
	FirestoreHost string `json:"firestoreHost"`
	AuthHost      string `json:"authHost"`
}

type Settings struct {
	Theme               string `json:"theme"`
	FontSize            int    `json:"fontSize"`
	AutoConnectEmulator bool   `json:"autoConnectEmulator"`
}

type Config struct {
	Settings Settings  `json:"settings"`
	Projects []Project `json:"projects"`
}

type Manager struct {
	config Config
	path   string
	mu     sync.RWMutex
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(home, ".claret")
	configPath := filepath.Join(configDir, "config.json")

	m := &Manager{
		path: configPath,
		config: Config{
			Settings: Settings{
				Theme:               "dark",
				FontSize:            14,
				AutoConnectEmulator: true,
			},
			Projects: []Project{},
		},
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, err
		}
		if err := m.Save(); err != nil {
			return nil, err
		}
	} else {
		if err := m.Load(); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.config)
}

func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

func (m *Manager) GetConfig() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

func (m *Manager) AddProject(p Project) error {
	m.mu.Lock()
	m.config.Projects = append(m.config.Projects, p)
	m.mu.Unlock()
	return m.Save()
}

func (m *Manager) RemoveProject(id string) error {
	m.mu.Lock()
	newProjects := []Project{}
	for _, p := range m.config.Projects {
		if p.ID != id {
			newProjects = append(newProjects, p)
		}
	}
	m.config.Projects = newProjects
	m.mu.Unlock()
	return m.Save()
}

func (m *Manager) UpdateSettings(s Settings) error {
	m.mu.Lock()
	m.config.Settings = s
	m.mu.Unlock()
	return m.Save()
}
