import { renderPage } from "../router.js";
import { clearCurrentUser } from "../auth.js";
import { apiFetch, errorMessage } from "../api.js";
import { renderSidebar, updateSidebarUsers } from "../pages/home/sidepar.js";
import { renderTopbar } from "../pages/home/topbar.js";
import { renderFeedView } from "../pages/home/feed.js";
import { renderCreatePostView } from "../pages/home/createPost.js";
import { renderPostView } from "../pages/home/postView.js";


export async function homepage(app, currentUser) {
  app.innerHTML = `
    <div id="home-page">
      ${renderSidebar()}

      <section class="home-main">
        ${renderTopbar(currentUser)}

        <main id="home-container">
        </main>
      </section>
    </div>

    `;

  setupHomeEvents();

  await updateSidebarUsers();
  await navigateHome("feed");
}

function setupHomeEvents() {
  document
  .getElementById("logout-button")
  .addEventListener("click", handleLogout);

  document
  .getElementById("create-post-button")
  .addEventListener("click",navigateHome("create-post"));
}

async function handleLogout() {
  const response = await apiFetch("/logout", {
    method: "POST",
  });

  if (!response.ok) {
    console.error(errorMessage(response, "Logout failed"));
    return;
  }
  clearCurrentUser();
  renderPage("login");
}

export async function navigateHome(view, data = {}) {
  const container = document.getElementById("home-container");

  if (!container) {
    console.error("Home container not found");
    return;
  }

  switch (view) {
    case "feed":
      renderFeedView(container);
      break;
    case "create-post":
      renderCreatePostView(container);
      break;
    case "post":
      renderPostView(container, data.postId);
      break;
    default:
      renderFeedView(container);
      break;
  }
}