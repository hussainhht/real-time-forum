import {escapeHtml} from "../../api.js";

export function renderTopbar(currentUser) {
  return `
    <header class="home-topbar">
      <div class="topbar-user-box">
        <div class="topbar-avatar"></div>

        <div class="topbar-user-info">
          <small>@${escapeHtml(currentUser.username )}</small>
          <strong>
            ${escapeHtml(currentUser.first_name )}
            ${escapeHtml(currentUser.last_name )}
          </strong>
        </div>

        <button type="button" id="logout-button" class="logout-btn">
          Logout
        </button>
      </div>

      <button type="button" id="create-post-button" class="create-post-btn">
        Create +
      </button>
    </header>
  `;
}