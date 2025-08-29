import { pb } from './Pocketbase.ts';
import { setContext,getContext } from "svelte";

const AUTH_STORE_KEY = import.meta.env.PB_AUTH_KEY

class AuthStore {
    user: typeof pb.authStore.record | null = $state(null);
    isSynced : boolean = $state(false);

    constructor() {
        // runs on mount
        $effect(() => {
            // subscribe to auth Store changes
            const unsubscribe = pb.authStore.onChange((token, model) => {
                this.user = model;
            }, true);

            this.isSynced = true;
            
            // for clean-up when destroyed
            return () => {
                unsubscribe();
            };
        });
        
    }

    async sign_in_with_google(){
        const authData = await pb.collection('users').authWithOAuth2({ provider: 'google' });
        this.user = await pb.collection('users').getOne(authData.record.id, { fields: 'id,username,status,colour' });
    }

    logout() {
        pb.authStore.clear();
        this.user = null;
    }
}

export { AuthStore };

// important if u are gonna have any SSR to prevent leakage to multiple users
export function set_auth_context() {
    const newAuthStore = new AuthStore();
    return setContext(AUTH_STORE_KEY, newAuthStore);
}

export function get_auth_context() : AuthStore {
    return getContext(AUTH_STORE_KEY);
}