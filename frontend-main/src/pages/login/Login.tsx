import React, { useState } from "react";
import axios, { AxiosError } from "axios";
import { useNavigate } from "react-router-dom";
import Cookie from "js-cookie";
import { useConfig } from "../../ConfigContext";

interface ApiErrorResponse {
  error: string;
}

const Login: React.FC = () => {
  const config = useConfig();
  const navigate = useNavigate();

  const [phoneNumber, setPhoneNumber] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string | null>("");

  const fetchLogin = async () => {
    try {
      console.log("config: ", config);
      const response = await axios.post(`${config.api.backend_url}/login`, {
        phone_number: phoneNumber,
        password: password,
      });

      if (response.data.token) {
        Cookie.set("authToken", response.data.token);
        navigate("/machines");
      } else if (response.data.error) {
        setError(response.data.error);
        console.error("Login failed:", response.data.error);
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const axiosError = error as AxiosError<ApiErrorResponse>;
        const status = axiosError.response?.status;
        const message = axiosError.response?.data?.error || axiosError.message;

        if (status) {
          console.error(`Request failed with status ${status}`);
        }

        setError(message);
      } else {
        console.error("An unexpected error occurred:", error);
        setError("An unexpected error occurred");
      }
    }
  };
  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Phone Number:", phoneNumber);
    console.log("Password:", password);

    fetchLogin();
  };

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="p-6 max-w-md mx-auto bg-gray-800 rounded-xl shadow-md space-y-6">
        <h1 className="text-2xl font-bold text-center text-white">
          Вход в аккаунт
        </h1>

        {error && (
          <div
            className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
            role="alert"
          >
            <strong className="font-bold">Ошибка:</strong>
            <span className="block sm:inline"> {error}</span>
            <span
              className="absolute top-0 bottom-0 right-0 px-4 py-3"
              onClick={() => setError(null)}
            >
              <svg
                className="fill-current h-6 w-6 text-red-500"
                role="button"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
              >
                <title>Close</title>
                <path d="M14.348 5.652a1 1 0 00-1.414 0L10 8.586 7.066 5.652a1 1 0 10-1.414 1.414L8.586 10l-2.934 2.934a1 1 0 101.414 1.414L10 11.414l2.934 2.934a1 1 0 001.414-1.414L11.414 10l2.934-2.934a1 1 0 000-1.414z" />
              </svg>
            </span>
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label
              htmlFor="phone"
              className="block text-sm font-medium text-gray-300"
            >
              Номер телефона
            </label>
            <input
              type="tel"
              id="phone"
              value={phoneNumber}
              onChange={(e) => setPhoneNumber(e.target.value)}
              required
              className="mt-1 px-3 py-2 w-full border border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 bg-gray-700 text-white placeholder-gray-400"
              placeholder="Введите номер телефона"
            />
          </div>
          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-300"
            >
              Пароль
            </label>
            <input
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="mt-1 px-3 py-2 w-full border border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 bg-gray-700 text-white placeholder-gray-400"
              placeholder="Введите пароль"
            />
          </div>
          <button
            type="submit"
            className="w-full py-2 px-4 bg-sky-500 text-white font-semibold rounded-md shadow hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:ring-opacity-75"
          >
            Войти
          </button>
        </form>
      </div>
    </div>
  );
};

export default Login;
