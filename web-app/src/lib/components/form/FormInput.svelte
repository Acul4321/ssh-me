<script lang="ts">
    let {
        label,
        type = 'text',
        value = $bindable(''),
    }: { label: string; type?: string; value?: string } = $props();
</script>

<div>
    <label for="{value}">{label}</label>
    <input {type} name="{value}" bind:value={value} />
</div>