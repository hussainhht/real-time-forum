import { apiFetch, errorMessage } from "../../api.js";
import { navigateHome } from "../home.js";

export async function renderPostView(box, postId) {
    if (!postId) {
    navigateHome("feed");
    return;
    }

    box.innerHTML = `
    <button id="back-to-feed" type="button">&larr; Back to feed</button>
    <div id="post-detail">loading...</div>
    `;

    let backButton = document.getElementById("back-to-feed");
    backButton.addEventListener("click", () => navigateHome("feed"));


    const detail = document.getElementById("post-detail");

    const response = await apiFetch(`/posts/${postId}`, {
    method: "GET",
    });

    if (!document.body.contains(detail)) return;

    if (!response.ok) {
    detail.textContent = errorMessage(response, "Failed to load post");
    return;
    }

    const post = response.data;

    detail.innerHTML = "";

    const categoryElement = document.createElement("span");
    categoryElement.className = "post-category";
    categoryElement.textContent = post.category;

    const titleElement = document.createElement("h2");
    titleElement.textContent = post.title;

    const authorElement = document.createElement("p");
    authorElement.textContent = `Posted by: ${post.username}`;

    const contentElement = document.createElement("p");
    contentElement.className = "post-full-content";
    contentElement.textContent = post.content;

    const commentsSection = document.createElement("section");
    commentsSection.id = "comments-section";
    commentsSection.innerHTML = `
    <h3>Comments</h3>
    <ul id="comments-list" class="comments-list">
    <li class="comments-list-message">loading...</li>
    </ul>
    <form id="add-comment-form">
    <textarea id="comment-content" placeholder="Write a comment..." required></textarea>
    <button type="submit">Comment</button>
    </form>
    <p id="comment-message" class="form-message"></p>
    `;

    detail.append(
    categoryElement,
    titleElement,
    authorElement,
    contentElement,
    commentsSection
    );

    document
    .getElementById("add-comment-form")
    .addEventListener("submit", (event) => handleAddComment(event, postId));

    renderComments(postId);
    }

    async function handleAddComment(event, postId) {
    event.preventDefault();

    const contentInput = document.getElementById("comment-content");
    const message = document.getElementById("comment-message");

    const content = contentInput.value.trim();
    if (content === "") {
    message.textContent = "Comment can't be empty";
    return;
    }

    const result = await apiFetch(`/posts/${postId}/comments`, {
    method: "POST",
    body: { content },
    });

    if (!document.body.contains(message)) return;

    if (!result.ok) {
    message.textContent = errorMessage(result, "Failed to add comment");
    return;
    }

    message.textContent = "";
    contentInput.value = "";

    renderComments(postId);
    }

    async function renderComments(postId) {
    const list = document.getElementById("comments-list");
    if (!list) return;

    list.innerHTML = `<li class="comments-list-message">loading...</li>`;

    const response = await apiFetch(`/posts/${postId}/comments`, {
    method: "GET",
    });

    if (!document.body.contains(list)) return;

    if (!response.ok) {
    list.innerHTML = "";
    const li = document.createElement("li");
    li.className = "comments-list-message";
    li.textContent = errorMessage(response, "Failed to load comments");
    list.appendChild(li);
    return;
    }

    const comments = response.data || [];
    list.innerHTML = "";

    if (comments.length === 0) {
    const li = document.createElement("li");
    li.className = "comments-list-message";
    li.textContent = "No comments yet";
    list.appendChild(li);
    return;
    }

    comments.forEach((comment) => {
    const li = document.createElement("li");
    li.className = "comment-item";

    const authorElement = document.createElement("p");
    authorElement.className = "comment-author";
    authorElement.textContent = comment.username;

    const contentElement = document.createElement("p");
    contentElement.className = "comment-content";
    contentElement.textContent = comment.content;

    li.append(authorElement, contentElement);
    list.appendChild(li);
    });
}