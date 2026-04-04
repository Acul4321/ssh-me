<script lang="ts">
    import Toast from "./Toast.svelte";
    import { toasts } from "./Toast";
    import "./Toaster.css";
</script>

{#if $toasts}
    <section class="toast-stack" aria-live="polite" aria-label="Notifications">
        {#each $toasts as toast (toast.id)}
            <Toast
                message = {toast.message}
                type = {toast.type}
            />
        {/each}
    </section>
{/if}
