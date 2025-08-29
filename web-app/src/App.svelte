<script lang="ts">
  import { pb } from './lib/Pocketbase'
  import { AuthStore, set_auth_context } from './lib/Auth.svelte'

  let auth : AuthStore = set_auth_context();

</script>

<main>
  <h1>ssh-me</h1>
  <p>SSH profiles for "humans"</p>

  {#if auth.user}
    <p>Signed in as {auth.user.username}</p>
    <button onclick={auth.logout}>Sign out</button>
  {:else} 
    <button onclick={auth.sign_in_with_google}> Continue with Google </button>
  {/if}

  <p>Share your profile via SSH.</p>
  <p>All info is public.</p>
</main>