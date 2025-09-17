import Pocketbase from "pocketbase"

export function create_instance() {
  return new Pocketbase(import.meta.env.PB_URL || "http://127.0.0.1:8090/")
}

export const pb : Pocketbase = create_instance()