import { escapeHtml, apiFetch, errorMessage } from "../../api.js";

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

  if (!document.body.contains(message)) {
    console.error("Message element is no longer in the DOM");
    return;
  }

  if (!response.ok) {
    message.textContent = errorMessage(response, "Failed to load users");  
    return;
  }

  const users = Array.isArray(response.data) ? response.data : [];
  renderChatUsers(users, usersList);
}

function renderChatUsers(users, usersList) {
  if (users.length === 0) {
    usersList.innerHTML = "<p>No users found</p>";
    return;
  }
  usersList.innerHTML = users
    .map(
      (user) => `
            <button type="button" class="chat-user-btn" data-user-id="${escapeHtml(user.id)}">
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
    button.addEventListener("click", () => {
      const userId = Number(button.dataset.userId);
      if (!Number.isInteger(userId) || userId <= 0) {
        console.error("Invalid user ID");
        return;
      }
      loadMessagesForUser(userId);
    });
  });
}

async function loadMessagesForUser(userId) {
  const messagesContainer = document.getElementById("messages-container");

  if (!messagesContainer) {
    console.error("Messages container not found");
    return;
  }

  messagesContainer.textContent = "Loading messages...";

  const response = await apiFetch(`/api/messages/${userId}`);

  if (!response.ok) {
    messagesContainer.innerHTML = `<p class="chat-error">${errorMessage(response, "Failed to load messages")}</p>`;
    return;
  }
  const messages = Array.isArray(response.data) ? response.data : [];

  renderMessages(messages, messagesContainer);
}

function renderMessages(messages, messagesContainer) {
  if (messages.length === 0) {
    messagesContainer.innerHTML = `<p class="empty-chat">No messages found</p>`;
    return;
  }

  messagesContainer.innerHTML = messages
    .map(
      (message) => `
            <div class="chat-message">
                <span class="chat-message-sender">${escapeHtml(message.sender_id)}</span>
                <span class="chat-message-content">${escapeHtml(message.content)}</span>
                <span class="chat-message-timestamp">${escapeHtml(message.created_at)}</span>
            </div>
        `,
    )
    .join("");
}
