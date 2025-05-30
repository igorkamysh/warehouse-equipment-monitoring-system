import React, { useState, useEffect } from "react";
import axios from "axios";
import { useSearchParams, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { useConfig } from "../../ConfigContext";


const QrCode: React.FC = () => {
  const config = useConfig();
  const navigate = useNavigate();

  const [error, setError] = useState<string | null>(null);
  const [searchParams] = useSearchParams();
  const [loading] = useState<boolean>(true);

  const token = Cookies.get("authToken");

  useEffect(() => {
    const tryEndSession = async () => {
      try {
        console.log("config: ", config);
        console.log("key: ", searchParams.get('key'));
        console.log("token: ", token);
        console.log("parking_name: ", searchParams.get('parking_name'));

        const url = `${config.api.backend_url}/finish_session`;
        const data = {
          key: searchParams.get("key"),
          parking_name: searchParams.get("parking_name")
        };

        const headers = {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        };

        const response = await axios.post(url, data, { headers });
        console.log("response: ", response);

        if (response.data.msg) {
          navigate("/machines");
        } else if (response.data.error) {
          console.error("End session failed:", response.data.error);
        }
      } catch (error) {
        if (axios.isAxiosError(error)) {
          // Проверяем, есть ли у ошибки ответ от сервера
          if (error.response) {
            const { status, data } = error.response;
  
            // Устанавливаем сообщение об ошибке в стейт
            if (data && data.error) {
              setError(data.error); // Сообщение от сервера
            } else {
              setError(`Request failed with status code ${status}`); // Общая ошибка
            }
  
            console.error(`Error: ${status} - ${data?.error || "Unknown error"}`);
          } else if (error.request) {
            // Ошибка при отправке запроса (нет ответа)
            setError(
              "No response received from the server. Please try again later."
            );
            console.error("No response from server:", error.request);
          } else {
            // Ошибка при настройке запроса
            setError(`Request setup error: ${error.message}`);
            console.error("Request error:", error.message);
          }
        } else {
          // Любая другая ошибка
          setError("An unknown error occurred.");
          console.error("Unknown error:", error);
        }
      }
    }

    tryEndSession();
  }, [navigate, loading]);

  if (token === undefined) {
    navigate("/login");
    return;
  }

  if (error) {
    return <p className="text-center text-red-500">{error}</p>;
  }

  return('');
};

export default QrCode;
