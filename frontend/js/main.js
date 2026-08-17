import { startApp, renderPage } from "./router.js";
import { currentUser } from "./auth.js";

// Keep the rendered page in sync with the URL hash whenever it changes
// (covers link clicks, and the browser back/forward buttons).
window.addEventListener("hashchange", () => {
  const page = window.location.hash.slice(1) || "home";
  renderPage(page);
});

// const path = window.location.pathname;

// if (path !== "/") {
//   renderErrorPage(
//     app,
//     "404 Not Found",
//     "The page you are looking for does not exist.",
//   );
//   return;
// }

// Initialize the application: check auth and render the first page.
startApp().then(() => {
  const hash = window.location.hash.slice(1);

  if (hash) {
    // Refresh or a bookmarked/shared link: render whatever the hash points to.
    renderPage(hash);
  } else {
    // No hash yet (first load at "/"): reflect the page startApp() just
    // rendered in the URL so it's bookmarkable and back/forward-safe.
    window.location.hash = currentUser ? "#home" : "#login";
  }
});
