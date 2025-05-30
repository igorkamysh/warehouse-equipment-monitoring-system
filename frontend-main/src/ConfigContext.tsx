import React, { createContext, useContext } from "react";
import { AppConfig, loadConfig } from "./config";

const ConfigContext = createContext<AppConfig | null>(null);

export const ConfigProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const config = loadConfig();

  return (
    <ConfigContext.Provider value={config}>{children}</ConfigContext.Provider>
  );
};

// Хук для удобного использования контекста конфигурации
export const useConfig = (): AppConfig => {
  const context = useContext(ConfigContext);
  if (!context) {
    throw new Error("useConfig must be used within a ConfigProvider");
  }
  return context;
};
