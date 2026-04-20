# Project Title & Overview: **Claret**

**Claret** is a high-performance, open-source desktop GUI application for managing the Firebase ecosystem (Firestore, Auth, Storage). Built with **Wails v2**, **Go 1.25**, and **Svelte 5**, Claret provides a blazing-fast, ultra-high-contrast developer experience for managing both production and emulator environments.

**Branding & Rationale:**
1.  **Claret**: A name representing a refined, sophisticated tool for data management.
2.  **Primary Branding**:
    - **Primary Colors**: `#FFCA28` (Firebase Yellow) & `#F57C00` (Firebase Orange).
    - **Background**: Slate-950 (Dark-first design).
    - **Accent**: Indigo-500 for primary interactions, Emerald-500 for data values.
    - **Sidebar**: High-contrast "Ultra" design (Brilliant White text on Deep Slate) with vertical guide lines.

---

## User Stories (Completed)

### High-Speed Project Navigation
- **Given** I manage multiple Firebase projects and databases.
- **When** I look at the sidebar.
- **Then** I see a unified tree grouping Firestore instances and Auth by project, with active selections highlighted in brilliant white.

### Real-Time Hierarchy Exploration
- **Given** I am browsing deeply nested Firestore collections.
- **When** I expand a document row or switch to Tree mode.
- **Then** I see synchronized, hierarchical data (Key/Value/Type) and sub-collections directly in the middle explorer.

### Advanced Authentication Management
- **Given** I am in the Auth Manager.
- **When** I search for a user by Email, UID, or Phone.
- **Then** Claret performs an optimized scan and displays users with detailed "Providers" icons and high-precision "Signed In" timestamps.

### Direct User Administration
- **Given** a user account needs a metadata update or suspension.
- **When** I click the edit icon in the Auth table.
- **Then** I can instantly update their Display Name, Avatar URL, or toggle their account status.

---

## Technical Stack (Finalized)

### Backend (Go + Wails)
- **Runtime**: Wails v2.
- **Go Version**: 1.25 (using 1.23.0 compatibility mode in `go.mod`).
- **Firebase SDK**: `firebase.google.com/go/v4`.
- **Search Engine**: Custom Go-based partial match scanner with safety capping (100 users) for performance.

### Frontend (Svelte 5 + Tailwind 4)
- **Framework**: Svelte 5 (Utilizing pure Class-based state with `$state` and reactive getters).
- **Styling**: Tailwind CSS v4 (Zero-config).
- **JSON Editor**: `svelte-jsoneditor` (Tree & Text modes).
- **Notifications**: `svelte-sonner`.
- **Iconography**: `lucide-svelte` + custom branding SVGs (e.g., Google G).

---

## Functional Requirements (Detailed)

### 1. Unified Sidebar & Project Engine (Completed)
- **Project Grouping**: All databases and services (Auth/Firestore) are nested under a single Project root.
- **High-Contrast Design**: Active nodes use Bold Brilliant White (`#FFFFFF`) with indigo ring highlights.
- **Database Discovery**: Auto-scans Service Accounts for all available Firestore databases during import.
- **Database Management**: Detach databases from the local view via hover-action.

### 2. Firestore Management (Completed)
- **Hybrid Modes**: Toggle between high-density Table mode and hierarchical Tree mode with synchronized filters.
- **Recursive Navigation**: Deep-drill into documents and sub-collections without changing views.
- **Inline Editing**: Double-click support for cell value updates and Document ID renaming (Clone/Delete).
- **Advanced Querying**: Server-side Query Builder with multi-condition logic and sort persistence.

### 3. Authentication Manager (Completed)
- **Unified Identifier**: Shows Phone, Email, or Anonymous status at a glance.
- **Provider Tracking**: Displays authentication source icons (Google, Password, Phone, Anonymous).
- **Full Metadata CRUD**: Edit Display Name, Photo URL, and Enable/Disable status.
- **Performance**: Optimized table rendering using CSS containment and fixed layout for zero-lag hover effects.

---

## UI/UX Design Standards
- **Ultra-Contrast**: Primary text is Brilliant White; data values use Emerald/Indigo accents.
- **High Density**: 11px-12px font sizes with minimal padding for professional "IDE-style" data density.
- **Safety Protocol**: Mandatory "DELETE" confirmation for destructive actions.
- **Instant Feedback**: Toast notifications for all server-side write operations.

---
**Claret v1.5.0-PRO** - *The professional standard for Firebase desktop administration.*
