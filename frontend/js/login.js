import { renderPage } from "./router.js";

export function LoginPage(app) {
  app.innerHTML = ` 
    
    <form id="login-form">
        <h1>login</h1>
        <input type="text" id="identifier" placeholder="Email or username" required>
        <input type="password" id="password" placeholder="password" required>
        <button type="submit">Login</button>

    </form>
    <button type="button" id="go-register">Creat a new acount</button>
    
    `;

  const form = document.getElementById("login-form");
  const registerButton = document.getElementById("go-register");
  form.addEventListener("submit", handleLogin);

  registerButton.addEventListener("click", () => {
    renderPage("register");
  });
}

async function handleLogin(event) {
  event.preventDefault(); //* this for cansling the prouser defolt stings >> like spase tap things the man perpos here is to cansle the relode after submit the form .. in spa when the user refresh the page the html desper for a moment then when the user try to supmit the form it refrish but we can cansle this when we using this  in the coect spa this is must not happend ....

  const identifier = document.getElementById("identifier").value.trim();

  const password = document.getElementById("password").value.trim();

  if (identifier === "" || password === "") {
    console.log("this is not ok you need to fill all fields");
    return;
  }

  const response = await fetch("/login", {
    //* why we need asinc and awiat here becose the respons need time if we remove it the result will be <pending> soo we need to add await to wwait the sarver respons then conteno and there ways to the  Promise (pending,fulfilled,rejected) this will store in the varible if we dont use asinc and await >>
    method: "POST",

    //? this is wtifht heders in the page
    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      identifier,
      password,
    }),
  });

  if (!response.ok) {
    const errorMassage = await response.text();
    console.log("Login faild:", errorMassage);
    return;
  }

  console.log("Login successful");
}
