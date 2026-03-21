<script lang="ts">
  import { pb } from '../../pocketbase/Pocketbase.ts';

  let { profileId }: { profileId: string } = $props();

  const PLATFORMS = [
    { id: 'github',    label: 'GitHub' },
    { id: 'linkedin',  label: 'LinkedIn' },
    { id: 'youtube',   label: 'YouTube' },
    { id: 'twitter',   label: 'Twitter / X' },
    { id: 'instagram', label: 'Instagram' },
    { id: 'mastodon',  label: 'Mastodon' },
    { id: 'bluesky',   label: 'Bluesky' },
    { id: 'twitch',    label: 'Twitch' },
    { id: 'discord',   label: 'Discord' },
    { id: 'devto',     label: 'Dev.to' },
    { id: 'hashnode',  label: 'Hashnode' },
    { id: 'website',   label: 'Website' },
    { id: 'blog',      label: 'Blog' },
    { id: 'email',     label: 'Email' },
  ];

  let links: Record<string, any>[] = $state([]);

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

  // Map platform id -> existing link record (or null)
  let platformLinks = $derived(
    Object.fromEntries(PLATFORMS.map(p => [p.id, links.find(l => l.platform === p.id) ?? null]))
  );

  let customLinks = $derived(links.filter(l => l.platform === 'custom'));

  // Platform link editing
  let editingPlatform: string | null = $state(null);
  let editingUrl = $state('');
  let editingLabel = $state('');

  function openPlatformEditor(platformId: string) {
    const existing = platformLinks[platformId];
    editingPlatform = platformId;
    editingUrl = existing?.url ?? '';
    editingLabel = existing?.label ?? '';
  }

  async function savePlatformLink(platformId: string) {
    const existing = platformLinks[platformId];
    const platformLabel = PLATFORMS.find(p => p.id === platformId)?.label ?? platformId;
    const data = {
      profile: profileId,
      platform: platformId,
      url: editingUrl,
      label: editingLabel || platformLabel,
      visible: true,
      sort_order: existing?.sort_order ?? links.length,
    };
    if (existing) {
      await pb.collection('links').update(existing.id, data);
    } else {
      await pb.collection('links').create(data);
    }
    editingPlatform = null;
    editingUrl = '';
    editingLabel = '';
    await loadLinks();
  }

  async function deletePlatformLink(id: string) {
    await pb.collection('links').delete(id);
    await loadLinks();
  }

  // Custom link editing
  let showAddCustom = $state(false);
  let newCustomLabel = $state('');
  let newCustomUrl = $state('');
  let editingCustomId: string | null = $state(null);
  let editingCustomLabel = $state('');
  let editingCustomUrl = $state('');

  async function addCustomLink() {
    if (!newCustomUrl) return;
    await pb.collection('links').create({
      profile: profileId,
      platform: 'custom',
      label: newCustomLabel,
      url: newCustomUrl,
      visible: true,
      sort_order: links.length,
    });
    newCustomLabel = '';
    newCustomUrl = '';
    showAddCustom = false;
    await loadLinks();
  }

  function openCustomEditor(link: Record<string, any>) {
    editingCustomId = link.id;
    editingCustomLabel = link.label;
    editingCustomUrl = link.url;
  }

  async function saveCustomLink() {
    if (!editingCustomId) return;
    await pb.collection('links').update(editingCustomId, {
      label: editingCustomLabel,
      url: editingCustomUrl,
    });
    editingCustomId = null;
    await loadLinks();
  }

  async function deleteCustomLink(id: string) {
    await pb.collection('links').delete(id);
    await loadLinks();
  }
</script>

<section>
  <h3>Platform Links</h3>
  <div class="platform-grid">
    {#each PLATFORMS as platform}
      {@const existing = platformLinks[platform.id]}
      <div class="platform-row">
        {#if editingPlatform === platform.id}
          <div class="platform-editor">
            <strong>{platform.label}</strong>
            <input type="url" placeholder="URL" bind:value={editingUrl} />
            <input type="text" placeholder="Custom label (optional)" bind:value={editingLabel} />
            <button onclick={() => savePlatformLink(platform.id)}>Save</button>
            <button onclick={() => editingPlatform = null}>Cancel</button>
          </div>
        {:else}
          <span class="platform-name">{platform.label}</span>
          {#if existing}
            <span class="platform-url">{existing.url}</span>
            <button onclick={() => openPlatformEditor(platform.id)}>Edit</button>
            <button onclick={() => deletePlatformLink(existing.id)}>Remove</button>
          {:else}
            <button onclick={() => openPlatformEditor(platform.id)}>+ Add</button>
          {/if}
        {/if}
      </div>
    {/each}
  </div>

  <h3>Custom Links</h3>
  <div class="custom-links">
    {#each customLinks as link}
      <div class="custom-row">
        {#if editingCustomId === link.id}
          <input type="text" placeholder="Label" bind:value={editingCustomLabel} />
          <input type="url" placeholder="URL" bind:value={editingCustomUrl} />
          <button onclick={saveCustomLink}>Save</button>
          <button onclick={() => editingCustomId = null}>Cancel</button>
        {:else}
          <span>{link.label || link.url}</span>
          <span class="platform-url">{link.url}</span>
          <button onclick={() => openCustomEditor(link)}>Edit</button>
          <button onclick={() => deleteCustomLink(link.id)}>Remove</button>
        {/if}
      </div>
    {/each}

    {#if showAddCustom}
      <div class="custom-row">
        <input type="text" placeholder="Label (e.g. My Blog)" bind:value={newCustomLabel} />
        <input type="url" placeholder="URL" bind:value={newCustomUrl} />
        <button onclick={addCustomLink}>Add</button>
        <button onclick={() => showAddCustom = false}>Cancel</button>
      </div>
    {:else}
      <button onclick={() => showAddCustom = true}>+ Add custom link</button>
    {/if}
  </div>
</section>

<style>
  section {
    margin-top: 1.5rem;
  }

  h3 {
    margin-bottom: 0.5rem;
  }

  .platform-grid, .custom-links {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    margin-bottom: 1rem;
  }

  .platform-row, .custom-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .platform-name {
    width: 8rem;
    font-weight: bold;
  }

  .platform-url {
    color: #666;
    font-size: 0.9em;
    flex: 1;
  }

  .platform-editor {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
    width: 100%;
  }

  input {
    padding: 0.25rem 0.5rem;
    font-size: 0.9em;
  }
</style>
