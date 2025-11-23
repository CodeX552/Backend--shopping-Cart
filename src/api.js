import axios from "axios";

const API = axios.create({
  baseURL: "http://backend-shopping-cart-8jog.onrender.com",
});

export function setAuthToken(token) {
  if (token) {
    API.defaults.headers.common["Authorization"] = `Bearer ${token}`;
  } else {
    delete API.defaults.headers.common["Authorization"];
  }
}

export default API;
