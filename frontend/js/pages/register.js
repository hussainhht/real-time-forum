import { apiFetch, errorMessage } from "../api.js";
import { renderErrorPage } from "./error.js";
import { app } from "../router.js";

export function RegisterPage(app) {
  app.innerHTML = `
    <main class="auth-page">
      <section class="auth-card">
        <h1>Create account</h1>
        <p class="auth-subtitle">Register to join the forum</p>

        <form id="register-form">
          <div class="name-row">
            <input
              type="text"
              id="first-name"
              placeholder="First name"
              required
            />
            <input
              type="text"
              id="last-name"
              placeholder="Last name"
              required
            />
          </div>

          <input type="number" id="age" placeholder="Age" required />

          <select id="gender" required>
            <option value="">Choose gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
          </select>

          <input type="email" id="email" placeholder="Email" required />
          <input type="text" id="username" placeholder="Username" required />
          <input
            type="password"
            id="password"
            placeholder="Password"
            required
          />
          <input
          type="password"
           id="confirm-password"
           placeholder="Confirm password"
           required
          />
          <p id="register-message" class="form-message"></p>
          <button class="primary-button" type="submit">Register</button>
        </form>
        <button id="go-login" type="button" class="login-link">
          Already have an account?
        </button>
      </section>

      
    </main>
  `;

  const form = document.getElementById("register-form");
  const loginButton = document.getElementById("go-login");

  form.addEventListener("submit", handleRegister);

  loginButton.addEventListener("click", () => {
    window.location.hash = "#login";
  });
}

async function handleRegister(event) {
  event.preventDefault();

  const first_name = document.getElementById("first-name").value.trim();
  const last_name = document.getElementById("last-name").value.trim();
  const age = Number(document.getElementById("age").value);
  const gender = document.getElementById("gender").value;
  const email = document.getElementById("email").value.trim();
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value;
  const message = document.getElementById("register-message");
  const confirm_password = document.getElementById("confirm-password").value;
  if (
    first_name === "" ||
    last_name === "" ||
    !age ||
    gender === "" ||
    email === "" ||
    username === "" ||
    password === "" ||
    confirm_password === ""
  ) {
    message.textContent = "Please fill all fields";
    return;
  }

  if (!Number.isInteger(age) || age <= 0) {
    message.textContent = "Please enter a valid age";
    return;
  }

  if (password !== confirm_password) {
    message.textContent = "Passwords do not match";
    return;
  }

  const response = await apiFetch("/register", {
    method: "POST",
    body: { first_name, last_name, age, gender, email, username, password, confirm_password },
  });


  if (response.status === 500) {
    renderErrorPage(
      app,
      "500 Internal Server Error",
      "Something went wrong on our end. Please try again later."
    );
    return;
  } else if (!response.ok) {
    message.textContent = errorMessage(response, "Registration failed");
    return;
  }

  message.textContent = "Registration successful";
  window.location.hash = "#login";
}
