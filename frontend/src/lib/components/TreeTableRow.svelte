<script>
    import { appState } from '../state.svelte.js'
    import { ChevronRight, ChevronDown, Folder, FileText, Braces, Layers, Edit2, Check, X } from 'lucide-svelte'
    import TreeTableRow from './TreeTableRow.svelte'

    let { id, path, type, depth = 0, value = null, data = null, parentDoc = null } = $props()

    let isExpanded = $state(false)
    let isEditing = $state(false)
    let editValue = $state('')
    let editInput = $state(null)
    
    // Type detection and styling
    let displayType = $derived.by(() => {
        if (type === 'collection') return 'Collection'
        if (type === 'document') return 'Document'
        if (value === null) return 'Null'
        if (Array.isArray(value)) return 'Array'
        if (typeof value === 'object') return 'Map'
        if (typeof value === 'string') {
            if (/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/.test(value)) return 'Timestamp'
            return 'String'
        }
        if (typeof value === 'number') return 'Number'
        if (typeof value === 'boolean') return 'Boolean'
        return 'Value'
    })

    let typeColor = $derived.by(() => {
        switch (displayType) {
            case 'Collection': return 'text-blue-400'
            case 'Document': return 'text-purple-400'
            case 'Map': return 'text-slate-300'
            case 'Array': return 'text-slate-300'
            case 'String': return 'text-indigo-400'
            case 'Boolean': return 'text-emerald-400'
            case 'Number': return 'text-amber-400'
            case 'Timestamp': return 'text-rose-400'
            default: return 'text-slate-500'
        }
    })

    let children = $derived.by(() => {
        if (!isExpanded) return []
        if (type === 'collection') {
            if (path === appState.selectedPath) {
                return appState.documents.map(d => ({ id: d.id, path: d.path, type: 'document', data: d.data, parentDoc: d }))
            }
            return (appState.collectionDocuments[path] || []).map(d => ({ id: d.id, path: d.path, type: 'document', data: d.data, parentDoc: d }))
        }
        if (type === 'document') {
            const fields = Object.entries(data || {}).map(([k, v]) => ({
                id: k,
                path: `${path}/${k}`,
                type: 'field',
                value: v,
                parentDoc: { id, path, data }
            }))
            const subs = (appState.docSubCollections[path] || []).map(s => ({
                id: s,
                path: `${path}/${s}`,
                type: 'collection'
            }))
            return [...subs, ...fields]
        }
        if (type === 'field' && (displayType === 'Map' || displayType === 'Array')) {
            return Object.entries(value).map(([k, v]) => ({
                id: k,
                path: `${path}/${k}`,
                type: 'field',
                value: v,
                parentDoc
            }))
        }
        return []
    })

    async function toggleExpand(e) {
        e.stopPropagation()
        isExpanded = !isExpanded
        if (isExpanded) {
            if (type === 'collection') {
                if (path !== appState.selectedPath && !appState.collectionDocuments[path]) {
                    await appState.toggleCollectionExpansion(path)
                }
            } else if (type === 'document') {
                if (!appState.docSubCollections[path]) {
                    await appState.toggleDocExpansion(path)
                }
            }
        }
    }

    function startEditing(e) {
        e.stopPropagation()
        if (type !== 'field' || displayType === 'Map' || displayType === 'Array') return
        editValue = String(value)
        isEditing = true
        setTimeout(() => editInput?.focus(), 0)
    }

    async function saveEdit() {
        if (!parentDoc) return
        
        let parsedValue = editValue
        if (displayType === 'Boolean') parsedValue = editValue.toLowerCase() === 'true'
        else if (displayType === 'Number') parsedValue = Number(editValue)
        else if (displayType === 'Null') parsedValue = null

        // Need to update the nested path in the parent document data
        // For simplicity in this recursive structure, we update the top-level field
        // Since Firestore Set with Merge handles nested keys if we use dot notation, 
        // but here we are working with the local 'data' object.
        
        const newData = { ...parentDoc.data }
        newData[id] = parsedValue // This only works for 1-level depth fields currently
        
        await appState.saveDocument({ ...parentDoc, data: newData })
        isEditing = false
    }

    function handleKeyDown(e) {
        if (e.key === 'Enter') saveEdit()
        if (e.key === 'Escape') isEditing = false
    }

    function handleRowClick() {
        if (type === 'document') appState.selectedDoc = { id, path, data }
    }
</script>

<div class="flex flex-col border-b border-slate-800/30">
    <div 
        class="grid grid-cols-[1.5fr_2fr_1fr] hover:bg-slate-800/40 transition-colors group cursor-pointer py-1.5 px-4"
        onclick={handleRowClick}
        ondblclick={startEditing}
    >
        <!-- Key Column -->
        <div class="flex items-center gap-1 overflow-hidden" style="padding-left: {depth * 20}px">
            {#if displayType === 'Collection' || displayType === 'Document' || displayType === 'Map' || displayType === 'Array'}
                <button onclick={toggleExpand} class="p-0.5 hover:bg-slate-700 rounded text-slate-500 transition-colors shrink-0">
                    {#if isExpanded}<ChevronDown size={14} />{:else}<ChevronRight size={14} />{/if}
                </button>
            {:else}
                <div class="w-5 shrink-0"></div>
            {/if}

            {#if type === 'collection'}<Folder size={14} class="text-blue-500 shrink-0" />
            {:else if type === 'document'}<FileText size={14} class="text-purple-500 shrink-0" />
            {:else if displayType === 'Map' || displayType === 'Array'}<Braces size={14} class="text-slate-400 shrink-0" />
            {:else}<Layers size={14} class="text-slate-600 shrink-0 opacity-40" />{/if}

            <span class="truncate {type === 'document' ? 'text-purple-300 font-bold' : 'text-slate-300'}">{id}</span>
        </div>

        <!-- Value Column -->
        <div class="flex items-center px-2 overflow-hidden border-l border-slate-800/30 relative group/val">
            {#if isEditing}
                <div class="flex-1 flex items-center gap-1">
                    <input 
                        bind:this={editInput}
                        bind:value={editValue}
                        onkeydown={handleKeyDown}
                        class="flex-1 bg-slate-950 border-none text-white focus:ring-0 p-0 text-[12px] font-mono"
                    />
                    <button onclick={saveEdit} class="text-emerald-500 hover:text-emerald-400"><Check size={14} /></button>
                    <button onclick={() => isEditing = false} class="text-slate-500 hover:text-white"><X size={14} /></button>
                </div>
            {:else}
                {#if displayType !== 'Map' && displayType !== 'Array' && displayType !== 'Collection' && displayType !== 'Document'}
                    <span class="truncate {typeColor}">{String(value)}</span>
                    <button onclick={startEditing} class="absolute right-2 opacity-0 group-hover/val:opacity-100 p-1 rounded hover:bg-slate-700 text-slate-500 transition-all">
                        <Edit2 size={10} />
                    </button>
                {:else if displayType === 'Array'}<span class="text-[10px] text-slate-600 font-bold uppercase">{value.length} items</span>
                {:else if displayType === 'Map'}<span class="text-[10px] text-slate-600 font-bold uppercase">{Object.keys(value).length} keys</span>{/if}
            {/if}
        </div>

        <!-- Type Column -->
        <div class="flex items-center px-2 border-l border-slate-800/30">
            <span class="text-[10px] uppercase font-bold tracking-tighter {typeColor} opacity-80">{displayType}</span>
        </div>
    </div>

    {#if isExpanded}
        <div class="flex flex-col">
            {#each children as child}
                <TreeTableRow 
                    id={child.id}
                    path={child.path}
                    type={child.type}
                    depth={depth + 1}
                    value={child.value}
                    data={child.data}
                    parentDoc={child.parentDoc}
                />
            {/each}
        </div>
    {/if}
</div>
