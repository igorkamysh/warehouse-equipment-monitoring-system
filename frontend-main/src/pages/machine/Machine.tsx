import React, { useState, useEffect } from "react";
import MachineInfo, {
  IMachineInfo,
} from "../../components/machines/MachineInfo";
import ControlButtons from "../../components/machines/ControlButtons";
import { useNavigate, useSearchParams } from "react-router-dom";
import { MachineState } from "../../components/machines/types";
import { useConfig } from "../../ConfigContext";

import Cookies from "js-cookie";
import axios from "axios";

const Machine: React.FC = () => {
  const config = useConfig();
  const navigate = useNavigate();

  const [searchParams] = useSearchParams();
  const machineId = searchParams.get("id");

  const [info, setInfo] = useState<IMachineInfo | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const token = Cookies.get("authToken");

  useEffect(() => {
    console.log("machine info fetching");
    const fetchMachineInfo = async () => {
      try {
        if (!token) {
          navigate("/login");
        }

        const response = await axios.get(
          `${config.api.backend_url}/get_machine?machine_id=${machineId}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "application/json",
            },
          }
        );

        console.log("response", response);
        setInfo(response.data);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch machines data");
        console.log(err);
        setLoading(false);
      }
    };

    fetchMachineInfo();
  }, [navigate, config, loading, machineId, token]);

  const sendUnlockMachine = async (machineId: string) => {
    try {
      const url = `${config.api.backend_url}/unlock_machine`;
      const data = {
        machine_id: machineId,
      };

      const headers = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      };

      const response = await axios.post(url, data, { headers });
      console.log(response);
      console.log(info);

      if (info !== null) {
        info.state = MachineState.InUse;
        setInfo(info);
      }
      setLoading(true);
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
  };

  const sendStopMachine = async (machineId: string) => {
    try {
      const url = `${config.api.backend_url}/stop_machine`;
      const data = {
        machine_id: machineId,
      };

      const headers = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      };

      const response = await axios.post(url, data, { headers });
      console.log(response);
      console.log(info);

      if (info !== null) {
        info.state = MachineState.Stopped;
        setInfo(info);
      }
      setLoading(true);
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
  };

  const sendUnstopMachine = async (machineId: string) => {
    try {
      const url = `${config.api.backend_url}/unstop_machine`;
      const data = {
        machine_id: machineId,
      };

      const headers = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      };

      const response = await axios.post(url, data, { headers });
      console.log(response);
      console.log(info);

      if (info !== null) {
        info.state = MachineState.InUse;
        setInfo(info);
      }
      setLoading(true);
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
  };

  const sendLockMachine = async (machineId: string) => {
    try {
      const url = `${config.api.backend_url}/lock_machine`;
      const data = {
        machine_id: machineId,
      };

      const headers = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      };

      const response = await axios.post(url, data, { headers });
      console.log(response);

      if (info !== null) {
        info.state = MachineState.Free;
        setInfo(info);
      }
      setLoading(true);
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
  };

  const handleStart = () => {
    if (!machineId) {
      return <p className="text-center text-gray-500">Wrong machine 'id'</p>;
    }
    sendUnlockMachine(machineId);
  };

  const handleStop = () => {
    if (!machineId) {
      return <p className="text-center text-gray-500">Wrong machine 'id'</p>;
    }
    sendStopMachine(machineId);
  };

  const handleUnstop = () => {
    if (!machineId) {
      return <p className="text-center text-gray-500">Wrong machine 'id'</p>;
    }
    sendUnstopMachine(machineId);
  };

  const handleFinish = () => {
    if (!machineId) {
      return <p className="text-center text-gray-500">Wrong machine 'id'</p>;
    }
    sendLockMachine(machineId);
  };

  if (token === undefined) {
    navigate("/login");
    return;
  }

  if (!machineId) {
    return <p>Machine Id is incorrect</p>;
  }

  if (loading || info === null) {
    return <p className="text-center text-gray-500">Loading...</p>;
  }

  if (error) {
    return <p className="text-center text-red-500">{error}</p>;
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="w-full max-w-md">
        <MachineInfo info={info} />
        <ControlButtons
          machineState={info?.state}
          onStart={handleStart}
          onStop={handleStop}
          onUnstop={handleUnstop}
          onFinish={handleFinish}
        />
      </div>
    </div>
  );
};

export default Machine;
