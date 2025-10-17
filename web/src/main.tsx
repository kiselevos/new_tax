import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import App from "./App";
import RegionalInfo from "./pages/RegionalInfo";
import "./App.css";
import HowItWorks from "./pages/HowItWorks";
import Results from "./pages/Results";
import SpecialTaxModes from "./pages/SpecialTaxModes";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/regional-info" element={<RegionalInfo />} />
        <Route path="/special-tax-modes" element={<SpecialTaxModes />} />
        <Route path="/how-it-works" element={<HowItWorks />} />
        <Route path="/results" element={<Results />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
