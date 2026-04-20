<script>
    import { appState } from '../state.svelte.js'
    import { ChevronRight, ChevronDown, Folder, FileText, Braces, Layers } from 'lucide-svelte'
    import TreeTableRow from './TreeTableRow.svelte'

    let { id, path, type, depth = 0, value = null, data = null } = $props()

    let isExpanded = $state(false)
    
    // Type detection and styling
    let displayType = $derived.by(() => {
        if (type === 'collection') return 'Collection'
        if (type === 'document') return 'Document'
        if (value === null) return 'Null'
        if (Array.isArray(value)) return 'Array'
        if (typeof value === 'object') return 'Map'
        if (typeof value === 'string') {
            // Check for ISO date string
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

    // Content for expansion
    let children = $derived.by(() => {
        if (!isExpanded) return []
        
        if (type === 'collection') {
            // If it's the root path of the view, use appState.documents
            if (path === appState.selectedPath) {
                return appState.documents.map(d => ({ id: d.id, path: d.path, type: 'document', data: d.data }))
            }
            // Otherwise check nested collection documents
            return (appState.collectionDocuments[path] || []).map(d => ({ id: d.id, path: d.path, type: 'document', data: d.data }))
        }
        
        if (type === 'document') {
            const fields = Object.entries(data || {}).map(([k, v]) => ({
                id: k,
                path: `${path}/${k}`,
                type: 'field',
                value: v
            }))
            // Also include sub-collections if any
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
                value: v
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
                    // Trigger load for nested collection
                    await appState.toggleCollectionExpansion(path)
                }
            } else if (type === 'document') {
                if (!appState.docSubCollections[path]) {
                    await appState.toggleDocExpansion(path)
                }
            }
        }
    }

    function handleRowClick() {
        if (type === 'document' || (type === 'field' && value && typeof value === 'object')) {
            // Load into full inspector
            if (type === 'document') {
                appState.selectedDoc = { id, path, data }
            }
        }
    }
</script>

<div class="flex flex-col border-b border-slate-800/30">
    <div 
        class="grid grid-cols-[1.5fr_2fr_1fr] hover:bg-slate-800/40 transition-colors group cursor-pointer py-1.5 px-4"
        onclick={handleRowClick}
    >
        <!-- Key Column -->
        <div class="flex items-center gap-1 overflow-hidden" style="padding-left: {depth * 20}px">
            {#if displayType === 'Collection' || displayType === 'Document' || displayType === 'Map' || displayType === 'Array'}
                <button 
                    onclick={toggleExpand}
                    class="p-0.5 hover:bg-slate-700 rounded text-slate-500 transition-colors shrink-0"
                >
                    {#if isExpanded}
                        <ChevronDown size={14} />
                    {:else}
                        <ChevronRight size={14} />
                    {/if}
                </button>
            {:else}
                <div class="w-5 shrink-0"></div>
            {/if}

            {#if type === 'collection'}
                <Folder size={14} class="text-blue-500 shrink-0" />
            {:else if type === 'document'}
                <FileText size={14} class="text-purple-500 shrink-0" />
            {:else if displayType === 'Map' || displayType === 'Array'}
                <Braces size={14} class="text-slate-400 shrink-0" />
            {:else}
                <Layers size={14} class="text-slate-600 shrink-0 opacity-40" />
            {/if}

            <span class="truncate {type === 'document' ? 'text-purple-300 font-bold' : 'text-slate-300'}">{id}</span>
        </div>

        <!-- Value Column -->
        <div class="flex items-center px-2 overflow-hidden border-l border-slate-800/30">
            {#if displayType !== 'Map' && displayType !== 'Array' && displayType !== 'Collection' && displayType !== 'Document'}
                <span class="truncate {typeColor}">{String(value)}</span>
            {:else if displayType === 'Array'}
                <span class="text-[10px] text-slate-600 font-bold uppercase">{value.length} items</span>
            {:else if displayType === 'Map'}
                <span class="text-[10px] text-slate-600 font-bold uppercase">{Object.keys(value).length} keys</span>
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
                />
            {/each}
        </div>
    {/if}
</div>
