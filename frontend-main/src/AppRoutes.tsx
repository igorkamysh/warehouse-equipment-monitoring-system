import React from "react";
import App from "./App";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Machines from "./pages/machines/Machines";
import Login from "./pages/login/Login";
import Machine from "./pages/machine/Machine";
import QrCode from "./pages/qrCode/QrCode";

const AppRoutes: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" Component={App} />
        <Route path="/machines" Component={Machines} />
        <Route path="/login" Component={Login} />
        <Route path="/machine" Component={Machine} />
        <Route path="/finish_session" Component={QrCode} />
      </Routes>
    </BrowserRouter>
  );
};

export default AppRoutes;
