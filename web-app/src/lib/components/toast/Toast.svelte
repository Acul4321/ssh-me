<script lang="ts">
    import { toastTypes } from './Toast';

    let {
        message,
        type,
        dismissible = true,
    }: { message: string; type : toastTypes; dismissible? : boolean } = $props();
</script>

<output class="toast" data-type="{type}">
    <div>
        {message}
    </div>

    {#if dismissible}
    <button class="close" onclick={() => console.log("delete toast")}>
        close
    </button>
    {/if}
</output>