import { LoginPage } from "./pages/login.js";
import { RegisterPage } from "./pages/register.js";
import { homepage } from "./pages/home.js";
import { renderErrorPage } from "./pages/error.js";
import { getCurrentUser, currentUser } from "./auth.js";

export const app = document.getElementById("app");

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
        renderErrorPage(
          app,
          "404 Not Found",
          "The page you are looking for does not exist.",
        );
      } else {
        LoginPage(app);
      }
      return;
  }
}

export async function startApp() {
  await getCurrentUser();
  if (currentUser) {
    renderPage("home");
    return;
  }
  renderPage("login");
}
