import { useState } from "react";
import "./App.css";
import { taxClient } from "./client";
import { Link } from "react-router-dom";
import { CalculatePrivateResponse } from "./gen/api/tax_pb";

interface FormState {
  grossSalary: string;
  territorialMultiplier: string;
  northernCoefficient: string;
  startDate: string;
  hasTaxPrivilege: boolean;
  isNotResident: boolean;
}

const months = [
  { value: "01", label: "Январь" },
  { value: "02", label: "Февраль" },
  { value: "03", label: "Март" },
  { value: "04", label: "Апрель" },
  { value: "05", label: "Май" },
  { value: "06", label: "Июнь" },
  { value: "07", label: "Июль" },
  { value: "08", label: "Август" },
  { value: "09", label: "Сентябрь" },
  { value: "10", label: "Октябрь" },
  { value: "11", label: "Ноябрь" },
  { value: "12", label: "Декабрь" },
];

const currentYear = new Date().getFullYear();

export default function App() {
  const [form, setForm] = useState<FormState>({
    grossSalary: "",
    territorialMultiplier: "",
    northernCoefficient: "",
    startDate: "",
    hasTaxPrivilege: false,
    isNotResident: false,
  });

  const [result, setResult] = useState<CalculatePrivateResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // универсальный обработчик для всех полей
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target;
    const checked =
      e.target instanceof HTMLInputElement && e.target.type === "checkbox"
        ? e.target.checked
        : undefined;

    setForm((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked ?? false : value,
    }));
  };

  const formatDate = (month: string) => `${currentYear}-${month}-01`;

  const handleMonthSelect = (month: string) => {
    setForm((prev) => ({
      ...prev,
      startDate: formatDate(month),
    }));
  };

  // обработка ввода оклада
  const handleSalaryInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value.replace(",", ".").replace(/[^0-9.]/g, "");
    const normalized =
      raw.split(".").length > 2 ? raw.replace(/\.+$/, "") : raw;

    setForm((prev) => ({
      ...prev,
      grossSalary: normalized,
    }));
  };

  // форматирование оклада при потере фокуса
  const handleSalaryBlur = () => {
    const val = form.grossSalary.replace(",", ".").replace(/\s/g, "");
    if (!val) return;

    const num = parseFloat(val);
    if (isNaN(num)) return;

    const formatted = num.toLocaleString("ru-RU", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });

    setForm((prev) => ({
      ...prev,
      grossSalary: formatted,
    }));
  };

  // отправка формы
  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    const salary = parseFloat(
      form.grossSalary.replace(/\s/g, "").replace(",", ".")
    );

    if (isNaN(salary) || salary <= 0) {
      setError("Некорректное значение оклада");
      setLoading(false);
      return;
    }

    try {
  const res = await taxClient.calculatePrivate({
    grossSalary: BigInt(Math.round(salary * 100)), // ✅ BigInt
    territorialMultiplier: form.territorialMultiplier
      ? BigInt(Math.round(parseFloat(form.territorialMultiplier) * 100))
      : BigInt(0),
    northernCoefficient: form.northernCoefficient
      ? BigInt(Number(form.northernCoefficient) + 100)
      : BigInt(0),
    startDate: form.startDate
      ? {
          seconds: BigInt(
            Math.floor(new Date(form.startDate).getTime() / 1000)
          ),
        }
      : undefined,
    hasTaxPrivilege: form.hasTaxPrivilege,
    isNotResident: form.isNotResident,
  });

  setResult(res);
} catch (err) {
  if (err instanceof Error) {
    setError(err.message);
  } else {
    setError("Ошибка при запросе");
  }
} finally {
  setLoading(false);
}

  }

  return (
    <div className="container">
      <h1>Расчёт налога на текущий год</h1>
      <h3>
        <Link to="/how-it-works" className="field-link under-title">
          (Как это работает)
        </Link>
      </h3>

      <form onSubmit={handleSubmit} className="form">
        {/* --- оклад --- */}
        <div className="field-row">
          <div className="field-left">
            <label htmlFor="grossSalary">Оклад:</label>
          </div>
          <div className="field-right">
            <div className="salary-field">
              <input
                type="text"
                name="grossSalary"
                id="grossSalary"
                inputMode="decimal"
                value={form.grossSalary}
                onChange={handleSalaryInput}
                onBlur={handleSalaryBlur}
                required
                placeholder="Введите сумму"
              />
              <span className="ruble">₽</span>
            </div>
          </div>
        </div>

        {/* --- территориальный коэффициент --- */}
        <div className="field-row">
          <div className="field-left">
            <label htmlFor="territorialMultiplier">
              Территориальный коэффициент:
            </label>
            <Link to="/regional-info" className="field-link under-title">
              справка →
            </Link>
          </div>
          <div className="field-right">
            <select
              id="territorialMultiplier"
              name="territorialMultiplier"
              value={form.territorialMultiplier}
              onChange={handleChange}
              className="field-select"
            >
              <option value=""> --без коэффициента-- </option>
              {Array.from({ length: 21 }, (_, i) =>
                (1 + i * 0.05).toFixed(2)
              ).map((val) => (
                <option key={val} value={val}>
                  {val}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* --- северная надбавка --- */}
        <div className="field-row">
          <div className="field-left">
            <label htmlFor="northernCoefficient">Северная надбавка:</label>
            <Link to="/regional-info" className="field-link under-title">
              справка →
            </Link>
          </div>
          <div className="field-right">
            <select
              id="northernCoefficient"
              name="northernCoefficient"
              value={form.northernCoefficient}
              onChange={handleChange}
              className="field-select"
            >
              <option value=""> --без надбавки-- </option>
              {Array.from({ length: 10 }, (_, i) => (i + 1) * 10).map(
                (val) => (
                  <option key={val} value={val}>
                    {val}%
                  </option>
                )
              )}
            </select>
          </div>
        </div>

        {/* --- дата начала расчёта --- */}
        <div className="field-row">
          <div className="field-left">
            <label>Дата начала расчёта:</label>
          </div>
          <div className="field-right">
            <div className="date-field">
              <span>01</span>
              <select
                name="startDate"
                value={
                  form.startDate
                    ? form.startDate.split("-")[1]
                    : String(new Date().getMonth() + 1).padStart(2, "0")
                }
                onChange={(e) => handleMonthSelect(e.target.value)}
              >
                {months.map((m) => (
                  <option key={m.value} value={m.value}>
                    {m.label}
                  </option>
                ))}
              </select>
              <span>{currentYear}</span>
            </div>
          </div>
        </div>

        {/* --- чекбоксы --- */}
        <div className="field-row">
          <div className="field-left">
            <label>Параметры:</label>
          </div>
          <div className="field-right">
            <div className="checkbox-column">
              <label className="checkbox-inline">
                <input
                  type="checkbox"
                  name="hasTaxPrivilege"
                  checked={form.hasTaxPrivilege}
                  onChange={handleChange}
                />
                <span>Льготы для силовых структур</span>
              </label>
              <label className="checkbox-inline">
                <input
                  type="checkbox"
                  name="isNotResident"
                  checked={form.isNotResident}
                  onChange={handleChange}
                />
                <span>Налоговый нерезидент</span>
              </label>
            </div>
          </div>
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Расчёт..." : "Рассчитать"}
        </button>
      </form>

      {error && <p className="error">{error}</p>}

      {result && (
        <div className="result">
          <h2>Результат</h2>
          <p>Годовой доход: {result.annualGrossIncome?.toString()}</p>
          <p>Налог: {result.annualTaxAmount?.toString()}</p>
          <p>Чистый доход: {result.annualNetIncome?.toString()}</p>
        </div>
      )}
    </div>
  );
}
