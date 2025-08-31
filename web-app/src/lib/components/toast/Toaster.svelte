<script lang="ts">
    import Toast from "./Toast.svelte";
    // import dismis function for onclick dismis button
    import { toasts } from "./Toast";
</script>

{#if $toasts}
    <section>
        {#each $toasts as toast (toast.id)}
            <Toast
                message = {toast.message}
                type = {toast.type}
                dismissible = {toast.dismissible}
            />
        {/each}
    </section>
{/if}