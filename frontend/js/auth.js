//auth
import { apiFetch } from "./api.js";

export let currentUser = null;

export async function getCurrentUser() {
  const result = await apiFetch("/api/session");
  if (!result.ok || !result.data || !result.data.authenticated) {
    currentUser = null;
    return null;
  }
  currentUser = result.data.user;
  return result.data.user;
}

export function setCurrentUser(user) {
  currentUser = user;
}

export function clearCurrentUser() {
  currentUser = null;
}
