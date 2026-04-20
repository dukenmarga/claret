<script>
    import { X, AlertTriangle, Trash2 } from 'lucide-svelte'
    import { appState } from '../state.svelte.js'

    let confirmationText = $state('')
    let isValid = $derived(confirmationText === 'DELETE')

    async function handleConfirm() {
        if (!isValid) return
        await appState.confirmDelete()
    }

    function handleClose() {
        appState.showDeleteDialog = false
        appState.confirmDeletionList = []
    }
</script>

<div class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[100] p-4">
    <div class="bg-slate-900 border border-slate-800 rounded-xl shadow-2xl w-full max-w-lg overflow-hidden animate-in fade-in zoom-in duration-200">
        <div class="p-4 border-b border-slate-800 flex items-center justify-between bg-red-500/5">
            <h2 class="text-sm font-bold text-red-400 flex items-center gap-2">
                <AlertTriangle size={18} />
                Confirm Deletion
            </h2>
            <button onclick={handleClose} class="text-slate-500 hover:text-white transition-colors">
                <X size={20} />
            </button>
        </div>

        <div class="p-6 space-y-4">
            <p class="text-xs text-slate-400 leading-relaxed">
                You are about to permanently delete <span class="text-white font-bold">{appState.confirmDeletionList.length}</span> document(s). 
                This action cannot be undone.
            </p>

            <div class="bg-slate-950/50 rounded-lg border border-slate-800 max-h-48 overflow-y-auto p-2">
                <table class="w-full text-left text-[10px]">
                    <thead class="sticky top-0 bg-slate-900 z-10 text-slate-500 uppercase">
                        <tr>
                            <th class="px-2 py-1">Document ID</th>
                            <th class="px-2 py-1">Path</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-800/50">
                        {#each appState.confirmDeletionList as doc}
                            <tr class="text-slate-300">
                                <td class="px-2 py-1 font-mono truncate max-w-[120px]">{doc.id}</td>
                                <td class="px-2 py-1 font-mono text-slate-500 truncate">{doc.path}</td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>

            <div>
                <label class="block text-[10px] font-semibold text-slate-500 uppercase tracking-wider mb-2">
                    Type <span class="text-red-400">DELETE</span> to confirm
                </label>
                <input 
                    bind:value={confirmationText} 
                    type="text" 
                    placeholder="DELETE" 
                    class="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-red-500/30 text-red-400 placeholder:text-slate-800"
                    onkeydown={(e) => e.key === 'Enter' && isValid && handleConfirm()}
                />
            </div>
        </div>

        <div class="p-4 bg-slate-900/50 border-t border-slate-800 flex gap-3">
            <button 
                onclick={handleClose} 
                class="flex-1 py-2 rounded-lg text-xs font-medium text-slate-400 hover:bg-slate-800 transition-colors"
            >
                Cancel
            </button>
            <button 
                onclick={handleConfirm}
                disabled={!isValid || appState.isLoading}
                class="flex-1 py-2 rounded-lg text-xs font-medium bg-red-600 hover:bg-red-500 disabled:opacity-30 disabled:cursor-not-allowed text-white shadow-lg shadow-red-900/20 transition-all flex items-center justify-center gap-2"
            >
                <Trash2 size={14} />
                Delete Permanently
            </button>
        </div>
    </div>
</div>
