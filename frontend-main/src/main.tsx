import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import AppRoutes from "./AppRoutes";
import { ConfigProvider } from "./ConfigContext";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ConfigProvider>
      <AppRoutes />
    </ConfigProvider>
  </StrictMode>
);
