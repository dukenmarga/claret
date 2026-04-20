<script>
    import { onMount } from 'svelte'
    import { GetUsers, DeleteUser } from '../../../wailsjs/go/main/App'
    import { Loader2, Trash2, Mail, Calendar, Key } from 'lucide-svelte'
    import { toast } from 'svelte-sonner'

    let users = $state([])
    let isLoading = $state(false)
    let error = $state(null)

    async function fetchUsers() {
        isLoading = true
        try {
            const resp = await GetUsers(100)
            if (resp.success) {
                users = resp.data || []
            } else {
                error = resp.error
            }
        } catch (e) {
            error = e.toString()
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
                error = resp.error
                toast.error(`Delete failed: ${resp.error}`)
            }
        } catch (e) {
            error = e.toString()
            toast.error(e.toString())
        }
    }

    onMount(fetchUsers)
</script>

<div class="flex flex-col h-full space-y-4">
    <div class="flex items-center justify-between">
        <h2 class="text-xl font-bold text-white">Auth Users</h2>
        <button onclick={fetchUsers} class="text-slate-400 hover:text-white transition-colors">
            <Loader2 size={18} class={isLoading ? 'animate-spin' : ''} />
        </button>
    </div>

    {#if isLoading && users.length === 0}
        <div class="flex-1 flex items-center justify-center text-slate-500">
            <Loader2 size={32} class="animate-spin text-indigo-500 mb-2" />
            <span>Loading users...</span>
        </div>
    {:else if error}
        <div class="p-4 bg-red-500/10 border border-red-500/20 text-red-400 rounded-lg">
            {error}
        </div>
    {:else}
        <div class="flex-1 border border-slate-800 rounded-lg overflow-auto bg-slate-900/20">
            <table class="w-full text-left text-sm">
                <thead class="sticky top-0 bg-slate-900 z-10">
                    <tr class="text-[11px] uppercase tracking-wider text-slate-500">
                        <th class="px-6 py-3 border-b border-slate-800">User</th>
                        <th class="px-6 py-3 border-b border-slate-800">UID</th>
                        <th class="px-6 py-3 border-b border-slate-800">Created</th>
                        <th class="px-6 py-3 border-b border-slate-800">Last Sign In</th>
                        <th class="px-6 py-3 border-b border-slate-800">Actions</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-800/50">
                    {#each users as user}
                        <tr class="hover:bg-slate-800/30 transition-colors">
                            <td class="px-6 py-4">
                                <div class="flex flex-col">
                                    <span class="text-slate-200 font-medium">{user.email || 'No Email'}</span>
                                    <span class="text-xs text-slate-500">{user.displayName || ''}</span>
                                </div>
                            </td>
                            <td class="px-6 py-4 font-mono text-xs text-slate-400">{user.uid}</td>
                            <td class="px-6 py-4 text-xs text-slate-400">
                                {new Date(user.userMetadata.creationTimestamp / 1000).toLocaleDateString()}
                            </td>
                            <td class="px-6 py-4 text-xs text-slate-400">
                                {user.userMetadata.lastSignInTimestamp ? new Date(user.userMetadata.lastSignInTimestamp / 1000).toLocaleDateString() : 'Never'}
                            </td>
                            <td class="px-6 py-4">
                                <div class="flex items-center gap-2">
                                    <button title="Reset Password" class="p-1.5 rounded bg-slate-800 text-slate-400 hover:text-white transition-colors">
                                        <Key size={14} />
                                    </button>
                                    <button 
                                        onclick={() => handleDeleteUser(user.uid)}
                                        title="Delete User" 
                                        class="p-1.5 rounded bg-slate-800 text-red-400 hover:bg-red-500/20 transition-colors"
                                    >
                                        <Trash2 size={14} />
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
