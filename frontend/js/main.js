import { startApp, renderPage } from "./router.js";
import { currentUser } from "./auth.js";
import { renderErrorPage } from "./pages/error.js";

window.addEventListener("hashchange", () => {
  const page = window.location.hash.slice(1) || "home";
  renderPage(page);
});

startApp().then(() => {
  const path = window.location.pathname;
  const hash = window.location.hash.slice(1);

  if (path !== "/") {
    renderErrorPage(
      app,
      "404 Not Found",
      "The page you are looking for does not exist."
    );
    return;
  }

  // Hash routing
  if (hash) {
    renderPage(hash);
    return;
  }

  window.location.hash = currentUser ? "#home" : "#login";
});