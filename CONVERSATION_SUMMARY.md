# Claret: Technical Development Summary & Architecture Reference

**Project Goal**: A high-performance, cross-platform desktop alternative to the Firebase Web Console, optimized for developer productivity and large dataset management.

## 1. Tech Stack Overview
- **Backend**: Go 1.25 (running with 1.23 compatibility in `go.mod` to bypass Wails CLI parser bugs).
- **Frontend**: Svelte 5 (utilizing Runes: `$state`, `$derived`, `$effect`), Tailwind CSS v4.
- **Bridge**: Wails v2 (v2.10.1+).
- **Persistence**: Local configuration stored in `~/.claret/config.json`.
- **Core Libraries**:
    - `firebase.google.com/go/v4`: Firebase Admin SDK.
    - `svelte-jsoneditor`: Advanced JSON tree/text editing.
    - `svelte-sonner`: Toast notifications.
    - `lucide-svelte`: Icon set.

---

## 2. Architecture & State Management
### Backend (Go)
- **`app.go`**: Acts as the main bridge. Handles OS-level operations (File Dialogs) and orchestrates the `firebase_client`.
- **`firebase_client.go`**: A wrapper for the Admin SDK. 
    - Implements **Server-Side Querying** with dynamic type conversion.
    - Handles **Emulator vs. Production** environment toggles via environment variables (`FIRESTORE_EMULATOR_HOST`).
    - **Crucial Fix**: Explicitly initializes empty slices (`make([]T, 0)`) before JSON serialization to ensure the frontend receives `[]` instead of `null`.

### Frontend (Svelte 5)
- **Centralized Store (`state.svelte.js`)**: Uses a singleton `AppState` class.
    - **Methods**: Implemented as arrow functions to ensure proper `this` binding through Svelte's reactive proxy.
    - **Reactivity**: Uses `$state` for data and `get` syntax for derived properties (like `filteredDocuments`) to ensure 100% reliable UI updates.
- **DataTable Component**: A high-density data grid supporting:
    - **Inline Tree Navigation**: Row-level expansion to browse sub-collections.
    - **Direct Editing**: Double-click cells to update Firestore.
    - **Bulk Actions**: Multi-select documents for deletion or duplication.

---

## 3. Implemented Features & Use-Cases

### Use-Case 1: Rapid Data Correction
- **Story**: A developer finds a typo in a production document field across multiple records.
- **Workflow**: 
    1. Select records via checkboxes.
    2. Use **Local Filter** to find specific values.
    3. **Double-click** the cell to correct.
    4. Hit **Enter**. Data persists immediately.

### Use-Case 2: Deep Hierarchy Exploration
- **Story**: Managing a "Chat" application where messages are in `channels/{id}/messages/{id}`.
- **Workflow**:
    1. Navigate to `channels`.
    2. Click the **Chevron (>)** on a channel document.
    3. Click the nested `messages` collection name.
    4. The table instantly drills down into the nested messages.

### Use-Case 3: Data Migration/Backup
- **Story**: Duplicating production configuration to a "staging" collection for testing.
- **Workflow**:
    1. Select all documents in `configs`.
    2. Click **Duplicate to...**.
    3. Type `configs_staging`.
    4. Claret performs a background batch copy.

### Use-Case 4: Advanced Search & Debugging
- **Story**: Finding all "Admin" users who haven't logged in since 2023.
- **Workflow**:
    1. Open **Server Query**.
    2. Add Condition: `isAdmin == true`.
    3. Add Condition: `lastLogin < 1704067200`.
    4. Click **Run Query**.

---

## 4. Technical Hurdles & Resolved Issues

### Wails CLI vs. Go 1.25
- **Problem**: Internal error `package "encoding/json" without types` during binding generation.
- **Solution**: Set `go 1.23.0` in `go.mod`. Use `go build -tags dev,debug` to bypass the Wails CLI parser when the binding generation bug persists.

### OpenSUSE WebKitGTK Versioning
- **Problem**: OpenSUSE provides `webkit2gtk-4.1.pc` while Wails expects `4.0`.
- **Solution**: 
    1. Use the `webkit2gtk41` build tag.
    2. Permanent Fix: Alias the `.pc` files via symlink:
       `ln -s /usr/lib64/pkgconfig/webkit2gtk-4.1.pc ~/.local/lib/pkgconfig/webkit2gtk-4.0.pc`.

### Firestore "Missing Field" Sorting
- **Problem**: Sorting by a field hides documents that don't have that field.
- **Solution**: Implemented a "Behavioral Hint" notification to inform users why document counts might change during sorting.

---

## 5. Deployment & Runtime
### Development Command
```bash
go run -tags dev,debug .
```
### Production Build
```bash
go build -tags production -o build/bin/claret .
```

---

## 6. Future Expansion Ideas
- **Export/Import**: JSON/CSV export of queried results.
- **Auth Custom Claims**: UI for managing user permissions directly in the Auth Manager.
- **Real-time Listener**: Optional "Watch" mode using Firestore Snapshots.
- **Schema Enforcement**: Saving collection-level schema definitions locally to provide better input validation during inline editing.
