// src/types/tax.types.ts

/** Формa запроса на расчёт налога */
export interface FormData {
  salary: number; // Оклад, указанный в запросе
  territorialMultiplier?: number; // Территориальный коэффициент
  northernCoefficient?: number; // Северная надбавка
  startDate?: string; // Дата начала расчёта (ISO string)
  hasTaxPrivilege?: boolean; // Льготы для силовых структур
  isNotResident?: boolean; // Статус налогового нерезидента
}

/** Детализация расчёта за месяц */
export interface MonthlyPrivateTax {
  month?: { seconds: bigint } | string; // google.protobuf.Timestamp или ISO строка

  monthlyGrossIncome?: bigint | string;
  monthlyNetIncome?: bigint | string;
  annualGrossIncome?: bigint | string;
  annualNetIncome?: bigint | string;

  taxRate?: bigint | string;
  monthlyTaxAmount?: bigint | string;
  annualTaxAmount?: bigint | string;

  monthlyNorthGrossIncome?: bigint | string;
  monthlyBaseGrossIncome?: bigint | string;

  monthlyNorthTaxAmount?: bigint | string;
  monthlyBaseTaxAmount?: bigint | string;

  annualNorthGrossIncome?: bigint | string;
  annualBaseGrossIncome?: bigint | string;

  annualNorthTaxAmount?: bigint | string;
  annualBaseTaxAmount?: bigint | string;
}

/** Общий ответ на расчёт (CalculatePrivateResponse) */
export interface ResultData {
  monthlyDetails: MonthlyPrivateTax[];

  annualTaxAmount?: bigint | string;
  annualGrossIncome?: bigint | string;
  annualNetIncome?: bigint | string;

  grossSalary?: bigint | string;
  territorialMultiplier?: bigint | string;
  northernCoefficient?: bigint | string;
}

/** Пропсы для MonthlyAccordion */
export interface MonthlyAccordionProps {
  monthlyDetails: MonthlyPrivateTax[];
  result: ResultData;
}

/** Пропсы для MainHeader */
export interface MainHeaderProps {
  result: ResultData;
  formData: FormData;
  currentYear: number;
}


export interface MonthlyAccordionProps {
  monthlyDetails: MonthlyPrivateTax[];
  result: ResultData;
}

export interface LocationState {
  result: ResultData;
  formData: FormData;
}