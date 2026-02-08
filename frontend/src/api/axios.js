import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080", // backend actuel
  withCredentials: true, // utile si cookies / auth plus tard
});

export default api;
