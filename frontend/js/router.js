import { LoginPage } from "./pages/login.js";
import {RegisterPage} from "./pages/register.js"
import { homepage } from "./pages/home.js";

export const app = document.getElementById("app");


export function renderPage(page) {
  if (page === "login") {
    LoginPage(app);
    return;
  }

  if (page === "register") {
    RegisterPage(app);
    return;
  }
  
  if (page === "home") {
    //todo: add home page here
    homepage(app)
    return
    
  }

  LoginPage(app);
}
