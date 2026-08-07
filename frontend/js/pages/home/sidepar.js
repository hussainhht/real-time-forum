import { escapeHtml, apiFetch, errorMessage } from "../../api.js";
import { renderMessagePage, setupMessageForm } from "./message.js";

export function renderSidebar(users = [], selectedUserId = null) {
  return `
        <aside class="chat-sidebar">
            <div class="chat-sidebar-header">
                <h2>Chats</h2>
            </div>

            <div
                id="chat-users-list"
                class="chat-users-list"
            >
                <p>Loading users...</p>
            </div>

            <p class="chat-sidebar-footer" id="message"></p>

        </aside>

    `;
}

export async function updateSidebarUsers() {
  const usersList = document.getElementById("chat-users-list");
  const message = document.getElementById("message");

  if (!usersList) {
    console.error("Users list element not found");
    return;
  }

  const response = await apiFetch("/api/chat-users", {
    method: "GET",
  });

  if (!document.body.contains(usersList)) {
    return;
  }

  if (!response.ok) {
    if (message) message.textContent = errorMessage(response, "Failed to load users");
    return;
  }

  const users = Array.isArray(response.data) ? response.data : [];
  renderUsersList(users, usersList);
}

function renderUsersList(users, usersList) {
  if (users.length === 0) {
    usersList.innerHTML = "<p>No users found</p>";
    return;
  }
  usersList.innerHTML = users
    .map(
      (user) => `
            <button type="button" class="chat-user-btn" data-user-id="${escapeHtml(user.id)}" data-username="${escapeHtml(user.username)}">
            <span class="chat-user-circle"></span>

            <span class="chat-user-name">${escapeHtml(user.username)}</span>
            </button>
        `,
    )
    .join("");

  setupChatUserEvents();
}

function setupChatUserEvents() {
  const userButtons = document.querySelectorAll(".chat-user-btn");

  userButtons.forEach((button) => {
    button.addEventListener("click", async () => {
      const userId = Number(button.dataset.userId);
      const username = button.dataset.username;

      if (!Number.isInteger(userId) || userId <= 0) {
        console.error("Invalid user ID");
        return;
      }

      const homeContainer = document.getElementById("home-container");
      if (!homeContainer) {
        console.error("Home container not found");
        return;
      }

      const user = {
        id: userId,
        username: username,
      };

      homeContainer.innerHTML = renderMessagePage(user);

      await setupMessageForm(userId);
    });
  });
}