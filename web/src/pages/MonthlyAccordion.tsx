import { safeConvertToNumber, safeGetDate } from "../utils";
import type { MonthlyAccordionProps } from "./types/tax.types";
import "./MonthlyAccordion.css";

export default function MonthlyAccordion({ monthlyDetails }: MonthlyAccordionProps) {
  if (!monthlyDetails || monthlyDetails.length === 0) {
    return (
      <div className="card">
        <h3>Нет данных для помесячной детализации</h3>
      </div>
    );
  }

  return (
    <div className="monthly-details card">
      <h2>Помесячная детализация</h2>

      <div className="monthly-accordion">
        {monthlyDetails.map((month, index) => {
          const monthDate = safeGetDate(month.month);
          const monthName = monthDate
            ? monthDate.toLocaleDateString("ru-RU", { month: "long", year: "numeric" })
            : "Не указан";

          const taxRate = safeConvertToNumber(month.taxRate);
          const monthlyGross = safeConvertToNumber(month.monthlyGrossIncome);
          const monthlyNet = safeConvertToNumber(month.monthlyNetIncome);
          const monthlyTax = safeConvertToNumber(month.monthlyTaxAmount);
          
          const monthlyBaseGross = safeConvertToNumber(month.monthlyBaseGrossIncome);
          const monthlyNorthGross = safeConvertToNumber(month.monthlyNorthGrossIncome);
          const monthlyBaseTax = safeConvertToNumber(month.monthlyBaseTaxAmount);
          const monthlyNorthTax = safeConvertToNumber(month.monthlyNorthTaxAmount);
          
          const annualGross = safeConvertToNumber(month.annualGrossIncome);
          const annualNet = safeConvertToNumber(month.annualNetIncome);
          const annualTax = safeConvertToNumber(month.annualTaxAmount);

          return (
            <div key={index} className="month-accordion-item">
              {/* Заголовок месяца */}
              <div className="month-header">
                <div className="month-summary">
                  <span className="month-name">{monthName.charAt(0).toUpperCase() + monthName.slice(1)}</span>
                  <span className="tax-info">
  Ставка НДФЛ:{" "}
  <span className="tax-rate-value">
    {taxRate !== null && taxRate !== undefined ? `${Math.round(taxRate * 100)}%` : "—"}
  </span>
</span>
                </div>
                <button
                  className="expand-button"
                  onClick={(e) =>
                    e.currentTarget.parentElement?.parentElement?.classList.toggle("expanded")
                  }
                >
                  ▼
                </button>
              </div>

              {/* Детали месяца */}
              <div className="month-details">
                <div className="details-content">
                  
                  {/* Доходы за месяц */}
                  <div className="section first-section">
                    <div className="section-header">
                      <span className="section-title">Доходы за месяц</span>
                    </div>
                    
                    <div className="section-items">
                      {monthlyBaseGross > 0 && (
                        <IncomeRow 
                          label="Базовый доход (оклад + тер. коэффициент)" 
                          value={monthlyBaseGross} 
                          color="normal"
                        />
                      )}
                      {monthlyNorthGross > 0 && (
                        <IncomeRow 
                          label="Северная надбавка" 
                          value={monthlyNorthGross} 
                          color="normal"
                        />
                      )}
                      <IncomeRow 
                        label="Доход до вычета НДФЛ" 
                        value={monthlyGross} 
                        color="blue"
                      />
                      <IncomeRow 
                        label="Чистый доход за месяц" 
                        value={monthlyNet} 
                        color="green"
                        important={true}
                      />
                    </div>
                  </div>

                  <div className="section-divider"></div>

                  {/* Налоги за месяц */}
                  <div className="section">
                    <div className="section-header">
                      <span className="section-title">Налоги за месяц</span>
                    </div>
                    
                    <div className="section-items">
                      {monthlyBaseTax > 0 && (
                        <IncomeRow 
                          label="НДФЛ с основного дохода" 
                          value={monthlyBaseTax} 
                          color="normal"
                        />
                      )}
                      {monthlyNorthTax > 0 && (
                        <IncomeRow 
                          label="НДФЛ с северной надбавки" 
                          value={monthlyNorthTax} 
                          color="normal"
                        />
                      )}
                      <IncomeRow 
                        label="Удержанный НДФЛ" 
                        value={monthlyTax} 
                        color="red"
                        important={true}
                      />
                    </div>
                  </div>

                  <div className="section-divider"></div>

                  {/* Накопительно с начала периода */}
                  <div className="section">
                    <div className="section-header">
                      <span className="section-title">Накопительно с начала периода</span>
                    </div>
                    
                    <div className="section-items">
                      <IncomeRow 
                        label="Доход до вычета НДФЛ" 
                        value={annualGross} 
                        color="blue"
                      />
                      <IncomeRow 
                        label="Удержанный НДФЛ" 
                        value={annualTax} 
                        color="red"
                      />
                      <IncomeRow 
                        label="Чистый доход с начанала периода" 
                        value={annualNet} 
                        color="green"
                        important={true}
                      />
                    </div>
                  </div>

                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// Компонент для отображения строки дохода/налога
interface IncomeRowProps {
  label: string;
  value: number;
  color: "normal" | "blue" | "green" | "red";
  important?: boolean;
}

function IncomeRow({ label, value, color, important = false }: IncomeRowProps) {
  const colorMap = {
    normal: "#666",
    blue: "#007bff",
    green: "#28a745", 
    red: "#dc3545"
  };

  const backgroundColorMap = {
    normal: "transparent",
    blue: "#f8f9ff",
    green: "#f8fff9",
    red: "#fff5f5"
  };

  const currentColor = colorMap[color];
  const currentBackground = important ? backgroundColorMap[color] : "transparent";

  return (
    <div className={`income-row ${important ? 'income-row-important' : ''}`} 
         style={{ 
           backgroundColor: currentBackground,
           borderColor: important ? `${currentColor}20` : 'transparent'
         }}>
      <span className={`income-label ${important ? 'income-label-important' : ''}`}>
        {important ? label : `• ${label}`}
      </span>
      <span className="income-value" style={{ color: currentColor }}>
        {value.toLocaleString("ru-RU", {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2
        })} ₽
      </span>
    </div>
  );
}