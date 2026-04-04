<script lang="ts">
  import { pb } from './lib/pocketbase/Pocketbase';
  import { AuthStore, set_auth_context } from './lib/pocketbase/Auth.svelte';
  import LinksManager from './lib/components/links/LinksManager.svelte';
  import Toaster from './lib/components/toast/Toaster.svelte';
  import { add_toast, toastTypes } from './lib/components/toast/Toast';

  const auth: AuthStore = set_auth_context();

  interface InputInfo {
    label: string;
    field: string;
    type: string;
    value: string;
    placeholder?: string;
  }

  let formInputs: InputInfo[] = $state([
    { label: 'Handle', field: 'username', type: 'text', value: '', placeholder: 'your-handle' },
    { label: 'Display Name', field: 'display_name', type: 'text', value: '', placeholder: 'Alan Turing' },
    { label: 'Status', field: 'status', type: 'text', value: '', placeholder: "What you're up to" },
    { label: 'Bio', field: 'bio', type: 'textarea', value: '', placeholder: 'Short description' },
    { label: 'Pronouns', field: 'pronouns', type: 'text', value: '', placeholder: 'they/them' },
    { label: 'Location', field: 'location', type: 'text', value: '', placeholder: 'London' },
    { label: 'Colour', field: 'colour', type: 'color', value: '#ffffff' },
    { label: 'Accent Colour', field: 'accent_colour', type: 'color', value: '#cccccc' },
  ]);

  let layout = $state('full');
  let border = $state('thick');

  $effect(() => {
    if (!auth.user) return;

    formInputs.forEach((input) => {
      const nextValue = auth.user?.[input.field];
      if (nextValue !== undefined && nextValue !== null && nextValue !== '') {
        input.value = String(nextValue);
      }
    });

    if (auth.user.layout) layout = auth.user.layout;
    if (auth.user.border) border = auth.user.border;
  });

  let profileState = $derived({
    ...formInputs.reduce((acc, input) => {
      acc[input.field] = input.value;
      return acc;
    }, {} as Record<string, string>),
    layout,
    border,
  });

  let sshLine = $derived(auth.user?.username ? `ssh ${auth.user.username}@ssh-me.com` : 'ssh your-handle@ssh-me.com');
  const demoSshLine = 'ssh demo@ssh-me.com';

  async function updateUserFields() {
    const profileId = auth.user?.id;
    if (!profileId) return;

    try {
      const updatedProfile = await pb.collection('profiles').update(profileId, profileState);
      auth.user = updatedProfile;
      add_toast('Saved profile', toastTypes.SUCCESS);
    } catch (err: any) {
      const usernameError = err?.data?.data?.username;
      if (usernameError) {
        add_toast('That handle is already taken', toastTypes.ERROR);
        return;
      }

      console.error('Failed to update profile:', err);
      add_toast('Failed to save profile', toastTypes.ERROR);
    }
  }

  async function copyCommand(command: string) {
    await navigator.clipboard.writeText(command);
    add_toast('Copied SSH command', toastTypes.SUCCESS);
  }
</script>

<Toaster />

<main class="page-shell">
  <section class="shell-card">
    {#if auth.user}
      <h1>profile</h1>
      <p>Signed in as {auth.user.username || auth.user.id}. Edit your profile and links.</p>

      <form class="profile-form" onsubmit={(event) => {
        event.preventDefault();
        updateUserFields();
      }}>
        <div class="form-grid">
          {#each formInputs as input}
            <label class:wide={input.field === 'status' || input.field === 'bio'}>
              {input.label}

              {#if input.type === 'textarea'}
                <textarea
                  bind:value={input.value}
                  placeholder={input.placeholder ?? ''}
                ></textarea>
              {:else if input.type === 'color'}
                <div class="color-field">
                  <span class="color-swatch" style:background={input.value}>
                    <input
                      type="color"
                      bind:value={input.value}
                      aria-label={input.label}
                    >
                  </span>
                </div>
              {:else}
                <input
                  type={input.type}
                  bind:value={input.value}
                  placeholder={input.placeholder ?? ''}
                >
              {/if}
            </label>
          {/each}

          <label>
            Layout
            <select bind:value={layout}>
              <option value="full">full</option>
              <option value="compact">compact</option>
              <option value="minimal">minimal</option>
            </select>
          </label>

          <label>
            Border
            <select bind:value={border}>
              <option value="thick">thick</option>
              <option value="rounded">rounded</option>
              <option value="normal">normal</option>
              <option value="double">double</option>
              <option value="ascii">ascii</option>
              <option value="hash">hash</option>
              <option value="block">block</option>
              <option value="star">star</option>
            </select>
          </label>
        </div>

        <button type="submit" class="primary-action">Save profile</button>
      </form>

      <LinksManager profileId={auth.user.id} />

      <footer class="footer-row">
        <div class="copy-command-wrap">
          <button type="button" class="copy-command" onclick={() => copyCommand(sshLine)}>
            {sshLine}
          </button>
          <button
            type="button"
            class="copy-icon-button"
            onclick={() => copyCommand(sshLine)}
            aria-label="Copy SSH command"
          >
            <i class="bi bi-copy" aria-hidden="true"></i>
          </button>
        </div>
        <button type="button" class="text-button" onclick={() => auth.logout()}>Sign out</button>
      </footer>
    {:else}
      <h1>ssh-me</h1>
      <p>Create a profile, add links, and share it with one SSH command.</p>

      <button type="button" class="primary-action" onclick={() => auth.sign_in_with_google()}>
        Log in with Google
      </button>

      <footer class="demo-hook">
        <span>Try this in your terminal:</span>
        <div class="copy-command-wrap demo-command-wrap">
          <code class="demo-command">{demoSshLine}</code>
          <button
            type="button"
            class="copy-icon-button"
            onclick={() => copyCommand(demoSshLine)}
            aria-label="Copy demo SSH command"
          >
            <i class="bi bi-copy" aria-hidden="true"></i>
          </button>
        </div>
      </footer>
    {/if}
  </section>
</main>
