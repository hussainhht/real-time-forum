import { LoginPage } from "./pages/login.js";
import { RegisterPage } from "./pages/register.js";
import { homepage } from "./pages/home.js";
import { getCurrentUser } from "./auth.js";

export const app = document.getElementById("app");
let currentUser = null;

export function renderPage(page) {

  const protectedPages = ["home"];
  const guestPages = ["login", "register"];
  //cases to render if auth is in wrong direction.
  if (protectedPages.includes(page) && !currentUser) {
    LoginPage(app);
    return;
  }
  if (guestPages.includes(page) && currentUser) {
    homepage(app, currentUser);
    return;
  }

  switch (page) {
    case "login":
      LoginPage(app);
      break;
    case "register":
      RegisterPage(app);
      break;
    case "home":
      homepage(app, currentUser);
      break;
    default:
      if (currentUser) {
        homepage(app, currentUser);
      } else {
        LoginPage(app);
      }
      return;
  }

}

export async function startApp() {
  currentUser = await getCurrentUser();
  if (currentUser) {
    renderPage("home");
    return;
  }
  renderPage("login");
}

export function setCurrentUser(user) {
  currentUser = user;
}

export function clearCurrentUser() {
  currentUser = null;
}