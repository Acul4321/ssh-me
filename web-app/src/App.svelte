<script lang="ts">
  import { pb } from './lib/pocketbase/Pocketbase'
  import { AuthStore, set_auth_context } from './lib/pocketbase/Auth.svelte'
  import Form from './lib/components/form/UserForm.svelte'
  import FormInput from './lib/components/form/FormInput.svelte';
  import LinksManager from './lib/components/links/LinksManager.svelte';
  import Toaster from './lib/components/toast/Toaster.svelte';
  import { add_toast, toastTypes } from './lib/components/toast/Toast';

  const auth : AuthStore = set_auth_context();

  interface InputInfo {
    label: string;
    field: string;
    type: string;
    value: any;
    placeholder?: string;
  }

  let formInputs: InputInfo[] = $state([
    { label: "Handle",        field: "username",      type: "text",   value: "",        placeholder: "your-handle" },
    { label: "Display Name",  field: "display_name",  type: "text",   value: "",        placeholder: "Alan Turing" },
    { label: "Status",        field: "status",        type: "text",   value: "",        placeholder: "What you're up to" },
    { label: "Bio",           field: "bio",           type: "text",   value: "",        placeholder: "short description of who you are" },
    { label: "Pronouns",      field: "pronouns",      type: "text",   value: "",        placeholder: "they/them" },
    { label: "Location",      field: "location",      type: "text",   value: "",        placeholder: "The Moon, Space" },
    { label: "Colour",        field: "colour",        type: "color",  value: "#000000" },
    { label: "Accent Colour", field: "accent_colour", type: "color",  value: "#000000" },
  ]);

  $effect(() => {
    if (auth.user) {
      formInputs.forEach(input => {
        if (auth.user[input.field] !== undefined) {
          input.value = auth.user[input.field];
        }
      });
    }
  });

  let layout = $state('full');

  $effect(() => {
    if (auth.user?.layout) layout = auth.user.layout;
  });

  let profileState = $derived({
    ...formInputs.reduce((acc, input) => {
      acc[input.field] = input.value;
      return acc;
    }, {} as Record<string, any>),
    layout,
  });

  async function update_user_fields() {
    const profileId = auth.user?.id;
    if (!profileId) return;

    try {
      await pb.collection('profiles').update(profileId, profileState);
      add_toast("Saved", toastTypes.SUCCESS);
    } catch (err: any) {
      const usernameError = err?.data?.data?.username;
      if (usernameError) {
        add_toast("That handle is already taken", toastTypes.ERROR);
      } else {
        console.error("Failed to update profile: ", err);
        add_toast("Failed to save", toastTypes.ERROR);
      }
    }
  }

  let sshLine : string = $derived(auth.user?.username ? `ssh ${auth.user.username}@ssh-me.com` : '');
</script>

<Toaster />

<header>
  <h1>ssh-me</h1>
  <p>SSH profiles for humans</p>
</header>

<main>
  {#if auth.user}
    <p>Signed in as {auth.user.username || auth.user.id}</p>
    <button onclick={auth.logout}>Sign out</button>

    <Form onSubmit={update_user_fields}>
      {#each formInputs as input}
        <FormInput
          label={input.label}
          type={input.type}
          bind:value={input.value}
          placeholder={input.placeholder ?? ''}
        />
      {/each}
      <div>
        <label for="layout">Terminal Layout</label>
        <select id="layout" bind:value={layout}>
          <option value="full">Full</option>
          <option value="compact">Compact</option>
          <option value="minimal">Minimal</option>
        </select>
      </div>
    </Form>

    <LinksManager profileId={auth.user.id} />

    {#if sshLine}
      <div>
        <p>{sshLine}</p>
        <button onclick={() => {
          navigator.clipboard.writeText(sshLine);
          add_toast("Copied", toastTypes.SUCCESS)
        }}>Copy</button>
      </div>
    {/if}

  {:else}
    <button onclick={auth.sign_in_with_google}>Continue with Google</button>
  {/if}
</main>

<footer>
  <p>Share your profile via SSH.</p>
  <p>All info is public.</p>
</footer>
