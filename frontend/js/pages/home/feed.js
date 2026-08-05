import { apiFetch, errorMessage } from "../../api.js";
import { navigateHome } from "../home.js"; //todo

export async function renderFeedView(box) {
    box.innerHTML = `
    <h2>Posts</h2>
    <div id="posts-container">loading...</div>
    `;

    const postsContainer = document.getElementById("posts-container");

    const response = await apiFetch("/posts", {
    method: "GET",
    });


    if (!document.body.contains(postsContainer)) return;

    if (!response.ok) {
        postsContainer.textContent = errorMessage(response, "Failed to load posts");
        return;
    }

    const posts = response.data || [];

    postsContainer.innerHTML = "";

    if (posts.length === 0) {
        postsContainer.textContent = "No posts yet";
        return;
    }

    posts.forEach((post) => {
        const postElement = document.createElement("div");
        postElement.className = "post-card";

        const categoryElement = document.createElement("span");
        categoryElement.className = "post-category";
        categoryElement.textContent = post.category;

        const titleElement = document.createElement("h3");
        titleElement.textContent = post.title;

        const authorElement = document.createElement("p");
        authorElement.textContent = `Posted by: ${post.username}`;

        postElement.append(categoryElement, titleElement, authorElement);

        postElement.addEventListener("click", () => {
            navigateHome("post", { postId: post.id }); //todo: navigate to post view
        });

        postsContainer.appendChild(postElement);
    });

}