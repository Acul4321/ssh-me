<script lang="ts">
    let {
        label,
        type = 'text',
        value = $bindable(''),
        placeholder = '',
    }: { label: string; type?: string; value?: string; placeholder?: string } = $props();
</script>

<div>
    <label for="{value}">{label}</label>
    <input {type} name="{value}" bind:value={value} {placeholder} />
</div>