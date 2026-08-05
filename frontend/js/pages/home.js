import { renderPage, clearCurrentUser } from "./router.js";
import { apiFetch } from "./api.js";

export async function homepage(app, currentUser) {
  app.innerHTML = `

    <div id="home-page">
      <button id="logout-button" type="button">
        Logout
      </button>

      <hr>
      <h2>posts</h2>
      <div id="posts-container">
        loading...
      </div>
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


  //get posts

  const postsContainer = document.getElementById("posts-container");

  try {
    const response = await apiFetch("/posts", {
      method: "GET",
    });

    if (!response.ok) {
      postsContainer.innerHTML = "Failed to load posts";
      return;
    }

    const posts = await response.json();

    postsContainer.innerHTML = "";

    if (posts.length === 0) {
      postsContainer.innerHTML = "No posts yet";
      return;
    }

    posts.forEach((post) => {
      const postElement = document.createElement("div");

      postElement.className = "post";

      const titleElement = document.createElement("h3");
      titleElement.textContent = post.title;

      const contentElement = document.createElement("p");
      contentElement.textContent = post.content;

      const authorElement = document.createElement("p");
      authorElement.textContent = `Posted by: ${post.username}`;

      postElement.append(titleElement, contentElement, authorElement);

      postsContainer.appendChild(postElement);

    });

  }catch (error) {
    console.error("Error in the posts:", error);
    postsContainer.innerHTML = "feailed to load posts";
  } 


  //todo create post form 
}
