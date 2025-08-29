<script lang="ts">
    let {
        label,
        type = 'text',
        value = '',
    }: { label: string; type?: string; value?: string } = $props();
</script>

<div>
    <label for="{value}">{label}</label>
    <input {type} name="{value}" bind:value />
</div>