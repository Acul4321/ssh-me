import Pocketbase from "pocketbase"

export function create_instance() {
  return new Pocketbase(import.meta.env.VITE_PB_URL || "http://localhost:8080")
}

export const pb : Pocketbase = create_instance()