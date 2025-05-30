import React, { ReactNode } from "react";
import { MachineState } from "./types";

interface ControlButtonsProps {
  onStart: () => void;
  onStop: () => void;
  onUnstop: () => void;
  onFinish: () => void;

  machineState: MachineState;
}

const ControlButtons: React.FC<ControlButtonsProps> = ({
  onStart,
  onStop,
  onUnstop,
  onFinish,
  machineState,
}) => {
  const renderPauseButton = (): ReactNode => {
    if (machineState === MachineState.Stopped) {
      return (
        <button
          disabled={false}
          onClick={onUnstop}
          className={
            "focus:outline-none text-white bg-yellow-400 hover:bg-yellow-500 focus:ring-4 focus:ring-yellow-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:focus:ring-yellow-900"
          }
        >
          Возобновить
        </button>
      );
    } else if (machineState == MachineState.Free) {
      return (
        <button
          disabled={true}
          onClick={onStop}
          className={
            "py-2.5 px-5 me-2 mb-2 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700"
          }
        >
          Пауза
        </button>
      );
    } else {
      return (
        <button
          disabled={false}
          onClick={onStop}
          className={
            "focus:outline-none text-white bg-yellow-400 hover:bg-yellow-500 focus:ring-4 focus:ring-yellow-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:focus:ring-yellow-900"
          }
        >
          Пауза
        </button>
      );
    }
  };

  return (
    <>
      <div className="flex flex-col w-48 mx-auto space-y-2">
        <button
          disabled={machineState === MachineState.InUse}
          onClick={onStart}
          className={
            machineState === MachineState.Free
              ? "focus:outline-none text-white bg-green-700 hover:bg-green-800 focus:ring-4 focus:ring-green-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-green-600 dark:hover:bg-green-700 dark:focus:ring-green-800"
              : "py-2.5 px-5 me-2 mb-2 text-sm font-medium text-gray-900 bg-white rounded-lg border border-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600"
          }
        >
          Начать
        </button>
        {renderPauseButton()}
        <button
          disabled={machineState !== MachineState.InUse}
          onClick={onFinish}
          className={
            machineState === MachineState.InUse
              ? "focus:outline-none text-white bg-red-700 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-900"
              : "py-2.5 px-5 me-2 mb-2 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700"
          }
        >
          Завершить
        </button>
      </div>
    </>
  );
};

export default ControlButtons;
