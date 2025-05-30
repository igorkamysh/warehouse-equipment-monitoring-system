import yaml from "js-yaml";
import fs from "fs";

export interface ApiConfig {
  backend_url: string;
}

export interface AppConfig {
  api: ApiConfig;
}

export const loadConfigOld = (): AppConfig => {
  try {
    const configPath = "config.yml";
    const fileContents = fs.readFileSync(configPath, "utf8");
    const config = yaml.load(fileContents) as AppConfig;

    return config;
  } catch (e) {
    console.error("Ошибка при загрузке конфигурации:", e);
    throw e;
  }
};

export const loadConfig = (): AppConfig => {
  return {
    api: {
      /* Change ip later */
      backend_url: "http://192.168.113.215:8080",
    },
  };
};

/*
export const loadConfig = async (): Promise<AppConfig> => {
  try {
    const response = await fetch("public/config.yml");
    if (!response.ok) {
      throw new Error("Ошибка при загрузке конфигурационного файла");
    }
    const text = await response.text();
    const config = yaml.load(text) as AppConfig;
    return config;
  } catch (e) {
    console.error("Ошибка при загрузке конфигурации:", e);
    throw e;
  }
};

  */
