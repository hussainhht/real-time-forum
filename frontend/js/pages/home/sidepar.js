import {escapeHtml} from '../../api.js';

export function renderSidebar(users = [], selectedUserId = null) {


    return `
    <aside class="sidebar">
      <div class="chat-sidebar-header">
        <div class="chat-sidebar-avatar"></div>
        <div class="chat-sidebar-title">Chats</div>
      </div>

      <div class="chat-users-list">
        ${
          users.length
            ? users
                .map((user) => {
                  const isActive =
                    user.id === selectedUserId ? "active" : "";

                  return `
                    <button
                      class="chat-user-item ${isActive}"
                      data-user-id="${user.id}"
                    >
                      <span class="chat-user-circle"></span>

                      <span class="chat-user-name">
                        ${escapeHtml(user.username)}
                      </span>
                    </button>
                  `;
                })
                .join("")
            : `<p class="chat-empty">No users yet</p>`
        }
      </div>
    </aside>` ;


    
}