import React from "react";
import MachinesTable from "../../components/machines/MachinesTable";
import "./index.css";

const Machines: React.FC = () => {
  return (
    <div className="table-wrapper ">
      <div>
        <p className="text-slate-50 text-center text-2xl py-5">
          Машины в системе
        </p>
        <MachinesTable />
      </div>
    </div>
  );
};

export default Machines;
