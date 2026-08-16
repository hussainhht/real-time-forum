import { escapeHtml, apiFetch, errorMessage } from "../../api.js";
import { renderMessagePage, setupMessageForm } from "./message.js";

const usersWithUnread = new Set();

export function markUserUnread(userId) {
  usersWithUnread.add(Number(userId));
  updateUnreadBadge(userId);
}

export function clearUserUnread(userId) {
  usersWithUnread.delete(Number(userId));
  updateUnreadBadge(userId);
}

function updateUnreadBadge(userId) {
  const button = document.querySelector(
    `.chat-user-btn[data-user-id="${Number(userId)}"]`,
  );
  if (!button) return;

  const info = button.querySelector(".chat-user-info");
  if (!info) return;

  const existingBadge = info.querySelector(".chat-user-badge-new");
  if (usersWithUnread.has(Number(userId))) {
    if (!existingBadge) {
      info.insertAdjacentHTML("beforeend", '<span class="chat-user-badge-new">NEW</span>');
    }
  } else if (existingBadge) {
    existingBadge.remove();
  }
}

export function renderSidebar(users = [], selectedUserId = null) {
  return `
    <aside class="chat-sidebar">
      <div class="chat-sidebar-header">
        <div class="pages-header">
          <button id="home-btn" type="button">Home</button>
        </div>
        <h2>Chats</h2>
      </div>

      <div id="chat-users-list" class="chat-users-list">
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
    if (message)
      message.textContent = errorMessage(response, "Failed to load users");
    return;
  }

  const users = Array.isArray(response.data) ? response.data : [];
  renderUsersList(users, usersList);
}

function getInitials(user) {
  const firstInitial = (user.first_name || "").charAt(0);
  const lastInitial = (user.last_name || "").charAt(0);
  const initials = `${firstInitial}${lastInitial}`;
  if (initials) {
    return initials.toUpperCase();
  }
  return (user.username || "").charAt(0).toUpperCase();
}

function renderUsersList(users, usersList) {
  if (users.length === 0) {
    usersList.innerHTML = "<p>No users found</p>";
    return;
  }
  usersList.innerHTML = users
    .map(
      (user) => `
            <button type="button" class="chat-user-btn" data-user-id="${escapeHtml(user.id)}" data-username="${escapeHtml(user.username)}" data-online="${user.online ? "true" : "false"}">
            <span class="chat-user-circle ${user.online ? "online" : "offline"}">${escapeHtml(getInitials(user))}</span>

            <span class="chat-user-info">
              <span class="chat-user-name">${escapeHtml(user.username)}</span>
              ${usersWithUnread.has(Number(user.id)) ? '<span class="chat-user-badge-new">NEW</span>' : ""}
            </span>
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
      const isOnline = button.dataset.online === "true";

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
        online: isOnline,
      };

      clearUserUnread(userId);

      homeContainer.innerHTML = renderMessagePage(user);

      await setupMessageForm(userId);
    });
  });
}
