import axios from "axios";

const http = axios.create({
  baseURL: "/api/v1/admin",
  timeout: 10000
});

http.interceptors.request.use(config => {
  const token = localStorage.getItem("admin-token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default http;
