<script>
    import { onMount } from 'svelte'
    import { GetUsers, DeleteUser } from '../../../wailsjs/go/main/App'
    import { Loader2, Trash2, Mail, Calendar, Key, UserCheck, Shield, RefreshCw, Users } from 'lucide-svelte'
    import { toast } from 'svelte-sonner'
    import { appState } from '../state.svelte.js'

    let users = $state([])
    let isLoading = $state(false)
    let error = $state(null)

    async function fetchUsers() {
        if (!appState.currentProjectId) return
        
        isLoading = true
        error = null
        try {
            const resp = await GetUsers(100)
            if (resp.success) {
                users = resp.data || []
            } else {
                error = resp.error
                toast.error("Auth fetch failed: " + resp.error)
            }
        } catch (e) {
            error = e.toString()
            toast.error(e.toString())
        } finally {
            isLoading = false
        }
    }

    async function handleDeleteUser(uid) {
        if (!confirm('Are you sure you want to delete this user?')) return
        try {
            const resp = await DeleteUser(uid)
            if (resp.success) {
                users = users.filter(u => u.uid !== uid)
                toast.success('User deleted')
            } else {
                toast.error(`Delete failed: ${resp.error}`)
            }
        } catch (e) {
            toast.error(e.toString())
        }
    }

    // Refresh users whenever the project changes
    $effect(() => {
        if (appState.currentProjectId) {
            fetchUsers()
        }
    })

    function formatDate(ts) {
        if (!ts) return 'Never'
        return new Date(ts).toLocaleDateString() + ' ' + new Date(ts).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }
</script>

<div class="flex flex-col h-full space-y-4">
    <div class="flex items-center justify-between bg-slate-900/40 p-4 rounded-xl border border-slate-800 shadow-sm">
        <div class="flex items-center gap-3">
            <div class="p-2 bg-indigo-500/10 rounded-lg text-indigo-400">
                <Shield size={20} />
            </div>
            <div>
                <h2 class="text-sm font-black text-white uppercase tracking-widest">Authentication</h2>
                <p class="text-[10px] text-slate-500 font-medium uppercase">Manage users for {appState.currentProject?.name}</p>
            </div>
        </div>
        <button 
            onclick={fetchUsers} 
            disabled={isLoading}
            class="p-2 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white transition-all active:rotate-180"
        >
            <RefreshCw size={16} class={isLoading ? 'animate-spin' : ''} />
        </button>
    </div>

    {#if isLoading && users.length === 0}
        <div class="flex-1 flex flex-col items-center justify-center text-slate-500 gap-4">
            <Loader2 size={40} class="animate-spin text-indigo-500" />
            <span class="text-xs font-bold uppercase tracking-widest animate-pulse">Retrieving Users...</span>
        </div>
    {:else if error}
        <div class="p-8 text-center bg-red-500/5 border border-red-500/10 rounded-2xl">
            <p class="text-red-400 text-sm font-bold mb-2">Error Loading Users</p>
            <p class="text-red-400/60 text-xs font-mono">{error}</p>
        </div>
    {:else}
        <div class="flex-1 border border-slate-800 rounded-xl overflow-hidden bg-slate-900/20 shadow-inner">
            <div class="overflow-auto h-full custom-scrollbar">
                <table class="w-full text-left text-xs border-collapse">
                    <thead class="sticky top-0 bg-slate-900 z-10 shadow-sm">
                        <tr class="text-[10px] uppercase font-black tracking-widest text-slate-500">
                            <th class="px-6 py-4 border-b border-slate-800">User Identity</th>
                            <th class="px-6 py-4 border-b border-slate-800">Unique ID (UID)</th>
                            <th class="px-6 py-4 border-b border-slate-800">Signed Up</th>
                            <th class="px-6 py-4 border-b border-slate-800">Last Active</th>
                            <th class="px-6 py-4 border-b border-slate-800 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-800/50">
                        {#each users as user}
                            <tr class="hover:bg-white/5 transition-colors group">
                                <td class="px-6 py-4">
                                    <div class="flex items-center gap-3">
                                        <div class="w-8 h-8 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-400 font-bold overflow-hidden shadow-sm">
                                            {#if user.photoURL}
                                                <img src={user.photoURL} alt="" class="w-full h-full object-cover" />
                                            {:else}
                                                {user.email ? user.email[0].toUpperCase() : '?'}
                                            {/if}
                                        </div>
                                        <div class="flex flex-col">
                                            <span class="text-slate-100 font-bold text-sm tracking-tight">{user.email || 'No Email'}</span>
                                            <span class="text-[10px] text-slate-500 font-medium uppercase">{user.displayName || 'Unnamed User'}</span>
                                        </div>
                                    </div>
                                </td>
                                <td class="px-6 py-4">
                                    <span class="font-mono text-[11px] text-indigo-400/80 bg-indigo-500/5 px-2 py-1 rounded border border-indigo-500/10">{user.uid}</span>
                                </td>
                                <td class="px-6 py-4 text-slate-400 font-medium">
                                    {formatDate(user.created)}
                                </td>
                                <td class="px-6 py-4 text-slate-400 font-medium">
                                    {formatDate(user.lastSignIn)}
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <button title="Reset Password" class="p-2 rounded-lg bg-slate-800 text-slate-400 hover:text-white hover:bg-slate-700 transition-all">
                                            <Key size={14} />
                                        </button>
                                        <button 
                                            onclick={() => handleDeleteUser(user.uid)}
                                            title="Delete User" 
                                            class="p-2 rounded-lg bg-red-500/10 text-red-400 hover:bg-red-500 hover:text-white transition-all shadow-lg shadow-red-900/20"
                                        >
                                            <Trash2 size={14} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
                {#if users.length === 0}
                    <div class="p-20 text-center flex flex-col items-center gap-3 opacity-30">
                        <Users size={48} />
                        <p class="text-sm font-black uppercase tracking-widest">No Users Registered</p>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .custom-scrollbar::-webkit-scrollbar { width: 4px; }
    .custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.05); border-radius: 10px; }
</style>
