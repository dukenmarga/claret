import { GetConfig, ConnectProject, GetCollections, GetDocuments, UpdateDocument, DeleteDocuments } from '../../wailsjs/go/main/App'
import { toast } from 'svelte-sonner'

class AppState {
    projects = $state([])
    currentProjectId = $state(null)
    
    // Explorer / Tree State
    collections = $state([]) // Root collections
    docSubCollections = $state({}) // docPath -> string[] (sub-collection names)
    collectionDocuments = $state({}) // collectionPath -> Document[]
    
    expandedDocs = $state([]) // docPath[]
    expandedCollections = $state([]) // collectionPath[]
    
    // Global Selection/Path (Main focus)
    selectedPath = $state('')
    documents = $state([]) // This remains the "main" view list
    selectedDoc = $state(null)
    selectedDocIds = $state([])
    
    // UI State
    hasMore = $state(false)
    isLoading = $state(false)
    error = $state(null)
    
    // Query State
    filterField = $state('id')
    filterValue = $state('')
    activeQueries = $state([])
    orderByField = $state('')
    orderDescending = $state(false)

    confirmDeletionList = $state([])
    showDeleteDialog = $state(false)

    constructor() {
        this.refreshConfig()
    }

    get currentProject() {
        return this.projects.find(p => p.id === this.currentProjectId)
    }

    refreshConfig = async () => {
        try {
            const config = await GetConfig()
            if (config?.projects) this.projects = config.projects
        } catch (e) { console.error(e) }
    }

    selectProject = async (id, databaseId = '') => {
        this.isLoading = true
        this.error = null
        try {
            const resp = await ConnectProject(id, databaseId)
            if (resp.success) {
                this.currentProjectId = id
                this.selectedPath = ''
                this.documents = []
                this.selectedDoc = null
                this.docSubCollections = {}
                this.collectionDocuments = {}
                this.expandedDocs = []
                this.expandedCollections = []
                await this.refreshCollections('')
            } else toast.error(resp.error)
        } catch (e) { toast.error(e.toString()) }
        finally { this.isLoading = false }
    }

    refreshCollections = async (parentPath) => {
        this.isLoading = true
        try {
            const resp = await GetCollections(parentPath)
            if (resp.success) {
                if (parentPath === '') {
                    this.collections = resp.data || []
                } else {
                    this.docSubCollections[parentPath] = resp.data || []
                }
            } else toast.error(resp.error)
        } catch (e) { toast.error(e.toString()) }
        finally { this.isLoading = false }
    }

    selectCollection = async (path) => {
        this.selectedPath = path
        this.documents = []
        this.selectedDoc = null
        this.selectedDocIds = []
        this.activeQueries = []
        this.orderByField = ''
        this.orderDescending = false
        this.docSubCollections = {}
        this.collectionDocuments = {}
        this.expandedDocs = []
        this.expandedCollections = []
        await this.executeQuery(false)
    }

    async executeQuery(isLoadMore = false) {
        if (!this.selectedPath) return
        this.isLoading = true
        if (!isLoadMore) this.filterValue = '' 

        const processedQueries = this.activeQueries
            .filter(q => q.value !== undefined && q.value !== null && q.value !== '')
            .map(q => {
                let val = q.value
                if (typeof val === 'string') {
                    const lower = val.toLowerCase().trim()
                    if (lower === 'true') val = true
                    else if (lower === 'false') val = false
                    else if (lower === 'null') val = null
                    else if (!isNaN(Number(val)) && val !== '') val = Number(val)
                }
                return { ...q, value: val }
            })

        try {
            const limit = 100
            const resp = await GetDocuments(this.selectedPath, {
                conditions: processedQueries,
                orderBy: this.orderByField,
                descending: this.orderDescending,
                limit: limit
            })
            if (resp.success) {
                const newDocs = resp.data || []
                if (isLoadMore) this.documents = [...this.documents, ...newDocs]
                else this.documents = newDocs
                this.hasMore = newDocs.length === limit
            } else toast.error(resp.error)
        } catch (e) { toast.error(e.toString()) }
        finally { this.isLoading = false }
    }

    toggleDocExpansion = async (docPath) => {
        if (this.expandedDocs.includes(docPath)) {
            this.expandedDocs = this.expandedDocs.filter(p => p !== docPath)
        } else {
            this.expandedDocs = [...this.expandedDocs, docPath]
            if (!this.docSubCollections[docPath]) {
                const resp = await GetCollections(docPath)
                if (resp.success) {
                    this.docSubCollections[docPath] = resp.data || []
                }
            }
        }
    }

    toggleCollectionExpansion = async (collPath) => {
        if (this.expandedCollections.includes(collPath)) {
            this.expandedCollections = this.expandedCollections.filter(p => p !== collPath)
        } else {
            this.expandedCollections = [...this.expandedCollections, collPath]
            if (!this.collectionDocuments[collPath]) {
                const resp = await GetDocuments(collPath, { conditions: [], limit: 50 })
                if (resp.success) {
                    this.collectionDocuments[collPath] = resp.data || []
                }
            }
        }
    }

    toggleSelection(id) {
        if (this.selectedDocIds.includes(id)) this.selectedDocIds = this.selectedDocIds.filter(i => i !== id)
        else this.selectedDocIds = [...this.selectedDocIds, id]
    }

    selectAll(select) {
        this.selectedDocIds = select ? this.documents.map(doc => doc.id) : []
    }

    requestDelete = (docs) => {
        this.confirmDeletionList = docs
        this.showDeleteDialog = true
    }

    async confirmDelete() {
        const paths = this.confirmDeletionList.map(d => d.path)
        const ids = this.confirmDeletionList.map(d => d.id)
        this.isLoading = true
        try {
            const resp = await DeleteDocuments(paths)
            if (resp.success) {
                this.documents = this.documents.filter(doc => !ids.includes(doc.id))
                this.selectedDocIds = this.selectedDocIds.filter(id => !ids.includes(id))
                if (this.selectedDoc && ids.includes(this.selectedDoc.id)) this.selectedDoc = null
                toast.success(`Deleted ${paths.length} docs`)
            } else toast.error(resp.error)
        } catch (e) { toast.error(e.toString()) }
        finally {
            this.showDeleteDialog = false
            this.confirmDeletionList = []
            this.isLoading = false
        }
    }

    deleteSelected = () => {
        const docs = this.documents.filter(doc => this.selectedDocIds.includes(doc.id)).map(doc => ({ id: doc.id, path: doc.path }))
        if (docs.length > 0) this.requestDelete(docs)
    }

    deleteDocument = (path) => {
        const doc = this.documents.find(d => d.path === path) || { id: path.split('/').pop(), path }
        this.requestDelete([{ id: doc.id, path: doc.path }])
    }

    async duplicateTo(targetCollection, doc = null) {
        const docs = doc ? [doc] : this.documents.filter(d => this.selectedDocIds.includes(d.id))
        if (docs.length === 0) return
        let count = 0
        for (const d of docs) {
            try {
                const resp = await UpdateDocument(`${targetCollection}/${d.id}`, d.data)
                if (resp.success) count++
            } catch (e) {}
        }
        if (count > 0) {
            toast.success(`Duplicated ${count} docs`)
            if (this.selectedPath === targetCollection) await this.selectCollection(targetCollection)
        }
    }

    async saveDocument(doc) {
        try {
            const resp = await UpdateDocument(doc.path, doc.data)
            if (resp.success) {
                const idx = this.documents.findIndex(d => d.path === doc.path)
                if (idx !== -1) this.documents[idx] = doc
                else this.documents.unshift(doc)
                toast.success('Saved')
            } else toast.error(resp.error)
        } catch (e) { toast.error(e.toString()) }
    }

    addDocument = (collectionPath) => {
        const id = Math.random().toString(36).substring(2, 12)
        this.selectedDoc = { id, path: `${collectionPath}/${id}`, data: { createdAt: new Date().toISOString() } }
    }

    async renameDocument(doc, newId) {
        if (!newId || newId === doc.id) return
        const oldPath = doc.path
        const newPath = `${oldPath.substring(0, oldPath.lastIndexOf('/'))}/${newId}`
        this.isLoading = true
        try {
            if ((await UpdateDocument(newPath, doc.data)).success) {
                if ((await DeleteDocuments([oldPath])).success) {
                    const idx = this.documents.findIndex(d => d.path === oldPath)
                    if (idx !== -1) this.documents[idx] = { ...doc, id: newId, path: newPath }
                    if (this.selectedDoc?.path === oldPath) this.selectedDoc = { ...doc, id: newId, path: newPath }
                    toast.success('Renamed')
                }
            }
        } catch (e) { toast.error(e.toString()) }
        finally { this.isLoading = false }
    }

    addQueryFilter = () => { this.activeQueries.push({ field: 'id', operator: '==', value: '' }) }
    removeQueryFilter = (i) => { this.activeQueries = this.activeQueries.filter((_, idx) => idx !== i) }
    async setSort(field) {
        if (this.orderByField === field) this.orderDescending = !this.orderDescending
        else { this.orderByField = field; this.orderDescending = false }
        await this.executeQuery()
    }
}

export const appState = new AppState()
