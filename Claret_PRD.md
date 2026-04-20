# Project Title & Overview: **Claret**

**Claret** is a lightweight, high-performance open-source desktop GUI application for managing the Firebase ecosystem (Firestore, Auth, Storage). Built with **Wails v2**, **Go 1.25**, and **Svelte 5**, Claret provides a blazing-fast, high-contrast developer experience for managing both production and emulator environments.

**Branding & Rationale:**
1.  **Claret**: A name representing a refined, sophisticated tool for data management.
2.  **Primary Branding**:
    - **Primary Colors**: `#FFCA28` (Firebase Yellow) & `#F57C00` (Firebase Orange).
    - **Background**: Slate-950 (Dark-first design).
    - **Accent**: Indigo-500 for primary interactions, Emerald-500 for data values.
    - **Sidebar**: Ultra-high contrast (White text on Slate) for professional IDE feel.

---

## User Stories (Completed)

### Firestore Data Exploration
- **Given** I have a Service Account JSON file.
- **When** I browse the file in the Add Project dialog.
- **Then** Claret automatically detects my Project ID and lists all available Firestore Databases for selection.

### Hybrid Tree-Table Browsing
- **Given** I am in a collection.
- **When** I toggle between "Table" and "Tree" modes.
- **Then** I see a synchronized view where data is represented as a professional column-based tree (Key/Value/Type) or a high-density grid.

### Safe Data Deletion
- **Given** I need to delete one or more documents.
- **When** I trigger the delete action.
- **Then** I must type "DELETE" in a security dialog that lists the specific Document IDs and paths being removed.

### Advanced Authentication Management
- **Given** I am in the Auth Manager.
- **When** I view the user list.
- **Then** I can edit user metadata (Display Name, Photo URL) and toggle their account status (Enabled/Disabled) directly.

---

## Technical Stack (Finalized)

### Backend (Go + Wails)
- **Runtime**: Wails v2.
- **Go Version**: 1.25 (using 1.23.0 compatibility mode in `go.mod`).
- **Firebase SDK**: `firebase.google.com/go/v4`.
- **Key Backend Logic**: Custom serialization to ensure empty lists return `[]` instead of `null` to the frontend.

### Frontend (Svelte 5 + Tailwind 4)
- **Framework**: Svelte 5 (100% Runes: `$state`, `$derived`, `$effect`, `$bindable`).
- **Styling**: Tailwind CSS v4 (Zero-config Vite plugin).
- **JSON Editor**: `svelte-jsoneditor` (Replacing Monaco for better Wails compatibility and tree/text hybrid support).
- **Notifications**: `svelte-sonner`.
- **Icons**: Lucide Svelte.

---

## Functional Requirements (Detailed)

### 1. Project & Sidebar Management (Completed)
- **Unified Sidebar**: Grouped hierarchy (Project > Firestore Instance > Collections).
- **Database Discovery**: Auto-scans Service Accounts to list multiple Firestore databases.
- **Database Detach**: Hover-action to remove a database from local config without deleting server data.

### 2. Firestore Manager (Completed)
- **Hybrid View Modes**:
    - **Table View**: Compact grid with inline cell editing and infinite scroll ("Load More").
    - **Tree View**: Hierarchical columns (Key/Value/Type) with recursive expansion and inline editing.
- **Server-Side Query Builder**: 
    - Multi-condition support (`==`, `!=`, `>`, etc.).
    - Native `DocumentID` support for `__key__` queries.
    - Server-side `OrderBy` and sorting via column header clicks.
- **Data Operations**:
    - **Renaming**: Background Document ID rename (Clone + Delete logic).
    - **Duplication**: Copy selected documents to any target collection path.
    - **Bulk Actions**: Multi-select documents for mass delete or mass duplicate.

### 3. Authentication Manager (Completed)
- **User List**: Displays UID, Email, Status (Active/Disabled), and Metadata.
- **User Editor**: Dedicated dialog to update profile info and account status.
- **Security**: Full user deletion with confirmation.

---

## UI/UX Design Standards
- **High Contrast**: Primary navigation uses White (`#FFFFFF`) bold text for visibility.
- **Reactive UI**: All changes (filtering, sorting, path changes) are instantaneous with zero screen flickering.
- **Density**: Maximum data density in tables/trees, utilizing 11px-12px font sizes and minimal padding.
- **Safety**: Red-themed dialogs for destructive actions; blue/indigo for navigation.

---
**Claret v1.4.0** - *The professional desktop standard for Firebase data management.*
