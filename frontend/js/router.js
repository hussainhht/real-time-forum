import { LoginPage } from "./login.js";
import {RegisterPage} from "./register.js"
 import { homepage } from "./home.js";

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
