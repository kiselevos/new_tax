import { ExternalLink } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function RegionalInfo() {
  const navigate = useNavigate();

  return (
    <div
      style={{
        maxWidth: "800px",
        margin: "40px auto",
        fontFamily: "sans-serif",
        color: "#333",
      }}
    >
      <h1 style={{ fontSize: "28px", fontWeight: 700, textAlign: "center" }}>
        Региональные коэффициенты и северные надбавки
      </h1>

      <p style={{ fontSize: "17px", color: "#555", textAlign: "center", marginBottom: "24px" }}>
        Районный коэффициент и северная надбавка зависят от региона работы.
        Эти показатели устанавливаются федеральными и региональными актами.
        Ниже приведены официальные источники, где можно уточнить действующие значения.
      </p>

      {/* Районный коэффициент */}
      <div style={cardStyle}>
        <h2 style={cardTitle}>📍 Районный коэффициент (РК)</h2>
        <ul style={linkList}>
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_256207/660c6bc7428f5eba28690e62cb42131f66d3b514/"
              target="_blank"
              rel="noopener noreferrer"
              style={linkStyle}
            >
              <ExternalLink size={16} style={{ marginRight: 6 }} />
              Федеральная таблица районных коэффициентов (КонсультантПлюс)
            </a>
          </li>
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_118861/"
              target="_blank"
              rel="noopener noreferrer"
              style={linkStyle}
            >
              <ExternalLink size={16} style={{ marginRight: 6 }} />
              Закон РФ «О районных коэффициентах и надбавках»
            </a>
          </li>
        </ul>
      </div>

      {/* Северная надбавка */}
      <div style={cardStyle}>
        <h2 style={cardTitle}>❄️ Северная надбавка (СН)</h2>
        <ul style={linkList}>
          <li>
            <a
              href="https://www.nalog.gov.ru/rn29/news/activities_fts/15907099/"
              target="_blank"
              rel="noopener noreferrer"
              style={linkStyle}
            >
              <ExternalLink size={16} style={{ marginRight: 6 }} />
              Разъяснение ФНС о порядке налогообложения северных надбавок
            </a>
          </li>
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_118861/"
              target="_blank"
              rel="noopener noreferrer"
              style={linkStyle}
            >
              <ExternalLink size={16} style={{ marginRight: 6 }} />
              Обзор законодательства о северных надбавках (КонсультантПлюс)
            </a>
          </li>
          <li>
            <a
              href="https://www.buhsoft.ru/article/4542-severnaya-nadbavka-v-2025-godu-razmer-tablitsa"
              target="_blank"
              rel="noopener noreferrer"
              style={linkStyle}
            >
              <ExternalLink size={16} style={{ marginRight: 6 }} />
              Таблица северных надбавок по регионам (БухСофт)
            </a>
          </li>
        </ul>
      </div>

      {/* Примечание */}
      <div
        style={{
          background: "#fafafa",
          border: "1px dashed #ccc",
          padding: "12px 16px",
          borderRadius: "8px",
          marginBottom: "20px",
          fontSize: "14px",
          color: "#555",
        }}
      >
        ⚠️ Обратите внимание: даже внутри одного региона коэффициенты и надбавки
        могут различаться по районам. Для точных данных ориентируйтесь на
        постановления местных органов власти.
      </div>

      <div style={{ textAlign: "center" }}>
        <button
          onClick={() => navigate("/")}
          style={{
            background: "#2b70f7",
            color: "#fff",
            border: "none",
            borderRadius: "6px",
            padding: "10px 20px",
            fontSize: "16px",
            cursor: "pointer",
          }}
        >
          🔙 Вернуться к калькулятору
        </button>
      </div>
    </div>
  );
}

const cardStyle = {
  background: "#fff",
  border: "1px solid #eee",
  borderRadius: "10px",
  padding: "16px 20px",
  marginBottom: "20px",
  boxShadow: "0 2px 6px rgba(0, 0, 0, 0.05)",
};

const cardTitle = {
  fontSize: "20px",
  marginBottom: "8px",
  fontWeight: 600,
};

const linkList = {
  listStyle: "none",
  padding: 0,
  margin: 0,
};

const linkStyle = {
  display: "flex",
  alignItems: "center",
  color: "#2b70f7",
  textDecoration: "none",
  marginBottom: "6px",
};
