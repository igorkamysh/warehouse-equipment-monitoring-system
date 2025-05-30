import React, { useState, useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";
import Cookies from "js-cookie";
import axios from "axios";
import "./index.css";
import { useConfig } from "../../ConfigContext";

interface Machine {
  id: string;
  state: number;
  voltage: number;
  ipAddr: string;
}

const MachinesTable: React.FC = () => {
  const config = useConfig();
  const navigate = useNavigate();

  const [machines, setMachines] = useState<Machine[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Fetch machines data
    const fetchMachines = async () => {
      try {
        const token = Cookies.get("authToken");
        if (!token) {
          navigate("/login");
        }
        const response = await axios.get(
          `${config.api.backend_url}/get_all_machines`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "application/json",
            },
          }
        );

        setMachines(response.data);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch machines data");
        console.log(err);
        setLoading(false);
      }
    };

    fetchMachines();
  }, [navigate, config]);

  if (loading) {
    return <p className="text-center text-gray-500">Loading...</p>;
  }

  if (error) {
    return <p className="text-center text-red-500">{error}</p>;
  }

  return (
    <div className="table-container relative overflow-x-auto shadow-md sm:rounded-lg">
      <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
        <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
          <tr>
            <th scope="col" className="px-6 py-3">
              Machine Id
            </th>
            <th scope="col" className="px-6 py-3">
              State
            </th>
            <th scope="col" className="px-6 py-3">
              Voltage
            </th>
            <th scope="col" className="px-6 py-3">
              Ip Address
            </th>
          </tr>
        </thead>
        <tbody>
          {machines.map((machine) => (
            <tr
              key={machine.id}
              className="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700"
            >
              <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                <Link
                  className="underline text-blue-500 hover:text-blue-700 cursor-pointer"
                  to={`/machine?id=${machine.id}`}
                >
                  {machine.id}
                </Link>
              </td>
              {machine.state === 0 ? (
                <td className="px-6 py-4 flex flex-row gap-2 items-center text-green-600">
                  <div className="circle bg-green-600"></div>
                  <div>Free</div>
                </td>
              ) : (
                <td className="px-6 py-4 flex flex-row gap-2 items-center text-red-400">
                  <div className="circle bg-red-400"></div>
                  <div>In Use</div>
                </td>
              )}
              <td className="px-6 py-4">{machine.voltage} V</td>
              <td className="px-6 py-4">{machine.ipAddr}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default MachinesTable;
