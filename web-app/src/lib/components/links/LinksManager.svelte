<script lang="ts">
  import { pb } from '../../pocketbase/Pocketbase.ts';
  import './LinksManager.css';

  let { profileId }: { profileId: string } = $props();

  type PlatformInputMode = 'handle' | 'url' | 'email';

  interface PlatformConfig {
    id: string;
    label: string;
    inputMode: PlatformInputMode;
    placeholder: string;
    buildUrl: (value: string) => string;
    parseValue: (url: string) => string;
  }

  function trimValue(value: string) {
    return value.trim();
  }

  function stripLeadingAt(value: string) {
    return trimValue(value).replace(/^@+/, '');
  }

  function stripTrailingSlash(value: string) {
    return trimValue(value).replace(/\/+$/, '');
  }

  function extractHandleFromUrl(url: string, pattern: RegExp) {
    const match = stripTrailingSlash(url).match(pattern);
    return match?.[1] ?? trimValue(url);
  }

  function parseMastodonAddress(value: string) {
    const cleaned = trimValue(value);
    const match = cleaned.match(/^@?([^@/]+)@([^@/\s]+)$/);

    if (match) {
      return { username: match[1], host: match[2] };
    }

    return null;
  }

  function buildMastodonUrl(value: string) {
    const parsedAddress = parseMastodonAddress(value);
    if (parsedAddress) {
      return `https://${parsedAddress.host}/@${parsedAddress.username}`;
    }

    return trimValue(value);
  }

  function parseMastodonValue(url: string) {
    const match = stripTrailingSlash(url).match(/^https?:\/\/([^/]+)\/@([^/?#]+)$/i);
    if (match) {
      return `@${match[2]}@${match[1]}`;
    }

    return trimValue(url);
  }

  const PLATFORMS = [
    {
      id: 'github',
      label: 'GitHub',
      inputMode: 'handle',
      placeholder: 'username',
      buildUrl: (value) => `https://github.com/${stripLeadingAt(value)}`,
      parseValue: (url) => extractHandleFromUrl(url, /^https?:\/\/(?:www\.)?github\.com\/([^/?#]+)$/i),
    },
    {
      id: 'linkedin',
      label: 'LinkedIn',
      inputMode: 'handle',
      placeholder: 'profile-slug',
      buildUrl: (value) => `https://linkedin.com/in/${stripLeadingAt(value)}`,
      parseValue: (url) => extractHandleFromUrl(
        url,
        /^https?:\/\/(?:[a-z]{2}\.)?linkedin\.com\/in\/([^/?#]+)$/i,
      ),
    },
    {
      id: 'youtube',
      label: 'YouTube',
      inputMode: 'handle',
      placeholder: 'username',
      buildUrl: (value) => `https://youtube.com/@${stripLeadingAt(value)}`,
      parseValue: (url) => `@${extractHandleFromUrl(
        url,
        /^https?:\/\/(?:www\.)?youtube\.com\/@([^/?#]+)$/i,
      ).replace(/^@+/, '')}`,
    },
    {
      id: 'twitter',
      label: 'Twitter / X',
      inputMode: 'handle',
      placeholder: 'username',
      buildUrl: (value) => `https://x.com/${stripLeadingAt(value)}`,
      parseValue: (url) => `@${extractHandleFromUrl(
        url,
        /^https?:\/\/(?:www\.)?(?:x|twitter)\.com\/([^/?#]+)$/i,
      ).replace(/^@+/, '')}`,
    },
    {
      id: 'instagram',
      label: 'Instagram',
      inputMode: 'handle',
      placeholder: 'username',
      buildUrl: (value) => `https://instagram.com/${stripLeadingAt(value)}`,
      parseValue: (url) => `@${extractHandleFromUrl(
        url,
        /^https?:\/\/(?:www\.)?instagram\.com\/([^/?#]+)$/i,
      ).replace(/^@+/, '')}`,
    },
    {
      id: 'mastodon',
      label: 'Mastodon',
      inputMode: 'handle',
      placeholder: 'user@server',
      buildUrl: buildMastodonUrl,
      parseValue: parseMastodonValue,
    },
    {
      id: 'bluesky',
      label: 'Bluesky',
      inputMode: 'handle',
      placeholder: 'username.bsky.social',
      buildUrl: (value) => `https://bsky.app/profile/${stripLeadingAt(value)}`,
      parseValue: (url) => extractHandleFromUrl(url, /^https?:\/\/(?:www\.)?bsky\.app\/profile\/([^/?#]+)$/i),
    },
    {
      id: 'twitch',
      label: 'Twitch',
      inputMode: 'handle',
      placeholder: 'username',
      buildUrl: (value) => `https://twitch.tv/${stripLeadingAt(value)}`,
      parseValue: (url) => extractHandleFromUrl(url, /^https?:\/\/(?:www\.)?twitch\.tv\/([^/?#]+)$/i),
    },
    {
      id: 'discord',
      label: 'Discord',
      inputMode: 'url',
      placeholder: 'discord.gg/invite-code',
      buildUrl: trimValue,
      parseValue: trimValue,
    },
    {
      id: 'devto',
      label: 'Dev.to',
      inputMode: 'handle',
      placeholder: 'username',
      buildUrl: (value) => `https://dev.to/${stripLeadingAt(value)}`,
      parseValue: (url) => extractHandleFromUrl(url, /^https?:\/\/(?:www\.)?dev\.to\/([^/?#]+)$/i),
    },
    {
      id: 'hashnode',
      label: 'Hashnode',
      inputMode: 'handle',
      placeholder: 'username',
      buildUrl: (value) => `https://hashnode.com/@${stripLeadingAt(value)}`,
      parseValue: (url) => `@${extractHandleFromUrl(
        url,
        /^https?:\/\/(?:www\.)?hashnode\.com\/@([^/?#]+)$/i,
      ).replace(/^@+/, '')}`,
    },
    {
      id: 'website',
      label: 'Website',
      inputMode: 'url',
      placeholder: 'https://example.com',
      buildUrl: trimValue,
      parseValue: trimValue,
    },
    {
      id: 'blog',
      label: 'Blog',
      inputMode: 'url',
      placeholder: 'https://blog.example.com',
      buildUrl: trimValue,
      parseValue: trimValue,
    },
    {
      id: 'email',
      label: 'Email',
      inputMode: 'email',
      placeholder: 'name@example.com',
      buildUrl: (value) => {
        const emailAddress = trimValue(value).replace(/^mailto:/i, '');
        return `mailto:${emailAddress}`;
      },
      parseValue: (url) => trimValue(url).replace(/^mailto:/i, ''),
    },
    {
      id: 'custom',
      label: 'Custom',
      inputMode: 'url',
      placeholder: 'https://example.com',
      buildUrl: trimValue,
      parseValue: trimValue,
    },
  ] satisfies PlatformConfig[];

  let links: Record<string, any>[] = $state([]);
  let newPlatform = $state('github');
  let newLabel = $state('');
  let newLinkValue = $state('');
  let editingId: string | null = $state(null);
  let editingPlatform = $state('github');
  let editingLabel = $state('');
  let editingLinkValue = $state('');
  let isSavingOrder = $state(false);

  $effect(() => {
    if (profileId) loadLinks();
  });

  async function loadLinks() {
    const result = await pb.collection('links').getList(1, 100, {
      filter: `profile = "${profileId}"`,
      sort: 'sort_order,id',
    });
    links = result.items.map((link, index) => ({
      ...link,
      sort_order: link.sort_order === undefined || link.sort_order === null
        ? index
        : Number(link.sort_order),
    }));
  }

  function getPlatformLabel(platformId: string) {
    return PLATFORMS.find((platform) => platform.id === platformId)?.label ?? platformId;
  }

  function getPlatformConfig(platformId: string) {
    return PLATFORMS.find((platform) => platform.id === platformId) ?? PLATFORMS[0];
  }

  function getPlatformInputLabel(platformId: string) {
    return getPlatformConfig(platformId).inputMode === 'handle' ? 'Handle' : 'URL';
  }

  function getLinkUrl(platformId: string, value: string) {
    return getPlatformConfig(platformId).buildUrl(value);
  }

  function getLinkValue(platformId: string, url: string) {
    return getPlatformConfig(platformId).parseValue(url);
  }

  function resetCreateForm() {
    newPlatform = 'github';
    newLabel = '';
    newLinkValue = '';
  }

  function getNextSortOrder() {
    return links.reduce(
      (maxSortOrder, link) => Math.max(maxSortOrder, Number(link.sort_order ?? -1)),
      -1,
    ) + 1;
  }

  async function persistLinkOrder(nextLinks: Record<string, any>[]) {
    const orderUpdates = nextLinks
      .map((link, index) => ({ link, index }))
      .filter(({ link, index }) => Number(link.sort_order ?? -1) !== index)
      .map(({ link, index }) => pb.collection('links').update(link.id, {
        sort_order: index,
      }));

    links = nextLinks.map((link, index) => ({
      ...link,
      sort_order: index,
    }));

    if (orderUpdates.length) {
      await Promise.all(orderUpdates);
    }
  }

  async function addLink() {
    if (!newLinkValue) return;

    const selectedLabel = newPlatform === 'custom'
      ? newLabel || 'Custom'
      : getPlatformLabel(newPlatform);
    const nextUrl = getLinkUrl(newPlatform, newLinkValue);

    if (newPlatform !== 'custom') {
      const existingLink = links.find((link) => link.platform === newPlatform);
      if (existingLink) {
        await pb.collection('links').update(existingLink.id, {
          label: selectedLabel,
          url: nextUrl,
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
      url: nextUrl,
      visible: true,
      sort_order: getNextSortOrder(),
    });

    resetCreateForm();
    await loadLinks();
  }

  function startEditing(link: Record<string, any>) {
    editingId = link.id;
    editingPlatform = link.platform;
    editingLabel = link.label ?? getPlatformLabel(link.platform);
    editingLinkValue = getLinkValue(link.platform, link.url);
  }

  async function saveEditing() {
    if (!editingId || !editingLinkValue) return;

    await pb.collection('links').update(editingId, {
      platform: editingPlatform,
      label: editingPlatform === 'custom'
        ? editingLabel || 'Custom'
        : getPlatformLabel(editingPlatform),
      url: getLinkUrl(editingPlatform, editingLinkValue),
      visible: true,
    });

    editingId = null;
    editingLabel = '';
    editingLinkValue = '';
    await loadLinks();
  }

  async function removeLink(id: string) {
    const nextLinks = links.filter((link) => link.id !== id);

    await pb.collection('links').delete(id);
    if (editingId === id) editingId = null;
    await persistLinkOrder(nextLinks);
  }

  async function moveLink(fromIndex: number, offset: number) {
    if (isSavingOrder) return;

    const toIndex = fromIndex + offset;
    if (toIndex < 0 || toIndex >= links.length) return;

    isSavingOrder = true;

    try {
      const nextLinks = [...links];
      [nextLinks[fromIndex], nextLinks[toIndex]] = [nextLinks[toIndex], nextLinks[fromIndex]];
      await persistLinkOrder(nextLinks);
    } finally {
      isSavingOrder = false;
    }
  }
</script>

<section class="links-panel">
  <p class="section-label">Links</p>

  {#if links.length}
    <div class="links-list">
      {#each links as link, index (link.id)}
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

            <input
              type={getPlatformConfig(editingPlatform).inputMode === 'email' ? 'email' : 'text'}
              placeholder={getPlatformConfig(editingPlatform).placeholder}
              aria-label={getPlatformInputLabel(editingPlatform)}
              bind:value={editingLinkValue}
            >

            <div class="link-actions">
              <button type="submit" class="button-primary">Save</button>
              <button type="button" class="button-secondary" onclick={() => editingId = null}>Cancel</button>
            </div>
          </form>
        {:else}
          <div class="link-row">
            <div class="link-order-controls">
              <button
                type="button"
                class="link-order-button"
                onclick={() => moveLink(index, -1)}
                disabled={index === 0 || isSavingOrder}
                aria-label={`Move ${(link.label || getPlatformLabel(link.platform))} up`}
              >
                <i class="bi bi-chevron-up" aria-hidden="true"></i>
              </button>
              <button
                type="button"
                class="link-order-button"
                onclick={() => moveLink(index, 1)}
                disabled={index === links.length - 1 || isSavingOrder}
                aria-label={`Move ${(link.label || getPlatformLabel(link.platform))} down`}
              >
                <i class="bi bi-chevron-down" aria-hidden="true"></i>
              </button>
            </div>
            <span class="link-label">{link.label || getPlatformLabel(link.platform)}</span>
            <a class="link-url" href={link.url} target="_blank" rel="noreferrer">{link.url}</a>
            <div class="link-actions">
              <button
                type="button"
                class="icon-button"
                onclick={() => startEditing(link)}
                aria-label={`Edit ${(link.label || getPlatformLabel(link.platform))}`}
              >
                <i class="bi bi-pencil" aria-hidden="true"></i>
              </button>
              <button
                type="button"
                class="icon-button"
                onclick={() => removeLink(link.id)}
                aria-label={`Remove ${(link.label || getPlatformLabel(link.platform))}`}
              >
                <i class="bi bi-trash" aria-hidden="true"></i>
              </button>
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

    <input
      type={getPlatformConfig(newPlatform).inputMode === 'email' ? 'email' : 'text'}
      placeholder={getPlatformConfig(newPlatform).placeholder}
      aria-label={getPlatformInputLabel(newPlatform)}
      bind:value={newLinkValue}
    >
    <button type="submit" class="button-primary">Add</button>
  </form>
</section>
