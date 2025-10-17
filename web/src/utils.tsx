import type { ReactNode } from "react";

/**
 * Безопасное преобразование значений (BigInt, строка, число, объект) в число (в рублях)
 * Все значения приходят в копейках → делим на 100
 */
export const safeConvertToNumber = (value: unknown): number => {
  if (value == null || value === "") return 0;

  try {
    if (typeof value === "bigint") return Number(value) / 100;

    if (typeof value === "string") {
      const clean = value.replace(/[^\d-]/g, "");
      if (!clean) return 0;
      const num = Number(clean);
      return isNaN(num) ? 0 : num / 100;
    }

    if (typeof value === "number") return value / 100;

    if (typeof value === "object") {
      const obj = value as Record<string, unknown>;

      if ("low" in obj && "high" in obj) {
        const low = Number((obj as { low?: number }).low ?? 0);
        const high = Number((obj as { high?: number }).high ?? 0);
        return (low + high * 4294967296) / 100;
      }

      const raw = Object.values(obj)[0];
      if (typeof raw === "string" || typeof raw === "number" || typeof raw === "bigint") {
        return safeConvertToNumber(raw);
      }
    }

    const num = Number(value);
    return isNaN(num) ? 0 : num / 100;
  } catch (err) {
    console.warn("safeConvertToNumber error:", err, value);
    return 0;
  }
};

/** Безопасное получение даты */
export const safeGetDate = (
  month: { seconds?: bigint | number } | string | undefined
): Date | null => {
  if (!month) return null;

  if (typeof month === "string") return new Date(month);

  if (typeof month === "object" && "seconds" in month && month.seconds != null) {
    const seconds =
      typeof month.seconds === "bigint" ? Number(month.seconds) : (month.seconds as number);
    return new Date(seconds * 1000);
  }

  return null;
};

/** Форматирование числа в рубли */
export const formatRub = (value: number): string =>
  value.toLocaleString("ru-RU", { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + " ₽";

/** Универсальный JSX-рендер поля */
export const renderValue = (
  label: string,
  value: unknown,
  color?: string,
  showZero: boolean = false
): ReactNode => {
  const num = safeConvertToNumber(value);
  if (!showZero && num === 0) return null;

  return (
    <div className="field-row" key={label}>
      <div className="field-left">
        <span style={{ fontSize: "0.9rem", color: "#666" }}>{label}:</span>
      </div>
      <div className="field-right">
        <span
          style={{
            fontSize: "0.9rem",
            color: color || "#333",
            textAlign: "right",
            fontWeight: num !== 0 ? 500 : 400,
          }}
        >
          {formatRub(num)}
        </span>
      </div>
    </div>
  );
};
