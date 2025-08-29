import Pocketbase from "pocketbase"

export function createInstance() {
  return new Pocketbase(import.meta.env.PB_URL || "https://ssh-me.com/")
}

export const pb : Pocketbase = createInstance()