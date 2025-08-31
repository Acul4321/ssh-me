<script lang="ts">
    import { toastTypes, delete_toast } from './Toast';

    let {
        message,
        type
    }: { message: string; type : toastTypes} = $props();
</script>

<output class="toast" data-type="{type}">
    <div>
        {message}
    </div>
</output>