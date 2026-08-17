//auth
import { apiFetch } from "./api.js";
import { renderErrorPage } from "./pages/error.js";
import { app } from "./router.js";

export let currentUser = null;

export async function getCurrentUser() {
  const result = await apiFetch("/api/session");

  if (result.status === 500) {
    renderErrorPage(
      app,
      "500 Internal Server Error",
      "Something went wrong on our end. Please try again later."
    );
    currentUser = null;
    return null;
  }

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
