//auth
import { apiFetch } from "./api.js";

export async function getCurrentUser() {
    const result = await apiFetch("api/session");
    if (!result.ok || !result.data || !result.data.authenticated) {
        return null;
    }
    return result.data.user;
}
