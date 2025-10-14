import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import App from "./App";
import RegionalInfo from "./pages/RegionalInfo";
import "./App.css";
import HowItWorks from "./pages/HowItWorks";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/regional-info" element={<RegionalInfo />} />
        <Route path="/how-it-works" element={<HowItWorks />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
