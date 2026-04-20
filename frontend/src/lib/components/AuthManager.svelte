<script>
    import { onMount } from 'svelte'
    import { GetUsers, DeleteUser, UpdateUser } from '../../../wailsjs/go/main/App'
    import { Loader2, Trash2, Mail, Calendar, Key, UserCheck, Shield, RefreshCw, Users, Edit3, Check, X, ShieldAlert } from 'lucide-svelte'
    import { toast } from 'svelte-sonner'
    import { appState } from '../state.svelte.js'

    let users = $state([])
    let isLoading = $state(false)
    let error = $state(null)

    // Editing State
    let editingUser = $state(null) // User object being edited
    let showEditDialog = $state(false)
    let editForm = $state({
        displayName: '',
        photoURL: '',
        disabled: false
    })

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

    function openEdit(user) {
        editingUser = user
        editForm = {
            displayName: user.displayName || '',
            photoURL: user.photoURL || '',
            disabled: user.disabled || false
        }
        showEditDialog = true
    }

    async function handleUpdateUser() {
        if (!editingUser) return
        try {
            const resp = await UpdateUser(editingUser.uid, editForm)
            if (resp.success) {
                toast.success('User updated')
                showEditDialog = false
                fetchUsers() // Refresh list
            } else {
                toast.error("Update failed: " + resp.error)
            }
        } catch (e) {
            toast.error(e.toString())
        }
    }

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
                            <th class="px-6 py-4 border-b border-slate-800">Status</th>
                            <th class="px-6 py-4 border-b border-slate-800">Signed Up</th>
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
                                <td class="px-6 py-4">
                                    {#if user.disabled}
                                        <span class="text-[9px] font-black uppercase bg-red-500/10 text-red-500 px-2 py-0.5 rounded-full border border-red-500/20">Disabled</span>
                                    {:else}
                                        <span class="text-[9px] font-black uppercase bg-emerald-500/10 text-emerald-500 px-2 py-0.5 rounded-full border border-emerald-500/20">Active</span>
                                    {/if}
                                </td>
                                <td class="px-6 py-4 text-slate-400 font-medium">
                                    {formatDate(user.created)}
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <button onclick={() => openEdit(user)} title="Edit User" class="p-2 rounded-lg bg-slate-800 text-slate-400 hover:text-white hover:bg-slate-700 transition-all">
                                            <Edit3 size={14} />
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

{#if showEditDialog}
    <div class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-[100] p-4">
        <div class="bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200">
            <div class="p-6 border-b border-slate-800 flex items-center justify-between bg-slate-950/20">
                <div class="flex items-center gap-3">
                    <div class="p-2 bg-indigo-500/10 rounded-lg text-indigo-400">
                        <Edit3 size={18} />
                    </div>
                    <h2 class="text-sm font-black text-white uppercase tracking-widest">Edit User</h2>
                </div>
                <button onclick={() => showEditDialog = false} class="text-slate-500 hover:text-white transition-colors">
                    <X size={20} />
                </button>
            </div>

            <div class="p-8 space-y-6">
                <div>
                    <label class="block text-[10px] font-black text-slate-500 uppercase tracking-widest mb-2">Display Name</label>
                    <input 
                        bind:value={editForm.displayName} 
                        type="text" 
                        class="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all" 
                        placeholder="John Doe"
                    />
                </div>

                <div>
                    <label class="block text-[10px] font-black text-slate-500 uppercase tracking-widest mb-2">Photo URL</label>
                    <input 
                        bind:value={editForm.photoURL} 
                        type="text" 
                        class="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all font-mono" 
                        placeholder="https://..."
                    />
                </div>

                <div class="flex items-center justify-between p-4 bg-slate-950 rounded-xl border border-slate-800">
                    <div>
                        <p class="text-xs font-bold text-white">Disable Account</p>
                        <p class="text-[10px] text-slate-500">Prevent user from signing in</p>
                    </div>
                    <button 
                        onclick={() => editForm.disabled = !editForm.disabled}
                        class="w-10 h-5 rounded-full transition-all relative {editForm.disabled ? 'bg-red-600 shadow-[0_0_12px_rgba(239,68,68,0.3)]' : 'bg-slate-800'}"
                    >
                        <div class="absolute top-1 left-1 w-3 h-3 rounded-full bg-white transition-all {editForm.disabled ? 'translate-x-5' : ''}"></div>
                    </button>
                </div>
            </div>

            <div class="p-6 bg-slate-950/40 border-t border-slate-800 flex gap-4">
                <button onclick={() => showEditDialog = false} class="flex-1 py-2.5 rounded-xl text-xs font-black uppercase tracking-widest text-slate-500 hover:text-white hover:bg-slate-800 transition-all">Cancel</button>
                <button onclick={handleUpdateUser} class="flex-1 py-2.5 rounded-xl text-xs font-black uppercase tracking-widest bg-indigo-600 text-white shadow-lg shadow-indigo-900/40 hover:bg-indigo-500 transition-all active:scale-95">Save Changes</button>
            </div>
        </div>
    </div>
{/if}

<style>
    .custom-scrollbar::-webkit-scrollbar { width: 4px; }
    .custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.05); border-radius: 10px; }
</style>
