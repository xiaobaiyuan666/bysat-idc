import axios from "axios";
import { ElMessage } from "element-plus";

const http = axios.create({
  baseURL: "/api",
  withCredentials: true,
});

http.interceptors.response.use(
  (response) => response,
  (error) => {
    const message =
      error.response?.data?.message ?? error.message ?? "请求失败，请稍后重试";

    if (error.response?.status !== 401) {
      ElMessage.error(message);
    }

    return Promise.reject(error);
  },
);

export { http };
