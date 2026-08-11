import { apiFetch, errorMessage } from "../../api.js";
import { navigateHome } from "../home.js";

export function renderCreatePostView(box) {
    box.innerHTML = `
        <h2>New post</h2>
        <form id="create-post-form">
            <input type="text" id="post-title" placeholder="Title" required />
            <select id="post-category" required>
                <option value="" disabled selected>Choose category</option>
                <option value="General">General</option>
                <option value="Technology">Technology</option>
                <option value="Sports">Sports</option>
                <option value="Entertainment">Entertainment</option>
                <option value="Science">Science</option>
            </select>
            <textarea id="post-content" placeholder="What's on your mind?" required></textarea>
            <button type="submit">Post</button>
        </form>
        <p id="create-post-message" class="form-message"></p>
    `;

    document
    .getElementById("create-post-form")
    .addEventListener("submit", handleCreatePost);

}

async function handleCreatePost(event) {
    event.preventDefault();

    const titleInput = document.getElementById("post-title");
    const categoryInput = document.getElementById("post-category");
    const contentInput = document.getElementById("post-content");
    const message = document.getElementById("create-post-message");

    const title = titleInput.value.trim();
    const category = categoryInput.value;
    const content = contentInput.value.trim();

    if (title === "" || category === "" || content === "") {
        message.textContent = "Please fill all fields";
        return;
    }

    const respons = await apiFetch("/posts", {
        method: "POST",
        body: { title, content, category },
    });

    if (!document.body.contains(message)) return;

    if (!respons.ok) {
        message.textContent = errorMessage(respons, "Failed to create post");
        return;
    }

    navigateHome("post", { postId: respons.data.id });
}