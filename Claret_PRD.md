# Project Title & Overview: **Claret**

**Claret** is a lightweight, high-performance open-source desktop GUI (Graphical User Interface) application for managing the Firebase ecosystem (Firestore, Auth, Storage). Built as a modern and efficient alternative to "Firefoo," Claret focuses on data access speed, full support for the Firebase Local Emulator Suite, and local credential security.

**Branding & Rationale:**
1.  **Claret**: A name that sounds sophisticated and sharp, representing a refined tool for data management.
2.  **Primary Branding**:
    - **Primary Colors**: `#FFCA28` (Firebase Yellow) & `#F57C00` (Firebase Orange).
    - **Background**: Slate-950 (Dark-first design).
    - **Accent**: Indigo-500 or Cyan-400 for developer interactions.

---

## User Personas & Problem Statements

### User Personas
1.  **The Solo Indie Hacker**: Builds apps quickly, needs Firestore data access without having to open the heavy and slow Google web console.
2.  **The Local Emulator Addict**: A developer who works 90% in offline/emulator environments and needs a GUI that can switch between `localhost` and `production` with a single click.
3.  **The QA/Support Engineer**: Needs to perform minor data edits or view Auth user logs without touching the source code.

### Problem Statements
- Firebase Console (Web) is often slow, heavy, and consumes significant browser resources (RAM).
- Firefoo (Existing tool) is paid and can feel bloated for simple needs.
- Managing data in the Emulator Suite via the default Firebase UI is very limited (no advanced filtering, export, or mass delete).

---

## User Stories (Gherkin Style)

### Firestore Data Exploration
- **Given** I have a Service Account JSON file.
- **When** I upload that file to Claret.
- **Then** I can see the list of collections in the sidebar and explore documents in table or JSON format.

### Local Emulator Connection
- **Given** the Firebase Emulator is running on port 8080.
- **When** I select the "Connect to Local Emulator" option in Claret.
- **Then** Claret must automatically detect the project ID and display data from the emulator without requiring an internet connection.

### Inline Data Editing
- **Given** I am viewing a row of data in the table.
- **When** I double-click on a field.
- **Then** I can directly change its value and press `Enter` to save the changes to Firestore.

### Document Duplication and Migration
- **Given** I have selected one or more documents in a collection.
- **When** I select the "Copy to Collection" action and provide a target collection path.
- **Then** Claret must create exact copies of those documents in the destination collection, maintaining the field types.

---

## Technical Stack & Architecture

### Backend (Go + Wails)
- **Runtime**: Wails v2/v3.
- **Firebase SDK**: `firebase.google.com/go/v4` (Official Admin SDK).
- **Local Storage**: `badgerdb` or a simple JSON file in `AppData`/`Application Support` to store project configurations.
- **Key Features**: Go-routines for fetching large datasets without blocking the UI.

### Frontend (Svelte 5 + Vite)
- **Framework**: Svelte 5 (Runes: `$state`, `$derived`, `$effect`).
- **State Management**: Svelte Native Stores / Context API.
- **Table**: `@tanstack/svelte-table` v8.
- **Code Editor**: Monaco Editor (via `monaco-editor` npm).
- **Styling**: Tailwind CSS v4.
- **Icons**: Lucide Svelte.

### Architecture Diagram (Conceptual)
```text
[ Svelte UI ] <--> [ Wails Bindings ] <--> [ Go Logic ] <--> [ Firebase Admin SDK ]
      ^                  ^                      |                   |
      |                  |                      v                   v
[ Local Config ] <-------+-------------- [ Service Account ]   [ Google APIs / Local Emulator ]
```

---

## Functional Requirements (Detailed)

### Phase 1: MVP (Must-Have)
- **Project Manager**:
    - Add/Remove Firebase Projects via Service Account JSON path.
    - Toggle between `Live` (Production) and `Emulator` mode.
- **Firestore Browser & Manager**:
    - Hierarchical Sidebar (Collection -> Document -> Subcollection).
    - Advanced Data Table:
        - Server-side pagination (limit/offset via Firestore `StartAfter`).
        - Client-side quick filter.
        - Multi-row selection for bulk actions.
    - JSON Inspector: Side panel to view full document JSON.
    - **Data Operations**:
        - **Insert**: Dialog to add a new document (Auto-ID/Manual-ID) with field type validation.
        - **Edit**: Inline table cell editing and full JSON editing in the side panel.
        - **Delete**: Delete a single document or multiple selected documents with a confirmation prompt.
        - **Copy/Move**: Copy one or more selected documents from the current collection to a user-specified target collection.
- **Auth Manager**:
    - List Users (UID, Email, Created At, Last Sign In).
    - Delete User & Reset Password link.
- **Export**: Export current view to JSON/CSV.

### Phase 2: Power Tools (Should-Have)
- **Query Builder**: UI for `where`, `orderBy`, and `limit`.
- **Storage Browser**: List buckets, list files, upload/download (drag & drop).
- **Command Palette**: `Ctrl+K` to switch collections quickly.

---

## Data Models (Internal Schema)

Claret stores configuration data locally (e.g., `~/.claret/config.json`):

### Config Store (JSON)
```json
{
  "settings": {
    "theme": "dark",
    "fontSize": 14,
    "autoConnectEmulator": true
  },
  "projects": [
    {
      "id": "project-uuid-1",
      "name": "My Production App",
      "type": "service_account",
      "serviceAccountPath": "/users/path/to/key.json",
      "isEmulator": false,
      "color": "#ef4444"
    },
    {
      "id": "project-uuid-2",
      "name": "Local Dev",
      "type": "emulator",
      "emulatorConfig": {
        "firestoreHost": "localhost:8080",
        "authHost": "localhost:9099"
      },
      "isEmulator": true,
      "color": "#10b981"
    }
  ]
}
```

### Bridge Types (Go to TS)
```go
// Result Wrapper
type AppResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// Firestore Document Map
type FirestoreDoc struct {
    ID   string                 `json:"id"`
    Data map[string]interface{} `json:"data"`
}
```

---

## API & Integration Points (Wails Bindings)

| Method Name | Parameters | Description |
| :--- | :--- | :--- |
| `ConnectProject(id string)` | `projectId` | Initializes the Firebase Admin SDK for a specific project. |
| `GetCollections()` | - | Retrieves the list of root collections. |
| `GetDocuments(path string, limit int, startAfter string)` | `path`, `limit`, `cursor` | Fetches document data within a collection. |
| `UpdateDocument(path string, data string)` | `path`, `jsonData` | Updates fields within a document. |
| `DeleteDocuments(paths []string)` | `paths` | Deletes one or more documents by their full Firestore paths. |
| `CopyDocuments(sourcePaths []string, targetCol string)` | `sourcePaths`, `targetCol` | Copies selected documents into a new collection. |
| `GetUsers(max int)` | `max` | Fetches the user list from Firebase Auth. |
| `CheckEmulatorStatus()` | - | Pings default emulator ports (8080, 9099). |

---

## UI/UX Flow & Component Descriptions

### 1. Main Layout Structure
- **Left Sidebar (Narrow)**: Project Switcher icons, Settings, Help.
- **Secondary Sidebar (Wide)**: Tree View for Collection/Sub-collection.
- **Main Content Area**: 
    - **Top Bar**: Breadcrumbs (e.g., `users / u123 / posts`), Query Filter bar, Refresh button, Action Buttons (Add, Bulk Delete, Copy).
    - **Central Grid**: TanStack Table with checkboxes for multi-selection.
- **Right Inspector Panel (Collapsible)**: Monaco Editor for editing raw JSON of the selected row.

### 2. Interaction Details
- **Filtering**: Text input above table columns performs instant client-side filtering.
- **Inline Edit**: Click on a table cell -> change value -> press `Enter` -> Go Backend sends update via SDK.
- **Multi-Select**: Use checkboxes on the left-most column to select multiple documents. A floating action bar appears to provide "Delete Selected" and "Copy to..." options.
- **Copy Path**: Right-click on a field in the JSON viewer -> "Copy Firestore Path".

### 3. Component List (shadcn-svelte)
- `Table`: For the primary data display with selection row logic.
- `Resizable`: To divide the Sidebar, Main, and Inspector areas.
- `Dialog/Modal`: For "Add Document", "Copy Documents", or "Add Project" forms.
- `Command`: For the Command Palette feature (Ctrl+K).
- `Sonner`: For toast notifications (Error/Success).

---

## Potential Technical Challenges & Solutions
- **Memory Management**: Loading 10,000+ documents at once.
    - *Solution*: Implement virtual scrolling in TanStack Table and a default limit of 50 docs per fetch.
- **Deep Nested Objects**: Firestore supports maps within maps.
    - *Solution*: Use a JSON editor (Monaco) for complex map editing; use dot-notation for flat table visualization.
- **Auth Latency**: The Auth user list can be slow if there are millions of users.
    - *Solution*: Always use pagination and provide clear loading feedback.

## Roadmap
1.  **Week 1**: Setup Wails + Svelte 5 + Tailwind 4. Implement Project Manager.
2.  **Week 2**: Firestore Tree View & Table View. Implementation of multi-select logic.
3.  **Week 3**: CRUD Firestore: Insert, Edit, Bulk Delete, and Copy to Collection features.
4.  **Week 4**: Emulator support & Auth View.
5.  **Week 5**: Export features & UI Polishing (Dark Mode).

---
**Claret** - *Blazing fast Firebase management, right on your desktop.*