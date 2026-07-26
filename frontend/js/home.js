import { renderPage } from "./router.js";

function homepage(app) {
  app.innerHTML = `

        <button id="logout-button" type="button">Logout</button>



    
    `;

  const logoutButton = document.getElementById("logout-button");
  logoutButton.addEventListener("click", async () => {
    const respons = await fetch("/logout", {
      method: "POST",
      credentials: "include",
    });

    if (!respons.ok) {
      return;
    }
    renderPage("login");
    

  });
}
