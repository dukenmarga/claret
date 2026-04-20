<script>
    import { untrack } from 'svelte'
    import { GetUsers, DeleteUser, UpdateUser, SearchUsers } from '../../../wailsjs/go/main/App'
    import { Loader2, Trash2, Mail, Shield, RefreshCw, Users, Edit3, X, ShieldAlert, Search, Phone, User, Globe } from 'lucide-svelte'
    import { toast } from 'svelte-sonner'
    import { appState } from '../state.svelte.js'

    let users = $state([])
    let isLoading = $state(false)
    let error = $state(null)
    let searchQuery = $state('')

    async function fetchUsers() {
        if (!appState.currentProjectId) return
        
        isLoading = true
        error = null
        
        try {
            const query = searchQuery.trim()
            const resp = query 
                ? await SearchUsers(query)
                : await GetUsers(100)
                
            if (resp.success) {
                users = resp.data || []
            } else {
                error = resp.error
                users = []
                toast.error("Auth fetch failed: " + resp.error)
            }
        } catch (e) {
            error = e.toString()
            users = []
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

    // Editing State
    let editingUser = $state(null)
    let showEditDialog = $state(false)
    let editForm = $state({
        displayName: '',
        photoURL: '',
        disabled: false
    })

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
                fetchUsers()
            } else {
                toast.error("Update failed: " + resp.error)
            }
        } catch (e) {
            toast.error(e.toString())
        }
    }

    $effect(() => {
        if (appState.currentProjectId) {
            untrack(() => {
                fetchUsers()
            })
        }
    })

    function formatDate(ts) {
        if (!ts) return 'Never'
        return new Date(ts).toLocaleString()
    }

    function getIdentifier(user) {
        if (user.phoneNumber) return user.phoneNumber
        if (user.email) return user.email
        return '(anonymous)'
    }
</script>

<div class="flex flex-col h-full space-y-4 font-sans">
    <!-- Top Bar -->
    <div class="flex items-center justify-between bg-slate-900/40 p-4 rounded-xl border border-slate-800 shadow-sm shrink-0">
        <div class="flex items-center gap-3">
            <div class="p-2 bg-indigo-500/10 rounded-lg text-indigo-400">
                <Shield size={24} />
            </div>
            <div>
                <h2 class="text-sm font-black text-white uppercase tracking-widest">Authentication</h2>
                <p class="text-[10px] text-slate-500 font-medium uppercase tracking-tight">Manage users for {appState.currentProject?.name}</p>
            </div>
        </div>

        <div class="flex items-center gap-4">
            <div class="relative flex items-center bg-slate-950 border border-slate-800 rounded-lg px-3 py-1.5 focus-within:ring-1 focus-within:ring-indigo-500 transition-all shadow-inner">
                <Search size={14} class="text-slate-500 mr-2" />
                <input 
                    bind:value={searchQuery}
                    onkeydown={(e) => {
                        if (e.key === 'Enter') {
                            e.preventDefault();
                            fetchUsers();
                        }
                    }}
                    oninput={() => {
                        if (searchQuery.trim().length === 0) {
                            fetchUsers();
                        }
                    }}
                    placeholder="Search UID, Email, or Phone..."
                    class="bg-transparent border-none text-xs text-white focus:ring-0 w-64 selection:bg-white selection:text-black placeholder:text-slate-700"
                />
                {#if searchQuery}
                    <button onclick={() => { searchQuery = ''; fetchUsers(); }} class="ml-2 text-slate-500 hover:text-white transition-colors">
                        <X size={14} />
                    </button>
                {/if}
            </div>

            <button 
                onclick={fetchUsers} 
                disabled={isLoading}
                class="p-2 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-white transition-all shadow-sm"
            >
                <RefreshCw size={18} class={isLoading ? 'animate-spin text-indigo-400' : ''} />
            </button>
        </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 relative border border-slate-800 rounded-xl overflow-hidden bg-slate-900/20 shadow-inner mx-1">
        {#if isLoading}
            <div class="absolute inset-0 flex flex-col items-center justify-center text-slate-500 gap-4 bg-slate-900/40 z-30 backdrop-blur-[1px]">
                <Loader2 size={48} class="animate-spin text-indigo-500" />
                <span class="text-xs font-black uppercase tracking-widest animate-pulse">Syncing User Registry...</span>
            </div>
        {/if}

        {#if error && !isLoading}
            <div class="h-full flex flex-col items-center justify-center p-8 text-center gap-4">
                <ShieldAlert size={48} class="text-red-500" />
                <div>
                    <p class="text-red-400 text-sm font-bold mb-1 uppercase tracking-widest">Connection Failure</p>
                    <p class="text-red-400/60 text-xs font-mono max-w-md break-words">{error}</p>
                </div>
                <button onclick={fetchUsers} class="px-6 py-2 bg-red-500/20 text-red-400 rounded-lg text-[10px] font-black uppercase tracking-widest hover:bg-red-500/30 transition-all">Retry Connection</button>
            </div>
        {:else if users.length === 0 && !isLoading}
            <div class="h-full flex flex-col items-center justify-center text-slate-500 gap-4 opacity-50 m-4">
                <ShieldAlert size={64} strokeWidth={1.5} />
                <div class="text-center">
                    <p class="text-lg font-black uppercase tracking-widest text-slate-200">No Users Found</p>
                    <p class="text-xs font-medium mt-1 uppercase text-slate-500 tracking-tight">Check your search query or project configuration</p>
                </div>
                <button onclick={() => { searchQuery = ''; fetchUsers(); }} class="mt-4 px-6 py-2 rounded-lg bg-slate-800 text-white text-[10px] font-black uppercase hover:bg-slate-700 transition-all active:scale-95">Reset Search</button>
            </div>
        {:else if users.length > 0}
            <div class="overflow-auto h-full custom-scrollbar">
                <table class="w-full text-left text-xs border-collapse table-fixed">
                    <thead class="sticky top-0 bg-slate-900 z-10 shadow-md">
                        <tr class="text-[10px] uppercase font-black tracking-[0.1em] text-slate-500 bg-slate-900">
                            <th class="px-6 py-4 border-b border-slate-800 w-[250px]">Identifier</th>
                            <th class="px-6 py-4 border-b border-slate-800 w-[120px]">Providers</th>
                            <th class="px-6 py-4 border-b border-slate-800 w-[200px]">User ID</th>
                            <th class="px-6 py-4 border-b border-slate-800 w-[100px]">Status</th>
                            <th class="px-6 py-4 border-b border-slate-800 w-[180px]">Created</th>
                            <th class="px-6 py-4 border-b border-slate-800 w-[180px]">Signed In</th>
                            <th class="px-6 py-4 border-b border-slate-800 text-right w-[100px]">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-800/50">
                        {#each users as user (user.uid)}
                            <tr class="hover:bg-white/[0.03] transition-colors group contain-content">
                                <td class="px-6 py-4 truncate">
                                    <div class="flex items-center gap-4 overflow-hidden">
                                        <div class="w-10 h-10 rounded-xl bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-black overflow-hidden shadow-lg flex-shrink-0 ring-1 ring-white/5">
                                            {#if user.photoURL}
                                                <img src={user.photoURL} alt="" class="w-full h-full object-cover" />
                                            {:else}
                                                <span class="text-base">{user.email ? user.email[0].toUpperCase() : (user.phoneNumber ? '#' : '?')}</span>
                                            {/if}
                                        </div>
                                        <div class="flex flex-col gap-0.5 overflow-hidden">
                                            <span class="text-white font-black text-sm tracking-tight truncate">{getIdentifier(user)}</span>
                                            <span class="text-[10px] text-slate-500 font-bold uppercase tracking-tighter truncate">{user.displayName || 'Unnamed User'}</span>
                                        </div>
                                    </div>
                                </td>
                                <td class="px-6 py-4">
                                    <div class="flex items-center gap-2">
                                        {#if user.providers && user.providers.length > 0}
                                            {#each user.providers as provider}
                                                <div class="p-2 rounded-xl bg-slate-950/80 border border-slate-800 flex items-center justify-center shadow-lg ring-1 ring-white/5" title={provider}>
                                                    {#if provider === 'google.com'}
                                                        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
                                                            <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                                                            <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                                                            <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                                                            <path d="M12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
                                                        </svg>
                                                    {:else if provider === 'password'}<Mail size={16} class="text-blue-400" />
                                                    {:else if provider === 'phone'}<Phone size={16} class="text-emerald-400" />
                                                    {:else}<User size={16} class="text-slate-400" />{/if}
                                                </div>
                                            {/each}
                                        {:else}
                                            <div class="p-2 rounded-xl bg-slate-950/80 border border-slate-800 flex items-center justify-center opacity-30 shadow-lg ring-1 ring-white/5" title="Anonymous">
                                                <User size={16} />
                                            </div>
                                        {/if}
                                    </div>
                                </td>
                                <td class="px-6 py-4 truncate font-mono text-[11px]">
                                    <span class="text-indigo-300 bg-indigo-500/10 px-2.5 py-1 rounded-lg border border-indigo-500/20 shadow-sm">{user.uid}</span>
                                </td>
                                <td class="px-6 py-4">
                                    {#if user.disabled}
                                        <span class="text-[9px] font-black uppercase bg-red-500/10 text-red-500 px-2.5 py-1 rounded-full border border-red-500/20">Disabled</span>
                                    {:else}
                                        <span class="text-[9px] font-black uppercase bg-emerald-500/10 text-emerald-500 px-2.5 py-1 rounded-full border border-emerald-500/20">Active</span>
                                    {/if}
                                </td>
                                <td class="px-6 py-4 text-slate-400 font-bold tracking-tight truncate">
                                    {formatDate(user.created)}
                                </td>
                                <td class="px-6 py-4 text-slate-200 font-black tracking-tight truncate">
                                    {formatDate(user.lastSignIn)}
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-all transform translate-x-2 group-hover:translate-x-0">
                                        <button onclick={() => openEdit(user)} title="Edit User" class="p-2.5 rounded-xl bg-slate-800 text-slate-300 hover:text-white hover:bg-slate-700 transition-all border border-slate-700/50 shadow-md">
                                            <Edit3 size={16} />
                                        </button>
                                        <button 
                                            onclick={() => handleDeleteUser(user.uid)}
                                            title="Delete User" 
                                            class="p-2.5 rounded-xl bg-red-500/10 text-red-400 hover:bg-red-500 hover:text-white transition-all border border-red-500/20 shadow-lg shadow-red-950/30"
                                        >
                                            <Trash2 size={16} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>

{#if showEditDialog}
    <div class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-[100] p-4">
        <div class="bg-slate-900 border border-slate-800 rounded-3xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200 ring-1 ring-white/10">
            <div class="p-6 border-b border-slate-800 flex items-center justify-between bg-slate-950/20">
                <div class="flex items-center gap-3">
                    <div class="p-2.5 bg-indigo-500/10 rounded-xl text-indigo-400">
                        <Edit3 size={20} />
                    </div>
                    <h2 class="text-sm font-black text-white uppercase tracking-widest">Modify User</h2>
                </div>
                <button onclick={() => showEditDialog = false} class="p-2 rounded-lg hover:bg-slate-800 text-slate-500 hover:text-white transition-colors">
                    <X size={20} />
                </button>
            </div>

            <div class="p-8 space-y-8">
                <div>
                    <label class="block text-[10px] font-black text-slate-500 uppercase tracking-widest mb-3 ml-1">Profile Name</label>
                    <input bind:value={editForm.displayName} type="text" class="w-full bg-slate-950 border border-slate-800 rounded-2xl px-5 py-3 text-sm text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all shadow-inner" placeholder="e.g. Alex Rivera" />
                </div>
                <div>
                    <label class="block text-[10px] font-black text-slate-500 uppercase tracking-widest mb-3 ml-1">Avatar Resource (URL)</label>
                    <input bind:value={editForm.photoURL} type="text" class="w-full bg-slate-950 border border-slate-800 rounded-2xl px-5 py-3 text-sm text-white focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all font-mono shadow-inner" placeholder="https://images.com/user.jpg" />
                </div>
                <div class="flex items-center justify-between p-5 bg-slate-950 rounded-2xl border border-slate-800 shadow-inner">
                    <div>
                        <p class="text-xs font-black text-white uppercase tracking-tight">Account Access</p>
                        <p class="text-[10px] text-slate-500 font-bold uppercase mt-0.5">Toggle Login Authorization</p>
                    </div>
                    <button onclick={() => editForm.disabled = !editForm.disabled} class="w-12 h-6 rounded-full transition-all relative {editForm.disabled ? 'bg-red-600 shadow-[0_0_15px_rgba(239,68,68,0.4)]' : 'bg-slate-800 shadow-inner'}">
                        <div class="absolute top-1 left-1 w-4 h-4 rounded-full bg-white transition-all shadow-sm {editForm.disabled ? 'translate-x-6' : ''}"></div>
                    </button>
                </div>
            </div>

            <div class="p-6 bg-slate-950/40 border-t border-slate-800 flex gap-4">
                <button onclick={() => showEditDialog = false} class="flex-1 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest text-slate-500 hover:text-white hover:bg-slate-800 transition-all">Discard</button>
                <button onclick={handleUpdateUser} class="flex-1 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest bg-indigo-600 text-white shadow-xl shadow-indigo-900/40 hover:bg-indigo-500 transition-all active:scale-95 ring-1 ring-white/10">Apply Changes</button>
            </div>
        </div>
    </div>
{/if}

<style>
    .custom-scrollbar::-webkit-scrollbar { width: 4px; }
    .custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
    .custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.05); border-radius: 10px; }
    
    table { table-layout: fixed; width: 100%; border-spacing: 0; }
    tr { contain: content; }
</style>
