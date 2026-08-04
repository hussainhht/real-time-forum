//auth

export async function getCurrentUser() {
  try {
    const response = await fetch("/api/session", {
      method: "GET",
      credentials: "same-origin", // this is tha http://localhost:8080  the origin have three things like protocol and domain and port
    });

    if (!response.ok) {
      return null;
    }
    const data = await response.json();
    if (!data.authenticated) {
      return null;
    }

    return data.user; // this is user data from data
  } catch (error) {
    console.error("Failed to check session:", error);
    return null;
  }
}
