import { pb } from './Pocketbase.ts';
import { setContext, getContext } from "svelte";

const AUTH_STORE_KEY = "auth"

class AuthStore {
    user: Record<string, any> | null = $state(null);
    isSynced: boolean = $state(false);

    constructor() {
        $effect(() => {
            const unsubscribe = pb.authStore.onChange(async () => {
                if (pb.authStore.isValid && pb.authStore.record) {
                    const record = pb.authStore.record;
                    const avatarUrl = record.avatarUrl ?? record.avatar ?? '';
                    this.user = await this.fetchOrCreateProfile(record.id, avatarUrl)
                } else {
                    this.user = null
                }
            }, true);

            this.isSynced = true;

            return () => {
                unsubscribe();
            };
        });
    }

    private async fetchOrCreateProfile(userId: string, avatarUrl?: string): Promise<Record<string, any> | null> {
        try {
            const result = await pb.collection('profiles').getList(1, 1, {
                filter: `user = "${userId}"`
            })
            if (result.totalItems > 0) {
                return result.items[0]
            }
            try {
                return await pb.collection('profiles').create({ user: userId, avatar_url: avatarUrl ?? '' })
            } catch {
                // creation failed (likely a race condition hitting the unique constraint)
                // fall back to fetching the profile that the other call created
                const retry = await pb.collection('profiles').getList(1, 1, {
                    filter: `user = "${userId}"`
                })
                return retry.items[0] ?? null
            }
        } catch (err) {
            console.error("Failed to fetch/create profile:", err)
            return null
        }
    }

    async sign_in_with_google() {
        await pb.collection('users').authWithOAuth2({ provider: 'google' });
        // profile is created/fetched by the authStore onChange listener above
    }

    logout() {
        pb.authStore.clear();
        this.user = null;
    }
}

export { AuthStore };

export function set_auth_context() {
    const newAuthStore = new AuthStore();
    return setContext(AUTH_STORE_KEY, newAuthStore);
}

export function get_auth_context(): AuthStore {
    return getContext(AUTH_STORE_KEY);
}
