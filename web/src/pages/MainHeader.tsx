import { safeConvertToNumber } from "../utils";
import "./MainHeader.css";

interface MainHeaderProps {
  result: any;
  formData: any;
  currentYear: number;
}

export default function MainHeader({ result, formData, currentYear }: MainHeaderProps) {
  return (
    <div className="card main-header">
      <h2>Сводка расчета</h2>

      <div className="field-row">
        <div className="field-left">
          <span className="field-label">Период расчета:</span>
        </div>
        <div className="field-right">
          <span className="field-value">
            {formData?.startDate
              ? new Date(formData.startDate).toLocaleDateString("ru-RU", {
                  day: "numeric",
                  month: "long",
                })
              : "01 января"}{" "}
            – 31 декабря {currentYear}
          </span>
        </div>
      </div>

      {/* Оклад */}
      <div className="field-row">
        <div className="field-left">
          <span className="field-label">Оклад:</span>
        </div>
        <div className="field-right">
          <span className="field-value field-value-bold">
            {safeConvertToNumber(result.grossSalary ?? result.gross_salary).toLocaleString(
              "ru-RU",
              { minimumFractionDigits: 2 }
            )}{" "}
            ₽
          </span>
        </div>
      </div>

      {/* Параметры расчета */}
      <div className="parameters-section">
        {/* Территориальный коэффициент - показываем только если не 1.00 */}
        {result.territorialMultiplier && Number(result.territorialMultiplier) > 100 && (
          <div className="field-row">
            <div className="field-left">
              <span className="field-label">Территориальный коэффициент:</span>
            </div>
            <div className="field-right">
              <span className="field-value field-value-bold">
                {(Number(result.territorialMultiplier) / 100).toFixed(2)}
              </span>
            </div>
          </div>
        )}

        {/* Северная надбавка - показываем только если не 0% */}
        {result.northernCoefficient && Number(result.northernCoefficient) > 100 && (
          <div className="field-row">
            <div className="field-left">
              <span className="field-label">Северная надбавка:</span>
            </div>
            <div className="field-right">
              <span className="field-value field-value-bold">
                {Number(result.northernCoefficient) - 100}%
              </span>
            </div>
          </div>
        )}

        {/* Льготы для силовых структур */}
        {formData?.hasTaxPrivilege && (
          <div className="field-row">
            <div className="field-left">
              <span className="field-label">Льготы:</span>
            </div>
            <div className="field-right">
              <span className="field-value field-value-bold">
                Для силовых структур
              </span>
            </div>
          </div>
        )}

        {/* Налоговый нерезидент */}
        {formData?.isNotResident && (
          <div className="field-row">
            <div className="field-left">
              <span className="field-label">Статус:</span>
            </div>
            <div className="field-right">
              <span className="field-value field-value-bold">
                Нерезидент РФ
              </span>
            </div>
          </div>
        )}
      </div>

      {/* Финансовые результаты */}
      <div className="financial-results">
        <div className="field-row financial-row">
          <div className="field-left">
            <span className="financial-label">Общий доход:</span>
          </div>
          <div className="field-right">
            <span className="financial-value income-value">
              {safeConvertToNumber(result.annualGrossIncome).toLocaleString("ru-RU", { minimumFractionDigits: 2 })} ₽
            </span>
          </div>
        </div>
        
        <div className="field-row financial-row">
          <div className="field-left">
            <span className="financial-label">Налог к уплате:</span>
          </div>
          <div className="field-right">
            <span className="financial-value tax-value">
              {safeConvertToNumber(result.annualTaxAmount).toLocaleString("ru-RU", { minimumFractionDigits: 2 })} ₽
            </span>
          </div>
        </div>
        
        <div className="field-row financial-row">
          <div className="field-left">
            <span className="financial-label">Чистый доход:</span>
          </div>
          <div className="field-right">
            <span className="financial-value net-income-value">
              {safeConvertToNumber(result.annualNetIncome).toLocaleString("ru-RU", { minimumFractionDigits: 2 })} ₽
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}