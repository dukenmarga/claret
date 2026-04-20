<script>
    import { appState } from '../state.svelte.js'
    import { ChevronRight, ChevronDown, Database, FileText, Folder } from 'lucide-svelte'
    import TreeItem from './TreeItem.svelte'

    let { path, id, type, depth = 0 } = $props()
    
    let node = $derived(appState.treeData[path])
    let isSelected = $derived(appState.selectedPath === path)

    async function handleToggle(e) {
        e.stopPropagation()
        await appState.toggleNode(path, type)
    }

    function handleSelect(e) {
        e.stopPropagation()
        if (type === 'collection') {
            appState.selectCollection(path)
        } else {
            // If it's a doc, maybe we want to select it in the table?
            // For now, just expand it to see sub-collections
            appState.toggleNode(path, 'doc')
        }
    }
</script>

<div class="flex flex-col">
    <div 
        class="w-full flex items-center gap-1.5 px-2 py-1 rounded-md text-xs transition-colors hover:bg-slate-800/50 group cursor-pointer {isSelected ? 'bg-indigo-500/10 text-indigo-400' : 'text-slate-400'}"
        style="padding-left: {depth * 12 + 8}px"
        onclick={handleSelect}
    >
        <button 
            class="p-0.5 hover:bg-slate-700 rounded transition-colors text-slate-600"
            onclick={handleToggle}
        >
            {#if node?.expanded}
                <ChevronDown size={12} />
            {:else}
                <ChevronRight size={12} />
            {/if}
        </button>

        {#if type === 'collection'}
            <Folder size={12} class={isSelected ? 'text-indigo-400' : 'text-slate-500'} />
        {:else}
            <FileText size={12} class="text-slate-600" />
        {/if}

        <span class="truncate flex-1 text-left {isSelected ? 'font-bold' : ''}">{id}</span>
    </div>

    {#if node?.expanded}
        <div class="flex flex-col">
            {#each node.children as child}
                <TreeItem 
                    path={child.path} 
                    id={child.id} 
                    type={child.type} 
                    depth={depth + 1} 
                />
            {/each}
            {#if node.children.length === 0 && appState.isLoading}
                <div class="py-1 px-4 text-[10px] text-slate-600 animate-pulse" style="padding-left: {(depth + 1) * 12 + 24}px">
                    Loading...
                </div>
            {/if}
        </div>
    {/if}
</div>
