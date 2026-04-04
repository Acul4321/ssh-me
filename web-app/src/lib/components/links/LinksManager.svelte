<script lang="ts">
  import { pb } from '../../pocketbase/Pocketbase.ts';
  import './LinksManager.css';

  let { profileId }: { profileId: string } = $props();

  const PLATFORMS = [
    { id: 'github', label: 'GitHub' },
    { id: 'linkedin', label: 'LinkedIn' },
    { id: 'youtube', label: 'YouTube' },
    { id: 'twitter', label: 'Twitter / X' },
    { id: 'instagram', label: 'Instagram' },
    { id: 'mastodon', label: 'Mastodon' },
    { id: 'bluesky', label: 'Bluesky' },
    { id: 'twitch', label: 'Twitch' },
    { id: 'discord', label: 'Discord' },
    { id: 'devto', label: 'Dev.to' },
    { id: 'hashnode', label: 'Hashnode' },
    { id: 'website', label: 'Website' },
    { id: 'blog', label: 'Blog' },
    { id: 'email', label: 'Email' },
    { id: 'custom', label: 'Custom' },
  ];

  let links: Record<string, any>[] = $state([]);
  let newPlatform = $state('github');
  let newLabel = $state('');
  let newUrl = $state('');
  let editingId: string | null = $state(null);
  let editingPlatform = $state('github');
  let editingLabel = $state('');
  let editingUrl = $state('');

  $effect(() => {
    if (profileId) loadLinks();
  });

  async function loadLinks() {
    const result = await pb.collection('links').getList(1, 100, {
      filter: `profile = "${profileId}"`,
      sort: 'sort_order',
    });
    links = result.items;
  }

  function getPlatformLabel(platformId: string) {
    return PLATFORMS.find((platform) => platform.id === platformId)?.label ?? platformId;
  }

  function resetCreateForm() {
    newPlatform = 'github';
    newLabel = '';
    newUrl = '';
  }

  async function addLink() {
    if (!newUrl) return;

    const selectedLabel = newPlatform === 'custom'
      ? newLabel || 'Custom'
      : getPlatformLabel(newPlatform);

    if (newPlatform !== 'custom') {
      const existingLink = links.find((link) => link.platform === newPlatform);
      if (existingLink) {
        await pb.collection('links').update(existingLink.id, {
          label: selectedLabel,
          url: newUrl,
          visible: true,
        });
        resetCreateForm();
        await loadLinks();
        return;
      }
    }

    await pb.collection('links').create({
      profile: profileId,
      platform: newPlatform,
      label: selectedLabel,
      url: newUrl,
      visible: true,
      sort_order: links.length,
    });

    resetCreateForm();
    await loadLinks();
  }

  function startEditing(link: Record<string, any>) {
    editingId = link.id;
    editingPlatform = link.platform;
    editingLabel = link.label ?? getPlatformLabel(link.platform);
    editingUrl = link.url;
  }

  async function saveEditing() {
    if (!editingId || !editingUrl) return;

    await pb.collection('links').update(editingId, {
      platform: editingPlatform,
      label: editingPlatform === 'custom'
        ? editingLabel || 'Custom'
        : getPlatformLabel(editingPlatform),
      url: editingUrl,
      visible: true,
    });

    editingId = null;
    editingLabel = '';
    editingUrl = '';
    await loadLinks();
  }

  async function removeLink(id: string) {
    await pb.collection('links').delete(id);
    if (editingId === id) editingId = null;
    await loadLinks();
  }
</script>

<section class="links-panel">
  <p class="section-label">Links</p>

  {#if links.length}
    <div class="links-list">
      {#each links as link (link.id)}
        {#if editingId === link.id}
          <form class="link-edit-row" onsubmit={(event) => {
            event.preventDefault();
            saveEditing();
          }}>
            <select bind:value={editingPlatform}>
              {#each PLATFORMS as platform}
                <option value={platform.id}>{platform.label}</option>
              {/each}
            </select>

            {#if editingPlatform === 'custom'}
              <input type="text" placeholder="Custom label" bind:value={editingLabel}>
            {/if}

            <input type="url" placeholder="https://example.com" bind:value={editingUrl}>

            <div class="link-actions">
              <button type="submit" class="button-primary">Save</button>
              <button type="button" class="button-secondary" onclick={() => editingId = null}>Cancel</button>
            </div>
          </form>
        {:else}
          <div class="link-row">
            <span class="link-label">{link.label || getPlatformLabel(link.platform)}</span>
            <a class="link-url" href={link.url} target="_blank" rel="noreferrer">{link.url}</a>
            <div class="link-actions">
              <button type="button" class="button-secondary" onclick={() => startEditing(link)}>Edit</button>
              <button type="button" class="button-secondary" onclick={() => removeLink(link.id)}>Remove</button>
            </div>
          </div>
        {/if}
      {/each}
    </div>
  {:else}
    <p class="empty-state">No links yet.</p>
  {/if}

  <form class="new-link-form" onsubmit={(event) => {
    event.preventDefault();
    addLink();
  }}>
    <select bind:value={newPlatform}>
      {#each PLATFORMS as platform}
        <option value={platform.id}>{platform.label}</option>
      {/each}
    </select>

    {#if newPlatform === 'custom'}
      <input type="text" placeholder="Custom label" bind:value={newLabel}>
    {/if}

    <input type="url" placeholder="https://example.com" bind:value={newUrl}>
    <button type="submit" class="button-primary">Add</button>
  </form>
</section>
