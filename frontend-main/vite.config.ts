import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: "0.0.0.0", // Позволяет прослушивать запросы со всех IP-адресов
    port: 3000, // Задает номер порта, на котором будет запущен сервер
    strictPort: true, // Если true, сервер не будет пытаться найти доступный порт, если указанный занят
  },
});
