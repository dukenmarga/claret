<script>
    import { Plus, X, FolderOpen, Database, Loader2 } from 'lucide-svelte'
    import { AddProject, SelectServiceAccountFile, ParseServiceAccount, ListDatabases } from '../../../wailsjs/go/main/App'
    import { appState } from '../state.svelte.js'
    import { toast } from 'svelte-sonner'

    let { onclose } = $props()
    
    let name = $state('')
    let type = $state('service_account') 
    let serviceAccountPath = $state('')
    let firestoreHost = $state('localhost:8080')
    let authHost = $state('localhost:9099')
    let projectId = $state('')
    
    let availableDatabases = $state([])
    let selectedDatabase = $state('(default)')
    let isFetchingDatabases = $state(false)

    async function handleBrowse() {
        const path = await SelectServiceAccountFile()
        if (path) {
            serviceAccountPath = path
            const resp = await ParseServiceAccount(path)
            if (resp.success) {
                projectId = resp.data
                if (!name) name = resp.data
                
                // Fetch databases for this service account
                isFetchingDatabases = true
                try {
                    const dbResp = await ListDatabases(path)
                    if (dbResp.success) {
                        availableDatabases = dbResp.data || ['(default)']
                        selectedDatabase = availableDatabases[0]
                    }
                } finally {
                    isFetchingDatabases = false
                }
            }
        }
    }

    async function handleAdd() {
        if (!projectId) {
            toast.error("Project ID is required")
            return
        }

        const project = {
            id: `project-${Date.now()}`,
            firebaseProjectId: projectId,
            name: name || projectId,
            type,
            serviceAccountPath: type === 'service_account' ? serviceAccountPath : '',
            isEmulator: type === 'emulator',
            emulatorConfig: type === 'emulator' ? { firestoreHost, authHost } : null,
            color: '#FFCA28',
            databaseId: selectedDatabase
        }

        const resp = await AddProject(project)
        if (resp.success) {
            toast.success('Project added successfully')
            await appState.refreshConfig()
            onclose()
        } else {
            toast.error(`Failed to add project: ${resp.error}`)
        }
    }
</script>

<div class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
    <div class="bg-slate-900 border border-slate-800 rounded-xl shadow-2xl w-full max-w-md overflow-hidden">
        <div class="p-6 border-b border-slate-800 flex items-center justify-between">
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
                <Plus size={20} class="text-firebase-yellow" />
                Add Firebase Project
            </h2>
            <button onclick={onclose} class="text-slate-500 hover:text-white transition-colors">
                <X size={20} />
            </button>
        </div>

        <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">
            <div class="flex p-1 bg-slate-950 rounded-lg border border-slate-800">
                <button onclick={() => type = 'service_account'} class="flex-1 py-1.5 rounded-md text-xs font-medium transition-all {type === 'service_account' ? 'bg-slate-800 text-white' : 'text-slate-500 hover:text-slate-300'}">Service Account</button>
                <button onclick={() => type = 'emulator'} class="flex-1 py-1.5 rounded-md text-xs font-medium transition-all {type === 'emulator' ? 'bg-slate-800 text-white' : 'text-slate-500 hover:text-slate-300'}">Emulator</button>
            </div>

            {#if type === 'service_account'}
                <div>
                    <label class="block text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-1.5">Service Account JSON</label>
                    <div class="flex gap-2">
                        <input bind:value={serviceAccountPath} type="text" placeholder="/path/to/key.json" class="flex-1 bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-1 focus:ring-indigo-500" />
                        <button onclick={handleBrowse} class="bg-slate-800 p-2 rounded-lg text-slate-400 hover:text-white"><FolderOpen size={18} /></button>
                    </div>
                </div>

                {#if projectId}
                    <div class="animate-in fade-in slide-in-from-top-2 duration-300 space-y-4">
                        <div>
                            <label class="block text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-1.5">Project ID (Detected)</label>
                            <div class="bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-slate-400 font-mono">{projectId}</div>
                        </div>

                        <div>
                            <label class="block text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-1.5">Display Name</label>
                            <input bind:value={name} type="text" placeholder="e.g. Production" class="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-1 focus:ring-indigo-500" />
                        </div>

                        <div>
                            <label class="block text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-1.5">Select Database</label>
                            {#if isFetchingDatabases}
                                <div class="flex items-center gap-2 text-slate-500 text-xs py-2"><Loader2 size={14} class="animate-spin" /> Fetching databases...</div>
                            {:else}
                                <div class="grid grid-cols-1 gap-2">
                                    {#each availableDatabases as db}
                                        <button 
                                            onclick={() => selectedDatabase = db}
                                            class="flex items-center justify-between px-3 py-2 rounded-lg border transition-all text-sm {selectedDatabase === db ? 'bg-indigo-500/10 border-indigo-500 text-indigo-400' : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'}"
                                        >
                                            <div class="flex items-center gap-2"><Database size={14} /> {db}</div>
                                            {#if selectedDatabase === db}<div class="w-2 h-2 rounded-full bg-indigo-500 shadow-[0_0_8px_rgba(99,102,241,0.6)]"></div>{/if}
                                        </button>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    </div>
                {/if}
            {:else}
                <div class="space-y-4">
                    <div>
                        <label class="block text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-1.5">Project ID</label>
                        <input bind:value={projectId} type="text" placeholder="demo-project" class="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-1 focus:ring-indigo-500" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-1.5">Firestore Host</label>
                            <input bind:value={firestoreHost} type="text" class="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-1 focus:ring-indigo-500" />
                        </div>
                        <div>
                            <label class="block text-[10px] font-bold text-slate-500 uppercase tracking-widest mb-1.5">Auth Host</label>
                            <input bind:value={authHost} type="text" class="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:ring-1 focus:ring-indigo-500" />
                        </div>
                    </div>
                </div>
            {/if}
        </div>

        <div class="p-6 bg-slate-900/50 border-t border-slate-800 flex gap-3">
            <button onclick={onclose} class="flex-1 py-2 rounded-lg text-sm font-medium text-slate-400 hover:bg-slate-800 transition-colors">Cancel</button>
            <button onclick={handleAdd} class="flex-1 py-2 rounded-lg text-sm font-medium bg-firebase-orange hover:bg-orange-600 text-white shadow-lg shadow-firebase-orange/20 transition-all active:scale-95">Add Project</button>
        </div>
    </div>
</div>
