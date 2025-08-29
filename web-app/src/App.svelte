<script lang="ts">
  import { pb } from './lib/pocketbase/Pocketbase'
  import { AuthStore, set_auth_context } from './lib/pocketbase/Auth.svelte'
  import Form from './lib/components/form/UserForm.svelte'
  import FormInput from './lib/components/FormInput.svelte';

  let auth : AuthStore = set_auth_context();

  let FormData = $state({
    handle: '',
    status: '',
    colour: '#000000'
  });

  function handleSubmit() {
    console.log("form submit");
  }
  
</script>

<header>
  <h1>ssh-me</h1>
  <p>SSH profiles for "humans"</p>
</header>

<main>
  {#if auth.user}
    <p>Signed in as {auth.user.username}</p>
    <button onclick={auth.logout}>Sign out</button>

    <Form onSubmit={handleSubmit}>
      <FormInput label="Handle" type="text" value={FormData.handle}/>
      <FormInput label="Status" type="text" value={FormData.status}/>
      <FormInput label="Colour" type="color" value={FormData.colour}/>
    </Form>

  {:else} 
    <button onclick={auth.sign_in_with_google}> Continue with Google </button>
  {/if}
</main>

<footer>
  <p>Share your profile via SSH.</p>
  <p>All info is public.</p>
</footer>