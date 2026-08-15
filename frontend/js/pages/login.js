import { renderPage } from "../router.js";
import { apiFetch, errorMessage } from "../api.js";
import { getCurrentUser, setCurrentUser } from "../auth.js";
import { connectWebsocket } from "../websocket.js";

export function LoginPage(app) {
  app.innerHTML = `
    <main class="auth-page">
      <section class="auth-card">
        <form id="login-form">
          <h1>Login</h1>
          <input
            type="text"
            id="identifier"
            placeholder="Email or username"
            required
          />
          <input
            type="password"
            id="password"
            placeholder="Password"
            required
          />
          <button type="submit">Login</button>
        </form>
        <button type="button" id="go-register" class="login-link">
          Create a new account
        </button>
        <p id="login-message" class="form-message"></p>
      </section>
    </main>
  `;

  const form = document.getElementById("login-form");
  const registerButton = document.getElementById("go-register");

  form.addEventListener("submit", handleLogin);

  registerButton.addEventListener("click", () => {
    renderPage("register");
  });
}

async function handleLogin(event) {
  event.preventDefault();

  const identifier = document.getElementById("identifier").value.trim();
  const password = document.getElementById("password").value;
  const message = document.getElementById("login-message");

  if (identifier === "" || password === "") {
    message.textContent = "Please fill all fields";
    return;
  }

  const response = await apiFetch("/login", {
    method: "POST",
    body: { identifier, password },
  });

  if (!response.ok) {
    message.textContent = errorMessage(response, "Login failed");
    return;
  }

  const user = await getCurrentUser();

  if (!user) {
    message.textContent = "Login succeeded but session could not be verified";
    return;
  }

  message.textContent = "";
  setCurrentUser(user);
  connectWebsocket();
  renderPage("home");
}