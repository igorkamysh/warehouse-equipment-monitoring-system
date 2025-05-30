import "./App.css";
import { Link } from "react-router-dom";

function App() {
  return (
    <>
      <div>
        <p className="text-center text-3xl my-8">Smart Sharing</p>

        <h2 className="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
          Ссылки:
        </h2>
        <ul className="max-w-md space-y-1 text-gray-500 list-disc list-inside dark:text-gray-400">
          <li>
            <Link to="/login">Авторизация</Link>
          </li>
          <li>
            <Link to="/machines">Список техники</Link>
          </li>
        </ul>
      </div>
    </>
  );
}

export default App;
