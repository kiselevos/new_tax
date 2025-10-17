import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./App.css";
import { taxClient } from "./client";
import TaxForm, { type FormState } from "./pages/TaxForm";

export default function App() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(form: FormState) {
    setLoading(true);
    setError(null);

    

    const salary = parseFloat(form.grossSalary.replace(/\s/g, "").replace(",", "."));
    if (isNaN(salary) || salary <= 0) {
      setError("Некорректное значение оклада");
      setLoading(false);
      return;
    }

    try {
      const result = await taxClient.calculatePrivate({
        grossSalary: BigInt(Math.round(salary * 100)),
        territorialMultiplier: form.territorialMultiplier
          ? BigInt(Math.round(parseFloat(form.territorialMultiplier) * 100))
          : BigInt(0),
        northernCoefficient: form.northernCoefficient
          ? BigInt(Number(form.northernCoefficient) + 100)
          : BigInt(0),
        startDate: form.startDate
          ? {
              seconds: BigInt(Math.floor(new Date(form.startDate).getTime() / 1000)),
            }
          : undefined,
        hasTaxPrivilege: form.hasTaxPrivilege,
        isNotResident: form.isNotResident,
      });

      console.log('Navigation to results with:', { result, formData: form });
navigate("/results", { 
  state: { 
    result: JSON.parse(
      JSON.stringify(result, (_, v) => 
        typeof v === "bigint" ? v.toString() : v
      )
    ),
    formData: form
  } 
});    
} catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка при запросе");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="container">
      <h1>Расчёт налога на текущий год</h1>

      <TaxForm onSubmit={handleSubmit} loading={loading} />
      
      {error && <p className="error">{error}</p>}
    </div>
  );
}
