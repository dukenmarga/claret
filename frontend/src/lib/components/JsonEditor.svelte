<script>
    import { JSONEditor } from 'svelte-jsoneditor'
    
    let { value = '', onchange } = $props()
    
    // Convert string to object for the editor
    let content = $derived.by(() => {
        try {
            return { json: JSON.parse(value) }
        } catch (e) {
            return { text: value }
        }
    })

    function handleChange(updatedContent) {
        const newValue = updatedContent.json 
            ? JSON.stringify(updatedContent.json, null, 2) 
            : updatedContent.text
        if (onchange) onchange(newValue)
    }
</script>

<div class="w-full h-full flex flex-col bg-slate-950 overflow-hidden jse-theme-dark">
    <JSONEditor
        {content}
        onChange={handleChange}
        mode="tree"
        mainMenuBar={true}
        navigationBar={true}
        statusBar={true}
    />
</div>

<style>
    /* Dark theme customization for svelte-jsoneditor */
    .jse-theme-dark {
        --jse-theme-color: #6366f1;
        --jse-theme-color-highlight: #818cf8;
        --jse-background-color: #020617;
        --jse-panel-background: #0f172a;
        --jse-font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
        --jse-font-size: 13px;
        --jse-text-color: #d1d5db;
        --jse-text-color-italic: #9ca3af;
        --jse-delimiter-color: #6b7280;
        --jse-key-color: #818cf8;
        --jse-value-color: #10b981;
        --jse-string-color: #10b981;
        --jse-number-color: #f59e0b;
        --jse-boolean-color: #ec4899;
        --jse-null-color: #6b7280;
        
        /* High contrast selection and cursor */
        --jse-selection-background-color: #1e293b; /* Slate-800 - much safer for all colors */
        --jse-selection-inactive-background-color: #0f172a;
        --jse-hover-background-color: #1e293b;
    }

    /* Remove the "box" look and borders from the selected node */
    :global(.jse-focus), :global(.jse-selected) {
        background-color: #1e293b !important; /* Subtle slate background */
        border: none !important;
        outline: 1px solid #4f46e5 !important; /* Sharp indigo outline instead of thick box */
        border-radius: 2px !important;
    }

    /* Ensure text inside the selection is bright and clear */
    :global(.jse-selected .jse-value), 
    :global(.jse-selected .jse-key),
    :global(.jse-focus .jse-value) {
        color: #ffffff !important;
    }

    /* Keep the green for URLs/strings but make them bright */
    :global(.jse-string) {
        color: #34d399 !important; /* Bright Emerald */
    }

    /* Fix the text cursor */
    :global(.jse-editable-text) {
        caret-color: #ffffff !important;
        background: #0f172a !important;
        color: #ffffff !important;
        padding: 2px 4px !important;
    }

    :global(.jse-main) {
        border: none !important;
        background: transparent !important;
    }
</style>
