import { renderPage } from "./router.js";

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
          <button class="primary-button" type="submit">Register</button>
        </form>
     <button id="go-login" type="button" class="login-link">
        Already have an account?
      </button>
      </section>


      <p id="register-message" class="form-message"></p>
    </main>
    `;

  const form = document.getElementById("register-form");
  const loginButton = document.getElementById("go-login");

  form.addEventListener("submit", handleRegister);

  loginButton.addEventListener("click", () => {
    renderPage("login");
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
  const password = document.getElementById("password").value.trim();
  const message = document.getElementById("register-message");

  if (
    first_name === "" ||
    last_name === "" ||
    !age ||
    gender === "" ||
    email === "" ||
    username === "" ||
    password === ""
  ) {
    message.textContent = "pless fill all fiealds";
    return;
  }

  try {
    const response = await fetch("/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },

      body: JSON.stringify({
        first_name,
        last_name,
        age,
        gender,
        email,
        username,
        password,
      }),
    });

    if (!response.ok) {
      const errorMassage = await response.text();
      message.textContent = errorMassage;
      return;
    }
    message.textContent = "Registerion succesful";
    renderPage("login");
  } catch (error) {
    console.log(error);
    message.textContent = "could not coonect to sevver";
  }
}
