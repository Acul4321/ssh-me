<script lang="ts">
    export let onSubmit: (event: Event) => void;
</script>    

<form on:submit|preventDefault={onSubmit}>
    <slot />
    <button type="submit">Submit</button>
</form>