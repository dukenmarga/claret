<script>
    import { appState } from '../state.svelte.js'
    import { Trash2, Copy, Database, Clipboard, Save, X, Search, Filter, Plus, ChevronUp, ChevronDown, ChevronRight, Play, RefreshCw, Folder, FileText, LayoutGrid, ListTree, Braces } from 'lucide-svelte'

    let { path = '' } = $props()
    
    let docs = $derived.by(() => {
        const raw = appState.documents || []
        const val = appState.filterValue || ''
        const fld = appState.filterField || 'id'
        if (!val.trim()) return raw
        const q = val.toLowerCase().trim()
        return raw.filter(d => {
            if (!d) return false
            if (fld === 'id') return d.id.toLowerCase().includes(q)
            return String(d.data?.[fld] || '').toLowerCase().includes(q)
        })
    })

    let columns = $derived.by(() => {
        const fields = new Set(['id'])
        if (appState.documents && appState.documents.length > 0) {
            appState.documents.forEach(doc => {
                if (doc?.data) Object.keys(doc.data).forEach(key => fields.add(key))
            })
        }
        return Array.from(fields).sort()
    })

    function formatValue(val) {
        if (val === null) return 'null'
        if (typeof val === 'object') return JSON.stringify(val)
        return String(val)
    }

    let targetCollection = $state('')
    let showDuplicateDialog = $state(false)
    let showQueryBuilder = $state(false)
    const operators = ['==', '!=', '>', '<', '>=', '<=', 'array-contains', 'in']

    let editingCell = $state(null) 
    let editInput = $state(null)
    let nestedViewMode = $state({}) 

    async function handleDuplicate() {
        if (!targetCollection) return
        await appState.duplicateTo(targetCollection)
        showDuplicateDialog = false
    }

    function startEditing(doc, field, isId = false) {
        editingCell = { docId: doc.id, field, value: isId ? doc.id : formatValue(doc.data[field]), originalDoc: doc, isId }
        setTimeout(() => editInput?.focus(), 0)
    }

    async function saveEdit() {
        if (!editingCell) return
        const { docId, field, value, originalDoc, isId } = editingCell
        if (isId) {
            if (value !== originalDoc.id) await appState.renameDocument(originalDoc, value)
            editingCell = null; return
        }
        let parsedValue = value
        if (value.toLowerCase() === 'true') parsedValue = true
        else if (value.toLowerCase() === 'false') parsedValue = false
        else if (value === 'null') parsedValue = null
        else if (!isNaN(Number(value)) && value.trim() !== '') parsedValue = Number(value)
        else if (value.startsWith('{') || value.startsWith('[')) {
            try { parsedValue = JSON.parse(value) } catch (e) {}
        }
        const newData = { ...originalDoc.data, [field]: parsedValue }
        await appState.saveDocument({ ...originalDoc, data: newData })
        editingCell = null
    }

    function handleKeyDown(e) {
        if (e.key === 'Enter') saveEdit()
        if (e.key === 'Escape') editingCell = null
    }
</script>

<div class="flex flex-col h-full bg-slate-900/40 rounded-lg border border-slate-800 overflow-hidden">
    <!-- Toolbar -->
    <div class="p-2 border-b border-slate-800 flex items-center justify-between bg-slate-900/60 h-10 shrink-0">
        <div class="flex items-center gap-4 flex-1">
            <!-- View Switcher (Added here so it's visible in Table mode) -->
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

            <div class="flex items-center gap-2 bg-slate-950/50 rounded-md border border-slate-800 px-2 py-0.5">
                <Filter size={12} class="text-slate-500" />
                <select bind:value={appState.filterField} class="bg-transparent border-none text-[10px] text-slate-300 focus:ring-0 cursor-pointer py-0 pr-6">
                    {#each columns as field}<option value={field}>{field}</option>{/each}
                </select>
                <div class="w-px h-3 bg-slate-800"></div>
                <div class="relative flex items-center">
                    <Search size={12} class="absolute left-0 text-slate-600" />
                    <input bind:value={appState.filterValue} type="text" placeholder="Filter results..." class="bg-transparent border-none text-[11px] text-white focus:ring-0 py-0 pl-5 w-48" />
                </div>
            </div>

            {#if appState.selectedDocIds.length > 0}
                <div class="flex items-center gap-2 border-l border-slate-800 pl-4 h-6">
                    <button onclick={() => appState.deleteSelected()} class="p-1.5 rounded hover:bg-red-500/10 text-red-400 text-[10px] font-bold uppercase flex items-center gap-1">
                        <Trash2 size={12} /> Delete ({appState.selectedDocIds.length})
                    </button>
                    <button onclick={() => { targetCollection = appState.selectedPath; showDuplicateDialog = true }} class="p-1.5 rounded hover:bg-indigo-500/10 text-indigo-400 text-[10px] font-bold uppercase">
                        <Copy size={12} /> Duplicate
                    </button>
                </div>
            {/if}

            <button onclick={() => showQueryBuilder = !showQueryBuilder} class="p-1.5 rounded {showQueryBuilder ? 'bg-indigo-600 text-white' : 'bg-slate-800 text-slate-400'} text-[10px] uppercase font-bold ml-auto">
                Server Query
            </button>
        </div>
        <button onclick={() => appState.executeQuery()} class="p-1.5 rounded hover:bg-slate-800 text-slate-400 ml-2">
            <RefreshCw size={14} class={appState.isLoading ? 'animate-spin' : ''} />
        </button>
    </div>

    <!-- Server Query Builder -->
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

    <div class="flex-1 overflow-auto bg-slate-950/20">
        <table class="w-full text-left text-xs border-collapse">
            <thead class="sticky top-0 bg-slate-900 z-10 shadow-sm">
                <tr class="text-[10px] uppercase tracking-wider text-slate-500">
                    <th class="px-2 py-2 border-b border-slate-800 w-8 text-center">
                        <input type="checkbox" checked={appState.selectedDocIds.length > 0 && appState.selectedDocIds.length === docs.length} onchange={(e) => appState.selectAll(e.target.checked)} class="rounded border-slate-700 bg-slate-950 text-indigo-600" />
                    </th>
                    <th class="px-3 py-2 border-b border-slate-800 w-24">Actions</th>
                    <th class="px-3 py-2 border-b border-slate-800 w-48 cursor-pointer hover:bg-slate-800" onclick={() => appState.setSort('id')}>ID</th>
                    {#each columns.filter(f => f !== 'id') as col}
                        <th class="px-3 py-2 border-b border-slate-800 min-w-[120px] max-w-64 truncate cursor-pointer hover:bg-slate-800" onclick={() => appState.setSort(col)}>{col}</th>
                    {/each}
                </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/50">
                {#each docs as doc (doc.id)}
                    <tr class="hover:bg-slate-800/30 transition-colors group cursor-pointer {appState.selectedDoc?.id === doc.id ? 'bg-indigo-500/5 ring-1 ring-inset ring-indigo-500/20' : ''}" onclick={() => appState.selectedDoc = doc}>
                        <td class="px-2 py-1 text-center" onclick={(e) => e.stopPropagation()}>
                            <input type="checkbox" checked={appState.selectedDocIds.includes(doc.id)} onchange={() => appState.toggleSelection(doc.id)} class="rounded border-slate-700 bg-slate-950 text-indigo-600" />
                        </td>
                        <td class="px-3 py-1">
                            <div class="flex items-center gap-1">
                                <button onclick={(e) => { e.stopPropagation(); appState.toggleDocExpansion(doc.path) }} class="p-1 rounded hover:bg-indigo-500/10 {appState.expandedDocs.includes(doc.path) ? 'bg-indigo-500/20 text-indigo-400' : 'text-slate-500'}">
                                    {#if appState.expandedDocs.includes(doc.path)}<ChevronDown size={13} />{:else}<ChevronRight size={13} />{/if}
                                </button>
                                <button onclick={(e) => { e.stopPropagation(); navigator.clipboard.writeText(JSON.stringify(doc.data)) }} class="p-1 rounded hover:bg-slate-700 text-slate-400"><Clipboard size={13} /></button>
                                <button onclick={(e) => { e.stopPropagation(); appState.deleteDocument(doc.path) }} class="p-1 rounded hover:bg-red-500/10 text-red-400/70 hover:text-red-400"><Trash2 size={13} /></button>
                            </div>
                        </td>
                        <td class="px-3 py-1 text-slate-300 font-mono text-[11px] truncate relative group/cell" ondblclick={() => startEditing(doc, 'id', true)}>
                            {#if editingCell?.docId === doc.id && editingCell?.isId}
                                <input bind:this={editInput} bind:value={editingCell.value} onkeydown={(e) => e.key === 'Enter' && saveEdit()} onblur={saveEdit} class="w-full bg-slate-950 border-none text-white focus:ring-0 text-[11px] font-mono p-0" />
                            {:else}{doc.id}{/if}
                        </td>
                        {#each columns.filter(f => f !== 'id') as col}
                            <td class="px-3 py-1 text-slate-400 truncate max-w-64 relative group/cell" ondblclick={() => startEditing(doc, col)}>
                                {#if editingCell?.docId === doc.id && editingCell?.field === col}
                                    <input bind:this={editInput} bind:value={editingCell.value} onkeydown={(e) => e.key === 'Enter' && saveEdit()} onblur={saveEdit} class="w-full bg-slate-950 border-none text-white focus:ring-0 text-xs p-0" />
                                {:else}{formatValue(doc.data[col])}{/if}
                            </td>
                        {/each}
                    </tr>
                    {#if appState.expandedDocs.includes(doc.path)}
                        <tr class="bg-slate-900/60 shadow-inner">
                            <td colspan="100" class="p-0 border-l-2 border-indigo-500/50">
                                <div class="pl-12 pr-4 py-3 space-y-4">
                                    <div class="flex items-center gap-4 border-b border-slate-800/50 pb-2">
                                        <button onclick={() => nestedViewMode[doc.path] = 'collections'} class="text-[10px] uppercase font-bold tracking-widest flex items-center gap-1.5 transition-colors {!nestedViewMode[doc.path] || nestedViewMode[doc.path] === 'collections' ? 'text-indigo-400' : 'text-slate-500 hover:text-slate-300'}"><Folder size={12} /> Sub-collections</button>
                                        <button onclick={() => nestedViewMode[doc.path] = 'data'} class="text-[10px] uppercase font-bold tracking-widest flex items-center gap-1.5 transition-colors {nestedViewMode[doc.path] === 'data' ? 'text-indigo-400' : 'text-slate-500 hover:text-slate-300'}"><Braces size={12} /> Data Tree</button>
                                    </div>
                                    {#if !nestedViewMode[doc.path] || nestedViewMode[doc.path] === 'collections'}
                                        <div class="space-y-2">
                                            {#if appState.docSubCollections[doc.path]}
                                                {#each appState.docSubCollections[doc.path] as subColl}
                                                    {@const collPath = `${doc.path}/${subColl}`}
                                                    <div class="flex flex-col">
                                                        <button onclick={() => appState.toggleCollectionExpansion(collPath)} class="flex items-center gap-2 text-xs text-indigo-400 hover:text-indigo-300 transition-colors py-1 group/coll w-fit">
                                                            {#if appState.expandedCollections.includes(collPath)}<ChevronDown size={12} />{:else}<ChevronRight size={12} />{/if}
                                                            <span class="font-medium group-hover/coll:underline">{subColl}</span>
                                                        </button>
                                                        {#if appState.expandedCollections.includes(collPath)}
                                                            <div class="pl-6 border-l border-slate-800 mt-1 mb-2 space-y-1">
                                                                {#if appState.collectionDocuments[collPath]}
                                                                    {#each appState.collectionDocuments[collPath] as nestedDoc}
                                                                        <button onclick={() => appState.selectedDoc = nestedDoc} class="flex items-center gap-2 group/nested text-slate-400 hover:text-white py-1 px-2 rounded-md hover:bg-slate-800/40 w-full text-left transition-all">
                                                                            <FileText size={10} class="text-slate-600" /><span class="text-[11px] font-mono flex-1">{nestedDoc.id}</span><span class="text-[9px] text-slate-600 opacity-0 group-hover/nested:opacity-100 uppercase font-bold tracking-tighter">Inspect</span>
                                                                        </button>
                                                                    {/each}
                                                                    {#if appState.collectionDocuments[collPath].length === 0}<div class="text-[10px] text-slate-600 italic px-2">Empty collection</div>{/if}
                                                                {:else}<div class="text-[10px] text-slate-600 animate-pulse px-2">Fetching...</div>{/if}
                                                            </div>
                                                        {/if}
                                                    </div>
                                                {/each}
                                                {#if appState.docSubCollections[doc.path].length === 0}<div class="text-[10px] text-slate-600 italic">No sub-collections</div>{/if}
                                            {:else}<div class="text-[10px] text-slate-600 animate-pulse">Loading...</div>{/if}
                                        </div>
                                    {:else}
                                        <div class="bg-slate-950/40 rounded-lg p-4 font-mono text-[11px] border border-slate-800/50">
                                            {#each Object.entries(doc.data) as [key, value]}
                                                <div class="flex gap-2 py-0.5 group/field"><span class="text-indigo-400 font-bold w-32 shrink-0 truncate">{key}</span><span class="text-slate-500 mr-1">:</span><span class="text-emerald-400 break-all">{formatValue(value)}</span></div>
                                            {/each}
                                            <div class="mt-4 pt-2 border-t border-slate-800/50 flex justify-end"><button onclick={() => appState.selectedDoc = doc} class="text-[9px] text-indigo-500 uppercase font-bold hover:text-indigo-400">Open Full Inspector</button></div>
                                        </div>
                                    {/if}
                                </div>
                            </td>
                        </tr>
                    {/if}
                {/each}
            </tbody>
        </table>
        {#if docs.length === 0 && !appState.isLoading}
            <div class="p-24 text-center text-slate-500 flex flex-col items-center justify-center gap-2 w-full h-full min-h-[300px]">
                <LayoutGrid size={48} class="opacity-10 mb-2" />
                <p class="font-medium text-slate-400">No documents found in this view</p>
                <button onclick={() => appState.executeQuery()} class="mt-4 px-4 py-2 bg-slate-800 rounded text-xs hover:bg-slate-700 transition-colors">Force Reload</button>
            </div>
        {/if}
    </div>

    <div class="p-2 border-t border-slate-800 flex items-center justify-between bg-slate-900/50 shrink-0">
        <div class="text-[10px] text-slate-500 px-2 uppercase font-bold tracking-tighter">
            {docs.length} Documents Visible (Total: {appState.documents.length})
        </div>
        {#if appState.hasMore}
            <button onclick={() => appState.executeQuery(true)} disabled={appState.isLoading} class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs rounded border border-slate-700 transition-colors flex items-center gap-2">
                <Plus size={12} /> Load More
            </button>
        {/if}
    </div>
</div>

{#if showDuplicateDialog}
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[100] p-4">
        <div class="bg-slate-900 border border-slate-800 rounded-xl p-4 space-y-4 w-full max-w-sm shadow-2xl">
            <h2 class="text-sm font-bold text-white flex items-center gap-2"><Copy size={16} class="text-indigo-400" />Duplicate {appState.selectedDocIds.length} Document(s)</h2>
            <input bind:value={targetCollection} type="text" placeholder="Target Path..." class="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-1 focus:ring-indigo-500" />
            <div class="flex gap-3">
                <button onclick={() => showDuplicateDialog = false} class="flex-1 py-1.5 rounded text-xs text-slate-400 hover:bg-slate-800">Cancel</button>
                <button onclick={handleDuplicate} class="flex-1 py-1.5 rounded text-xs bg-indigo-600 text-white font-bold">Duplicate</button>
            </div>
        </div>
    </div>
{/if}

<style>
    table { border-spacing: 0; }
</style>
