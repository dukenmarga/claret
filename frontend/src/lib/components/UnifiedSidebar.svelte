<script>
    import { appState } from '../state.svelte.js'
    import { ChevronRight, ChevronDown, Database, Users, Plus, RefreshCw, Folder, Settings, ShieldCheck } from 'lucide-svelte'

    let { onAddProject } = $props()

    let groupedProjects = $derived.by(() => {
        const groups = {}
        appState.projects.forEach(p => {
            if (!groups[p.firebaseProjectId]) {
                groups[p.firebaseProjectId] = {
                    id: p.firebaseProjectId,
                    name: p.name,
                    color: p.color,
                    databases: []
                }
            }
            groups[p.firebaseProjectId].databases.push(p)
        })
        return Object.values(groups)
    })

    function getProjectDisplay(name) {
        return name ? name.substring(0, 2).toUpperCase() : '??'
    }
</script>

<div class="flex flex-col h-full bg-slate-900/80 border-r border-slate-800 shadow-2xl">
    <div class="p-4 border-b border-slate-800 flex items-center justify-between h-14 bg-slate-950/20 shrink-0">
        <div class="flex items-center gap-2">
            <div class="w-7 h-7 rounded-lg bg-firebase-orange flex items-center justify-center text-white font-black text-sm shadow-lg ring-1 ring-white/10">C</div>
            <h1 class="font-black text-xs tracking-[0.2em] text-white uppercase opacity-90">Claret</h1>
        </div>
        <button onclick={onAddProject} class="p-1.5 rounded-lg hover:bg-slate-800 text-slate-300 hover:text-white transition-all bg-slate-900/50 border border-slate-800">
            <Plus size={16} strokeWidth={3} />
        </button>
    </div>

    <div class="flex-1 overflow-y-auto p-3 space-y-2 custom-scrollbar">
        {#each groupedProjects as group}
            <div class="flex flex-col">
                <!-- Project Node -->
                <div 
                    class="group flex items-center gap-2 px-2.5 py-2 rounded-xl transition-all cursor-pointer {appState.currentProject?.firebaseProjectId === group.id ? 'bg-white/5 ring-1 ring-white/10' : 'hover:bg-white/5'}"
                    onclick={() => appState.toggleProjectExpansion(group.id)}
                >
                    <div class="p-0.5 rounded transition-colors {appState.expandedProjectIds.includes(group.id) ? 'text-white' : 'text-slate-500'}">
                        {#if appState.expandedProjectIds.includes(group.id)}<ChevronDown size={16} strokeWidth={2.5} />{:else}<ChevronRight size={16} strokeWidth={2.5} />{/if}
                    </div>
                    <div class="w-6 h-6 rounded-lg flex items-center justify-center text-[10px] font-black text-white shadow-md ring-1 ring-white/20" style="background-color: {group.color || '#FFCA28'}">
                        {getProjectDisplay(group.name)}
                    </div>
                    <span class="flex-1 text-xs font-black truncate tracking-tight {appState.currentProject?.firebaseProjectId === group.id ? 'text-white' : 'text-slate-300 group-hover:text-white'}">{group.name}</span>
                </div>

                {#if appState.expandedProjectIds.includes(group.id)}
                    <div class="ml-4 pl-4 border-l-2 border-slate-800/50 mt-1 space-y-1 animate-in fade-in slide-in-from-left-2 duration-200">
                        {#each group.databases as db}
                            <div class="flex flex-col">
                                <button 
                                    onclick={() => { appState.selectProject(db.id, db.databaseId); appState.activeTab = 'firestore' }}
                                    class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg text-[11px] transition-all {appState.currentProjectId === db.id && appState.activeTab === 'firestore' ? 'text-white bg-indigo-500/30 font-black shadow-lg' : 'text-slate-400 hover:text-white hover:bg-white/5'}"
                                >
                                    <Database size={14} class={appState.currentProjectId === db.id && appState.activeTab === 'firestore' ? 'text-white' : 'text-slate-500'} />
                                    <span>Firestore <span class="{appState.currentProjectId === db.id && appState.activeTab === 'firestore' ? 'text-white/60' : 'text-slate-600'} font-mono text-[10px] ml-1">[{db.databaseId || 'default'}]</span></span>
                                </button>

                                {#if appState.currentProjectId === db.id && appState.activeTab === 'firestore'}
                                    <div class="ml-4 pl-3 border-l-2 border-indigo-500/30 mt-1 space-y-0.5 mb-2">
                                        {#each appState.collections as coll}
                                            <button 
                                                onclick={() => appState.selectCollection(coll)}
                                                class="w-full flex items-center gap-2 px-2.5 py-1.5 rounded-md transition-all {appState.selectedPath === coll ? 'text-white bg-white/10 font-black shadow-sm' : 'text-slate-400 hover:text-white hover:bg-white/5'}"
                                            >
                                                <Folder size={12} class={appState.selectedPath === coll ? 'text-white' : 'text-slate-600'} />
                                                <span class="truncate text-xs">{coll}</span>
                                            </button>
                                        {/each}
                                    </div>
                                {/if}
                            </div>
                        {/each}

                        <button 
                            onclick={() => { appState.selectProject(group.databases[0].id, group.databases[0].databaseId); appState.activeTab = 'auth' }}
                            class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-lg text-[11px] transition-all {appState.currentProject?.firebaseProjectId === group.id && appState.activeTab === 'auth' ? 'text-white bg-indigo-500/30 font-black shadow-lg' : 'text-slate-400 hover:text-white hover:bg-white/5'}"
                        >
                            <ShieldCheck size={14} class={appState.currentProject?.firebaseProjectId === group.id && appState.activeTab === 'auth' ? 'text-white' : 'text-slate-500'} />
                            <span>Authentication</span>
                        </button>
                    </div>
                {/if}
            </div>
        {/each}
    </div>

    <div class="p-4 border-t border-slate-800 bg-slate-950/40 flex items-center justify-between shrink-0">
        <button class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-all">
            <Settings size={18} />
        </button>
        <span class="text-[8px] font-mono tracking-tighter text-white opacity-20 uppercase">Claret v1.4.0</span>
    </div>
</div>

<style>
    .custom-scrollbar::-webkit-scrollbar { width: 4px; }
    .custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.05); border-radius: 10px; }
</style>
