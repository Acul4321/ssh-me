<script lang="ts">
  import { pb } from './lib/pocketbase/Pocketbase'
  import { AuthStore, set_auth_context } from './lib/pocketbase/Auth.svelte'
  import Form from './lib/components/form/UserForm.svelte'
  import FormInput from './lib/components/form/FormInput.svelte';
  import Toaster from './lib/components/toast/Toaster.svelte';
    import { add_toast, toastTypes } from './lib/components/toast/Toast';

  const auth : AuthStore = set_auth_context();

  // Form Inputs structure
  interface InputInfo {
    label: string;
    field: string;
    type: string;
    value: any;
  }

  let formInputs: InputInfo[] = $state([
    { label: "Handle", field: "username", type: "text", value: "" },
    { label: "Status", field: "status", type: "text", value: "" },
    { label: "Colour", field: "colour", type: "color", value: "#000000" },
    { label: "LinkedIn", field: "linkedin", type: "url", value: "" }
  ]);

  // for popuating the input fields with user information
  $effect(() => {
    if (auth.user) {
      formInputs.forEach(input => {
        if (auth.user[input.field] !== undefined) {
          input.value = auth.user[input.field];
        }
      });
    }
  });

  // variable storing the profile information in the forms
  let profileState = $derived(formInputs.reduce((acc, input) => {
    acc[input.field] = input.value;
    return acc;
  }, {} as Record<string, any>));

  // for updating the user in the db with profileState TODO: move db functions into seperate file
  async function update_user_fields() {
    console.log(profileState);
    const userId = pb.authStore.record?.id;
    if (!userId) return;

    try {
      await pb.collection('users').update(userId, profileState);
      console.log("profile Successfully updated");
    } catch (err) {
      console.error("Failed to update profile: ", err);
    }
  }
</script>

<Toaster />
<button onclick={() => add_toast("hello World", toastTypes.SUCCESS, 3000)}>toast</button>

<header>
  <h1>ssh-me</h1>
  <p>SSH profiles for "humans"</p>
</header>

<main>
  {#if auth.user}
    <p>Signed in as {auth.user.username}</p>
    <button onclick={auth.logout}>Sign out</button>

    <!-- Profile Form -->
    <Form onSubmit={update_user_fields}>
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