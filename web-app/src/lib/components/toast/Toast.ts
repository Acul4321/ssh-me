import { writable, type Writable } from "svelte/store";

export const toasts : Writable<ToastInfo[]> = writable([]);

export enum toastTypes { SUCCESS, ERROR, INFO}

export interface ToastInfo {
    id : number;
    message : string;
    type : toastTypes;
    duration?: number;
    dismissible? : boolean;
}

export function add_toast(message : string, type : toastTypes, duration : number = 1000, dismissible : boolean = true){
    // id for easily finding the toast
    const id = (Math.random() * 10000);

    const info : ToastInfo = {
        id,
        message,
        type,
        duration,
        dismissible
    };

    // add new toast
    toasts.update((all) => [...all, info]);

    // set timeout for destroying of toast after timeout
    if(info.duration) setTimeout(() => delete_toast(info.id), info.duration);
}

export function delete_toast(id : number){
    toasts.update((all) => all.filter((t) => t.id !== id));
}