import { pb } from './Pocketbase.ts';
import { setContext, getContext } from "svelte";

const AUTH_STORE_KEY = import.meta.env.PB_AUTH_KEY

class AuthStore {
    user: Record<string, any> | null = $state(null);
    isSynced: boolean = $state(false);

    constructor() {
        $effect(() => {
            const unsubscribe = pb.authStore.onChange(async () => {
                if (pb.authStore.isValid && pb.authStore.record) {
                    this.user = await this.fetchOrCreateProfile(pb.authStore.record.id)
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
            return await pb.collection('profiles').create({ user: userId, avatar_url: avatarUrl ?? '' })
        } catch (err) {
            console.error("Failed to fetch/create profile:", err)
            return null
        }
    }

    async sign_in_with_google() {
        const result = await pb.collection('users').authWithOAuth2({ provider: 'google' });
        const avatarUrl = result.meta?.avatarUrl ?? '';
        if (result.record) {
            this.user = await this.fetchOrCreateProfile(result.record.id, avatarUrl);
        }
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
