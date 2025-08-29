import { pb } from '../lib/Pocketbase.ts';
import { setContext,getContext } from "svelte";

const AUTH_STORE_KEY = import.meta.env.PB_AUTH_KEY

export class AuthStore {
    user: typeof pb.authStore.record | null = $state(null);
    isSynced : boolean = $state(false);

    constructor() {
        // runs on mount
        $effect(() => {
            if(pb.authStore.isValid){
                this.user = pb.authStore.record;
            }
            this.isSynced = true;
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

// important if u are gonna have any SSR to prevent leakage to multiple users
export function set_auth_context() {
    const newAuthStore = new AuthStore();
    return  setContext(AUTH_STORE_KEY, newAuthStore);
}

export function get_auth_context() {
    return getContext(AUTH_STORE_KEY);
}