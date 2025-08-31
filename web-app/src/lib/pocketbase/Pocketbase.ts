import Pocketbase from "pocketbase"

export function create_instance() {
  return new Pocketbase(import.meta.env.PB_URL || "https://ssh-me.com/")
}

export const pb : Pocketbase = create_instance()