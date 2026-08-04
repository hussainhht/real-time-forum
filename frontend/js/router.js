import { LoginPage } from "./pages/login.js";
import {RegisterPage} from "./pages/register.js";
import { homepage } from "./pages/home";
import {getCurrentUser} from "./auth.js";

export const app = document.getElementById("app");
let currentUser = null;

export function renderPage(page) {

  const protectedPages = ["home"]
  const guestPages =["login","register"]

  if (protectedPages.includes(page) && !currentUser) {
    LoginPage(app);
    return;
  }
  if (guestPages.includes(page) && currentUser) {
    homepage(app,currentUser);
    return;
    
  }

  if (page === "login") {
    LoginPage(app);
    return;
  }

  if (page === "register") {
    RegisterPage(app);
    return;
  }
  
  if (page === "home") {
    homepage(app)
    return
    
  }

  if (currentUser) {
    renderPage("home");
    return;
    
  }

  LoginPage(app);
}


// it can move to auth.js
export async function startApp() {
  currentUser = await getCurrentUser();

  if (currentUser) {
    renderPage("home");
    return;
  }

  renderPage("login");
  
}

export function setCurrentUser(user) { 
  currentUser = user ;
}

export function clearCurrentUser(){
  currentUser = null;
}