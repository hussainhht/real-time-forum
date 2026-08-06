import { renderPage } from "./router.js";
import { clearCurrentUser } from "./auth.js";
import { apiFetch } from "./api.js";
import { sidebar } from "./pages/home/sidepar.js";
import { topbar } from "./pages/home/topbar.js";
import { feedPage } from "./pages/home/feed.js";


export async function homepage(app, currentUser) {
  app.innerHTML = `
    <div id="home-page">
      ${sidebar()}
      >
      <section class="home-main">
        ${topbar(currentUser)}

        <main id="home-container">
          ${feedPage()}
        </main>
      </section>
      <button id="logout-button" type="button">
        Logout
      </button>
    </div>

    `;

  const logoutButton = document.getElementById("logout-button");
  logoutButton.addEventListener("click", async () => {
    const response = await fetch("/logout", {
      method: "POST",
      credentials: "include",
    });

    if (!response.ok) {
      return;
    }
    clearCurrentUser();
    renderPage("login");
  });



}
