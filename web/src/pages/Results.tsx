import { useLocation, useNavigate } from "react-router-dom";
import MainHeader from "./MainHeader";
import type { LocationState, ResultData } from "./types/tax.types";
import MonthlyAccordion from "./MonthlyAccordion";
import "../App.css";

export default function Results() {
  const location = useLocation();
  const navigate = useNavigate();

  const state = location.state as LocationState | undefined;
  let result = state?.result;
  const formData = state?.formData;

  console.group("DEBUG: Results.tsx data inspection");
  console.log("Raw location.state:", location.state);
  console.log("Result type:", typeof result);
  console.log("Result object:", result);
  if (result) {
    console.log("Result keys:", Object.keys(result));
    console.log("monthlyDetails:", result.monthlyDetails);
  }
  console.log("Form data:", formData);
  console.groupEnd();

  const currentYear = new Date().getFullYear();

  // Если result пришёл как строка — пробуем распарсить
  if (typeof result === "string") {
    try {
      result = JSON.parse(result) as ResultData;
    } catch (e) {
      console.error("Ошибка парсинга result:", e);
    }
  }

  // Если данных нет
  if (!result || !formData) {
    return (
      <div className="container">
        <h2>Данные не найдены</h2>
        <p>Пожалуйста, вернитесь на главную страницу и выполните расчет.</p>
        <button onClick={() => navigate("/")} className="back-button">
          Вернуться к калькулятору
        </button>
      </div>
    );
  }

  const monthlyDetails = result.monthlyDetails || [];

  const handleNewCalculation = () => navigate("/");

  return (
    <div className="container">
      {/* Заголовок */}
      <div className="results-header">
        <h1>Результаты расчета</h1>
        <button onClick={handleNewCalculation} className="back-button">
          Новый расчет
        </button>
      </div>

      {/* Основная сводка */}
      <MainHeader result={result} formData={formData} currentYear={currentYear} />

      {/* Помесячная детализация */}
      <MonthlyAccordion monthlyDetails={monthlyDetails} result={result} />

      {/* Кнопка для нового расчета */}
      <div className="actions">
        <button onClick={handleNewCalculation} className="primary-button">
          Сделать новый расчет
        </button>
      </div>
    </div>
  );
}