import {escapeHtml} from '../../api.js';

export function renderErrorPage({ status, message }) {
    const app = document.getElementById("app");
    
    if (!app) {
        console.error('App element not found');
        return;
    }

    app.innerHTML = `
    <main class = "error-page">
      <h1 class="error-status">
        ${escapeHtml(String(status))} 
      </h1>

      <p class="error-message">
        ${escapeHtml(message)}
      </p>

      <div class="error-actions">
        <button type="submit" id="error-home-button">
          back to home
        </button>
      </div>

    </main>`

    const homeButton = document.getElementById("error-home-button");

    homeButton.addEventListener("click", () => {
        window.location.href = "/";
    })

}

