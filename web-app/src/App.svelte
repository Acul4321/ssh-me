<script lang="ts">
  import { pb } from './lib/pocketbase/Pocketbase'
  import { AuthStore, set_auth_context } from './lib/pocketbase/Auth.svelte'
  import Form from './lib/components/form/UserForm.svelte'
  import FormInput from './lib/components/form/FormInput.svelte';

  const auth : AuthStore = set_auth_context();

  interface InputInfo {
    label: string;
    type: string;
    value: any;
  }

  let formInputs: InputInfo[] = $state([
    { label: "Handle", type: "text", value: "" },
    { label: "Status", type: "text", value: "" },
    { label: "Colour", type: "color", value: "#000000" }
  ]);

  let profileState = $derived(formInputs.reduce((acc, input) => {
    acc[input.label.toLowerCase()] = input.value;
    return acc;
  }, {} as Record<string, any>));

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

    <!-- Profile Form -->
    <Form onSubmit={handleSubmit}>
      {#each formInputs as input}
        <FormInput
          label={input.label}
          type={input.type}
          bind:value={input.value}
        />
      {/each}
    </Form>

  {:else} 
    <button onclick={auth.sign_in_with_google}> Continue with Google </button>
  {/if}
</main>

<footer>
  <p>Share your profile via SSH.</p>
  <p>All info is public.</p>
</footer>