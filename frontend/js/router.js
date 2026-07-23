import { LoginPage } from "./login.js";
import {RegisterPage} from "./register.js"
// import { app } from "./main";

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

  LoginPage(app);
}
