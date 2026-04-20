<script>
    import { appState } from '../state.svelte.js'
    import TreeTableRow from './TreeTableRow.svelte'
    import { ListTree, LayoutGrid, RefreshCw } from 'lucide-svelte'

    let { path = '' } = $props()
    
    // The top-level collection nodes
    let collectionName = $derived(path.split('/').pop() || 'Root')
</script>

<div class="flex flex-col h-full bg-slate-900/40 rounded-lg border border-slate-800 overflow-hidden">
    <!-- View Switcher Toolbar -->
    <div class="p-2 border-b border-slate-800 flex items-center justify-between bg-slate-900/60 h-10 shrink-0">
        <div class="flex items-center bg-slate-950 border border-slate-800 rounded-md p-0.5">
            <button 
                onclick={() => appState.viewMode = 'table'}
                class="px-3 py-1 rounded text-[10px] uppercase font-bold transition-all flex items-center gap-1.5 {appState.viewMode === 'table' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-500 hover:text-slate-300'}"
            >
                <LayoutGrid size={12} /> Table
            </button>
            <button 
                onclick={() => appState.viewMode = 'tree'}
                class="px-3 py-1 rounded text-[10px] uppercase font-bold transition-all flex items-center gap-1.5 {appState.viewMode === 'tree' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-500 hover:text-slate-300'}"
            >
                <ListTree size={12} /> Tree
            </button>
        </div>

        <button onclick={() => appState.executeQuery()} class="p-1.5 rounded hover:bg-slate-800 text-slate-400">
            <RefreshCw size={14} class={appState.isLoading ? 'animate-spin' : ''} />
        </button>
    </div>

    <!-- Tree Header -->
    <div class="grid grid-cols-[1.5fr_2fr_1fr] border-b border-slate-800 bg-slate-900/80 text-[10px] uppercase font-bold tracking-wider text-slate-500 px-4 py-2 shrink-0">
        <div>Key</div>
        <div>Value</div>
        <div>Type</div>
    </div>

    <!-- Tree Content -->
    <div class="flex-1 overflow-auto bg-slate-950/20 font-mono text-[12px]">
        <TreeTableRow 
            id={collectionName} 
            path={path} 
            type="collection" 
            depth={0} 
        />
    </div>
</div>
