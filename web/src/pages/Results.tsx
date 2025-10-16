import { useLocation, useNavigate } from "react-router-dom";
import MainHeader from "./MainHeader";
import MonthlyAccordion from "./MonthlyAccordion";
import "../App.css";

interface ResultsFormState {
  grossSalary: string;
  territorialMultiplier: string;
  northernCoefficient: string;
  startDate: string;
  hasTaxPrivilege: boolean;
  isNotResident: boolean;
}

interface LocationState {
  result: any;
  formData: ResultsFormState;
}

export default function Results() {
  const location = useLocation();
  const navigate = useNavigate();

  let { result, formData } = (location.state as LocationState) || {};

  console.group("DEBUG: Results.tsx data inspection");
  console.log("Raw location.state:", location.state);
  console.log("Result type:", typeof result);
  console.log("Result object:", result);
  if (result) {
    console.log("Result keys:", Object.keys(result));
    console.log("monthlyDetails:", (result as any).monthlyDetails);
  }
  console.log("Form data:", formData);
  console.groupEnd();

  const currentYear = new Date().getFullYear();

  // Если result пришёл как строка — пробуем распарсить
  if (typeof result === "string") {
    try {
      result = JSON.parse(result);
    } catch (e) {
      console.error("Ошибка парсинга result:", e);
    }
  }

  // Если данных нет
  if (!result) {
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

  // Получаем массив деталей
  const monthlyDetails = (result as any).monthlyDetails || [];

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
