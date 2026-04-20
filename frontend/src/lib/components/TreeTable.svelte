<script>
    import { appState } from '../state.svelte.js'
    import TreeTableRow from './TreeTableRow.svelte'
    import { ListTree, LayoutGrid, RefreshCw, Filter, Search, X, Database, Plus, Play } from 'lucide-svelte'

    let { path = '' } = $props()
    
    let collectionName = $derived(path.split('/').pop() || 'Root')

    // Copy columns logic from DataTable for the filter dropdown
    let columns = $derived.by(() => {
        const fields = new Set(['id'])
        if (appState.documents && appState.documents.length > 0) {
            appState.documents.forEach(doc => {
                if (doc?.data) Object.keys(doc.data).forEach(key => fields.add(key))
            })
        }
        return Array.from(fields).sort()
    })

    let showQueryBuilder = $state(false)
    const operators = ['==', '!=', '>', '<', '>=', '<=', 'array-contains', 'in']
</script>

<div class="flex flex-col h-full bg-slate-900/40 rounded-lg border border-slate-800 overflow-hidden">
    <!-- Synchronized Toolbar (Matches DataTable) -->
    <div class="p-2 border-b border-slate-800 flex items-center justify-between bg-slate-900/60 h-10 shrink-0">
        <div class="flex items-center gap-4 flex-1">
            <!-- View Switcher -->
            <div class="flex items-center bg-slate-950 border border-slate-800 rounded-md p-0.5 mr-2">
                <button 
                    onclick={() => appState.viewMode = 'table'}
                    class="px-2 py-1 rounded text-[9px] uppercase font-bold transition-all flex items-center gap-1 {appState.viewMode === 'table' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-500 hover:text-slate-300'}"
                >
                    <LayoutGrid size={12} /> Table
                </button>
                <button 
                    onclick={() => appState.viewMode = 'tree'}
                    class="px-2 py-1 rounded text-[9px] uppercase font-bold transition-all flex items-center gap-1 {appState.viewMode === 'tree' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-500 hover:text-slate-300'}"
                >
                    <ListTree size={12} /> Tree
                </button>
            </div>

            <!-- Filter Bar (Synchronized) -->
            <div class="flex items-center gap-2 bg-slate-950/50 rounded-md border border-slate-800 px-2 py-0.5">
                <Filter size={12} class="text-slate-500" />
                <select bind:value={appState.filterField} class="bg-transparent border-none text-[10px] text-slate-300 focus:ring-0 cursor-pointer py-0 pr-6">
                    {#each columns as field}<option value={field}>{field}</option>{/each}
                </select>
                <div class="w-px h-3 bg-slate-800"></div>
                <div class="relative flex items-center">
                    <Search size={12} class="absolute left-0 text-slate-600" />
                    <input bind:value={appState.filterValue} type="text" placeholder="Filter tree..." class="bg-transparent border-none text-[11px] text-white focus:ring-0 py-0 pl-5 w-48" />
                </div>
                {#if appState.filterValue}
                    <button onclick={() => appState.filterValue = ''} class="text-slate-600 hover:text-white">
                        <X size={12} />
                    </button>
                {/if}
            </div>

            <button onclick={() => showQueryBuilder = !showQueryBuilder} class="p-1.5 rounded {showQueryBuilder ? 'bg-indigo-600 text-white' : 'bg-slate-800 text-slate-400'} text-[10px] uppercase font-bold ml-auto transition-colors">
                Server Query
            </button>
        </div>
        <button onclick={() => appState.executeQuery()} class="p-1.5 rounded hover:bg-slate-800 text-slate-400 ml-2">
            <RefreshCw size={14} class={appState.isLoading ? 'animate-spin' : ''} />
        </button>
    </div>

    <!-- Server Query Builder (Synchronized) -->
    {#if showQueryBuilder}
        <div class="p-3 border-b border-slate-800 bg-slate-900/80 space-y-2">
            <div class="flex items-center justify-between mb-2">
                <span class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Query Filters</span>
                <button onclick={() => appState.addQueryFilter()} class="p-1 rounded bg-indigo-600/20 text-indigo-400 text-[10px]">+ Condition</button>
            </div>
            {#each appState.activeQueries as query, i}
                <div class="flex items-center gap-2">
                    <select bind:value={query.field} class="bg-slate-950 border border-slate-800 rounded px-2 py-1 text-xs text-slate-200">
                        {#each columns as field}<option value={field}>{field}</option>{/each}
                    </select>
                    <select bind:value={query.operator} class="bg-slate-950 border border-slate-800 rounded px-2 py-1 text-xs text-slate-200">
                        {#each operators as op}<option value={op}>{op}</option>{/each}
                    </select>
                    <input bind:value={query.value} type="text" placeholder="Value..." class="flex-1 bg-slate-950 border border-slate-800 rounded px-2 py-1 text-xs text-slate-200" />
                    <button onclick={() => appState.removeQueryFilter(i)} class="p-1.5 text-slate-500 hover:text-red-400"><X size={14} /></button>
                </div>
            {/each}
            <div class="pt-2 flex justify-end">
                <button onclick={() => appState.executeQuery()} disabled={appState.isLoading} class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded text-xs font-bold">Run Query</button>
            </div>
        </div>
    {/if}

    <!-- Tree Header -->
    <div class="grid grid-cols-[1.5fr_2fr_1fr] border-b border-slate-800 bg-slate-900/80 text-[10px] uppercase font-black tracking-widest text-slate-400 px-4 py-2 shrink-0 shadow-sm">
        <div>Key / Document ID</div>
        <div>Value</div>
        <div>Type</div>
    </div>

    <!-- Tree Content -->
    <div class="flex-1 overflow-auto bg-slate-950/20 font-mono text-[12px] custom-scrollbar">
        <TreeTableRow 
            id={collectionName} 
            path={path} 
            type="collection" 
            depth={0} 
        />
        
        {#if appState.filteredDocuments.length === 0 && !appState.isLoading}
            <div class="p-20 text-center flex flex-col items-center gap-2 opacity-20">
                <Database size={48} />
                <p class="text-xs uppercase font-black tracking-tighter">No matching data in tree</p>
            </div>
        {/if}
    </div>

    <!-- Tree Footer -->
    <div class="p-2 border-t border-slate-800 bg-slate-900/50 flex items-center justify-between shrink-0">
        <div class="text-[10px] text-slate-500 px-2 uppercase font-black tracking-tighter">
            {appState.filteredDocuments.length} Documents in Hierarchy
        </div>
    </div>
</div>

<style>
    .custom-scrollbar::-webkit-scrollbar { width: 4px; }
    .custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.05); border-radius: 10px; }
</style>
