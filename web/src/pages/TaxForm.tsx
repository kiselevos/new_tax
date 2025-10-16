import { useState } from "react";
import { Link } from "react-router-dom";

export interface FormState {
  grossSalary: string;
  territorialMultiplier: string;
  northernCoefficient: string;
  startDate: string;
  hasTaxPrivilege: boolean;
  isNotResident: boolean;
}

interface TaxFormProps {
  onSubmit: (form: FormState) => void;
  loading: boolean;
}

interface TooltipProps {
  content: string;
  children: React.ReactNode;
}

const Tooltip = ({ content, children }: TooltipProps) => {
  return (
    <div className="tooltip-container">
      {children}
      <div className="tooltip">{content}</div>
    </div>
  );
};


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

export default function TaxForm({ onSubmit, loading }: TaxFormProps) {
  const [form, setForm] = useState<FormState>({
    grossSalary: "",
    territorialMultiplier: "",
    northernCoefficient: "",
    startDate: "",
    hasTaxPrivilege: false,
    isNotResident: false,
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
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

  const handleSalaryInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value.replace(",", ".").replace(/[^0-9.]/g, "");
    const normalized = raw.split(".").length > 2 ? raw.replace(/\.+$/, "") : raw;
    setForm((prev) => ({ ...prev, grossSalary: normalized }));
  };

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

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(form);
  };

  return (
    <form onSubmit={handleSubmit} className="form">
      {/* --- оклад --- */}
      <div className="field-row">
        <div className="field-left">
          <Tooltip content="Укажите ваш оклад, указанный в трудовом договоре. Это сумма до вычета налогов и без надбавок.">
            <label htmlFor="grossSalary">Оклад:</label>
          </Tooltip>
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
    <Tooltip content="Устанавливается в зависимости от региона, где вы работаете. Подробнее →">
    <label htmlFor="territorialMultiplier">Территориальный коэффициент:</label>
    </Tooltip>
    <Link to="/regional-info" className="field-link under-title">
      Подробнее →
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
      {Array.from({ length: 21 }, (_, i) => (1 + i * 0.05).toFixed(2))
        .filter((val) => parseFloat(val) >= 1.10) // ✅ убираем 1.00 и 1.05
        .map((val) => (
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
          <Tooltip content="Северная надбавка зависит от региона и стажа работы. Подробнее →">
            <label htmlFor="northernCoefficient">Северная надбавка:</label>
          </Tooltip>
          <Link to="/regional-info" className="field-link under-title">
            Подробнее →
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
            {Array.from({ length: 10 }, (_, i) => (i + 1) * 10).map((val) => (
              <option key={val} value={val}>
                {val}%
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* --- дата начала расчёта --- */}
<div className="field-row">
  <div className="field-left">
    <Tooltip content="На этом этапе расчёт выполняется с 1-го числа выбранного месяца текущего года.">
      <label>Дата начала расчёта:</label>
    </Tooltip>
  </div>
  <div className="field-right">
    <div className="date-field">
      <span>01</span>
      <select
        name="startDate"
        value={
          form.startDate
            ? form.startDate.split("-")[1] // если пользователь уже выбрал — показать выбранный месяц
            : "01" // ✅ по умолчанию всегда январь
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
          <Tooltip content="Укажите льготы или ограничения, если они есть. Они повлияют на итоговый расчёт. Подробнее →">
            <label>Параметры:</label>
          </Tooltip>
          <Link to="/special-tax-modes" className="field-link under-title">
            Подробнее →
          </Link>
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
      
      {/* Примечание о работе калькулятора */}
      <div className="calculator-note">
        <p>🔍 Все расчёты выполняются по открытым формулам и не сохраняются на сервере. Это инструмент для справочных целей, чтобы вы могли оценить свой доход «на руки» и понять, как именно формируется сумма налога.</p>
        <Link to="/how-it-works" className="calculator-link">
          подробнее о работе калькулятора →
        </Link>
      </div>
    </form>
  );
}