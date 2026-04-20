<script>
  import { appState } from './lib/state.svelte.js'
  import { onMount } from 'svelte'
  import { LayoutGrid, Database, Users, Settings as SettingsIcon, Plus, Trash2, Copy, RefreshCw, ChevronRight, Loader2, Save } from 'lucide-svelte'
  import DataTable from './lib/components/DataTable.svelte'
  import TreeTable from './lib/components/TreeTable.svelte'
  import AddProjectDialog from './lib/components/AddProjectDialog.svelte'
  import JsonEditor from './lib/components/JsonEditor.svelte'
  import AuthManager from './lib/components/AuthManager.svelte'
  import DeleteConfirmDialog from './lib/components/DeleteConfirmDialog.svelte'
  import { Toaster, toast } from 'svelte-sonner'
  
  let activeTab = $state('firestore') 
  let showAddProject = $state(false)
  let editedJson = $state('')

  $effect(() => {
    if (appState.selectedDoc) {
      editedJson = JSON.stringify(appState.selectedDoc.data, null, 2)
    } else {
      editedJson = ''
    }
  })

  async function handleSave() {
    if (!appState.selectedDoc) return
    try {
      const data = JSON.parse(editedJson)
      await appState.saveDocument({ ...appState.selectedDoc, data })
    } catch (e) {
      toast.error("Invalid JSON: " + e.message)
    }
  }
</script>

<div class="flex h-screen w-screen bg-slate-950 text-slate-200 overflow-hidden font-sans">
  <Toaster theme="dark" position="bottom-right" />
  
  <!-- Left Sidebar -->
  <aside class="w-16 border-r border-slate-800 flex flex-col items-center py-4 gap-4 bg-slate-900/50 shrink-0">
    <div class="w-10 h-10 rounded-xl bg-firebase-orange flex items-center justify-center text-white font-bold text-xl shadow-lg mb-4">C</div>
    {#each appState.projects as project}
      <button class="w-10 h-10 rounded-lg flex items-center justify-center transition-all hover:scale-110 {appState.currentProjectId === project.id ? 'ring-2 ring-firebase-yellow bg-slate-800' : 'bg-slate-800/50 hover:bg-slate-800'}" onclick={() => appState.selectProject(project.id)}>
        <div class="flex flex-col items-center">
            <span class="text-[10px] font-bold text-slate-400 mb-0.5 uppercase">{project.name.substring(0, 2)}</span>
            <span class="w-2 h-2 rounded-full" style="background-color: {project.color}"></span>
        </div>
      </button>
    {/each}
    <button onclick={() => showAddProject = true} class="w-10 h-10 rounded-lg bg-slate-800/30 flex items-center justify-center hover:bg-slate-800 transition-colors mt-2 text-slate-400 border border-dashed border-slate-700"><Plus size={20} /></button>
  </aside>

  <!-- Explorer Sidebar -->
  <aside class="w-64 border-r border-slate-800 flex flex-col bg-slate-900/30 shrink-0">
    <div class="p-4 border-b border-slate-800 flex items-center justify-between bg-slate-900/50 h-14">
      <h2 class="font-semibold text-sm uppercase tracking-wider text-slate-400">Explorer</h2>
      <button class="text-slate-500 hover:text-white" onclick={() => appState.refreshCollections('')}><RefreshCw size={14} class={appState.isLoading ? 'animate-spin' : ''} /></button>
    </div>
    <nav class="flex-1 overflow-y-auto p-2">
      {#if appState.currentProject}
        <div class="mb-4 px-2 py-2 bg-slate-800/40 rounded-lg border border-slate-700/50">
            <div class="flex items-center gap-2 mb-1">
                <div class="w-2 h-2 rounded-full" style="background-color: {appState.currentProject.color}"></div>
                <span class="text-xs font-bold text-white truncate">{appState.currentProject.name}</span>
            </div>
            <div class="text-[10px] text-slate-500 font-mono truncate">ID: {appState.currentProject.firebaseProjectId}</div>
            <div class="text-[10px] text-indigo-400 font-mono truncate mt-1">DB: {appState.currentProject.databaseId || '(default)'}</div>
        </div>
      {/if}
      <div class="space-y-1">
        <button class="w-full flex items-center gap-2 px-3 py-2 rounded-md text-sm transition-colors {activeTab === 'firestore' ? 'bg-indigo-500/10 text-indigo-400' : 'hover:bg-slate-800'}" onclick={() => activeTab = 'firestore'}><Database size={16} /> Firestore</button>
        <button class="w-full flex items-center gap-2 px-3 py-2 rounded-md text-sm transition-colors {activeTab === 'auth' ? 'bg-indigo-500/10 text-indigo-400' : 'hover:bg-slate-800'}" onclick={() => activeTab = 'auth'}><Users size={16} /> Auth</button>
      </div>
      {#if activeTab === 'firestore' && appState.currentProjectId}
        <div class="mt-6"><h3 class="px-3 text-xs font-semibold text-slate-500 uppercase mb-2">Collections</h3>
          <div class="space-y-0.5">
            {#each appState.collections as collection}
              <button class="w-full flex items-center gap-2 px-3 py-1.5 rounded-md text-sm transition-colors {appState.selectedPath === collection ? 'bg-indigo-500/10 text-indigo-400' : 'hover:bg-slate-800/50 text-slate-300'}" onclick={() => appState.selectCollection(collection)}>
                <ChevronRight size={14} class={appState.selectedPath === collection ? 'text-indigo-400' : 'text-slate-600'} />{collection}
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </nav>
  </aside>

  <!-- Main Content Area -->
  <main class="flex-1 flex flex-col min-w-0 bg-slate-950">
    <header class="h-14 border-b border-slate-800 flex items-center justify-between px-6 bg-slate-900/20 shrink-0">
      <div class="flex items-center gap-2 text-sm overflow-hidden flex-1 mr-4">
        <span class="text-slate-500">claret</span><ChevronRight size={14} class="text-slate-700" />
        <span class="text-slate-300 font-medium truncate">{appState.currentProject?.name || 'No Project'}</span>
        {#if appState.isLoading}<Loader2 size={14} class="animate-spin text-indigo-500 ml-2" />{/if}
        {#if appState.selectedPath}<ChevronRight size={14} class="text-slate-700" /><span class="text-indigo-400 truncate">{appState.selectedPath}</span>{/if}
      </div>
      <div class="flex items-center gap-3">
        {#if activeTab === 'firestore'}<button onclick={() => appState.addDocument(appState.selectedPath)} disabled={!appState.selectedPath} class="bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white px-3 py-1.5 rounded-md text-xs font-medium flex items-center gap-2 transition-all active:scale-95"><Plus size={14} /> Add Document</button>{/if}
      </div>
    </header>

    <div class="flex-1 overflow-hidden relative p-4">
      {#if activeTab === 'firestore'}
        {#if appState.selectedPath}
            {#if appState.viewMode === 'table'}
                <DataTable path={appState.selectedPath} />
            {:else}
                <TreeTable path={appState.selectedPath} />
            {/if}
        {:else}
            <div class="h-full flex flex-col items-center justify-center text-slate-600"><LayoutGrid size={48} class="mb-4 opacity-10" /><p>Select a collection to browse</p></div>
        {/if}
      {:else}
        <AuthManager />
      {/if}
      {#if appState.error}<div class="absolute bottom-4 left-6 right-6 bg-red-500/10 border border-red-500/20 text-red-400 p-3 rounded-lg text-sm flex justify-between z-30 shadow-xl"><span>{appState.error}</span><button onclick={() => appState.error = null}>✕</button></div>{/if}
    </div>
  </main>

  <!-- Inspector Side Panel -->
  {#if activeTab === 'firestore'}
  <aside class="w-[600px] border-l border-slate-800 bg-slate-900/50 flex flex-col shrink-0">
    <div class="p-4 border-b border-slate-800 flex items-center justify-between h-14 bg-slate-900/50">
      <h2 class="font-semibold text-sm uppercase tracking-wider text-slate-400">Inspector</h2>
      {#if appState.selectedDoc}<button onclick={handleSave} class="bg-indigo-600 hover:bg-indigo-500 text-white px-3 py-1 rounded text-xs font-medium transition-colors flex items-center gap-2 active:scale-95"><Save size={14} /> Save</button>{/if}
    </div>
    <div class="flex-1 overflow-hidden bg-slate-950">
        {#if appState.selectedDoc}<JsonEditor value={editedJson} onchange={(val) => editedJson = val} />{:else}<div class="h-full flex flex-col items-center justify-center text-slate-600 gap-4"><p>Select a document to inspect</p></div>{/if}
    </div>
  </aside>
  {/if}
</div>

{#if showAddProject}<AddProjectDialog onclose={() => showAddProject = false} />{/if}
{#if appState.showDeleteDialog}<DeleteConfirmDialog />{/if}
