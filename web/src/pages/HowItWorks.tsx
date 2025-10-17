import { useNavigate } from "react-router-dom";
import "./InfoPages.css";

export default function HowItWorks() {
  const navigate = useNavigate();

  return (
    <div className="info-page">
      <h1 className="info-page-title">
        Как работает этот калькулятор
      </h1>

      <p className="info-page-subtitle">
        Этот калькулятор помогает рассчитать подоходный налог (НДФЛ) в России с учётом
        вашего оклада, региональных коэффициентов и северных надбавок.
        Всё просто и прозрачно — вы можете увидеть, как меняется доход по месяцам
        и сколько налога удерживается.
      </p>

      {/* Основная идея */}
      <div className="info-card">
        <h2 className="info-card-title">🔹 Основная идея</h2>
        <p className="info-card-text">
          Калькулятор рассчитывает ваш доход и налог помесячно, с начала выбранного периода согласно законодательству РФ. 
          Все расчёты ведутся по календарному году — с января по декабрь.
        </p>
        <p className="info-card-text">
          С 1 января 2025 года в России действует <strong>пятиступенчатая прогрессивная шкала НДФЛ</strong>.
          Чем выше годовой доход, тем выше ставка:
        </p>
        <ul className="info-bullet-list">
          <li>до 2,4 млн ₽ — 13%</li>
          <li>от 2,4 до 5 млн ₽ — 15%</li>
          <li>от 5 до 20 млн ₽ — 18%</li>
          <li>от 20 до 50 млн ₽ — 20%</li>
          <li>свыше 50 млн ₽ — 22%</li>
        </ul>
        <p className="info-card-text">
          То есть налог рассчитывается постепенно: к каждому диапазону применяется своя ставка.
          Это и есть прогрессивное налогообложение — чем больше зарабатываешь, тем выше доля налога
          именно на превышающую часть дохода.
        </p>
        <p className="info-source">
          Источник:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_480709/"
            target="_blank"
            rel="noopener noreferrer"
            className="info-source-link"
          >
            Федеральный закон № 389-ФЗ от 25.07.2024 г.
          </a>{" "}
          (вступает в силу с 1 января 2025 года);{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_28165/"
            target="_blank"
            rel="noopener noreferrer"
            className="info-source-link"
          >
            статья 224 Налогового кодекса РФ
          </a>
          .
        </p>
      </div>

      {/* Коэффициенты */}
      <div className="info-card">
        <h2 className="info-card-title">📍 Районный коэффициент и северная надбавка</h2>
        <p className="info-card-text">
          В некоторых регионах России оклады повышаются с помощью{" "}
          <strong>районных коэффициентов</strong> и{" "}
          <strong>северных надбавок</strong>. Эти надбавки компенсируют
          сложные климатические условия и высокие расходы на жизнь.
        </p>
        <p className="info-card-text">
          Калькулятор сначала увеличивает оклад на районный коэффициент —
          это ваш «региональный оклад». Затем отдельно рассчитывается северная
          надбавка — с неё налог считается по упрощённой ставке 13% до 5 млн руб.
        </p>
        <p className="info-card-text">
          Подробнее о региональных коэффициентах можно узнать{" "}
          <a
            href="/regional-info"
            className="info-source-link"
          >
            здесь
          </a>
          .
        </p>
        <p className="info-source">
          Источники:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_118861/"
            target="_blank"
            rel="noopener noreferrer"
            className="info-source-link"
          >
            Закон РФ «О районных коэффициентах и надбавках»
          </a>
          ,{" "}
          <a
            href="https://www.nalog.gov.ru/rn29/news/activities_fts/15907099/"
            target="_blank"
            rel="noopener noreferrer"
            className="info-source-link"
          >
            Разъяснение ФНС России о северных надбавках
          </a>
          .
        </p>
      </div>

      {/* Как считается налог */}
      <div className="info-card">
        <h2 className="info-card-title">💰 Как считается налог</h2>
        <p className="info-card-text">
          Каждый месяц к сумме вашего дохода добавляется налог, рассчитанный по накопительной
          системе. Это значит, что ставка может постепенно меняться в течение года,
          если ваш доход переходит в следующий диапазон.
        </p>
        <p className="info-card-text">
          Налог округляется до целого рубля по правилам Налогового кодекса:
          суммы меньше 50 копеек отбрасываются, а 50 копеек и больше округляются вверх.
        </p>
        <p className="info-source">
          Источник:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_28165/27e8ceefb17f7ef8f4c7e8db918fe86b3543ab2b/"
            target="_blank"
            rel="noopener noreferrer"
            className="info-source-link"
          >
            п. 6 ст. 52 Налогового кодекса РФ
          </a>
          .
        </p>
      </div>

      {/* Особые случаи */}
      <div className="info-card">
        <h2 className="info-card-title">⚙️ Особые случаи</h2>
        <p className="info-card-text">
          В калькуляторе предусмотрены отдельные режимы для:
        </p>
        <ul className="info-bullet-list">
          <li>💼 сотрудников силовых структур — с льготной шкалой 13–15%;</li>
          <li>🌍 налоговых нерезидентов — ставка 30% на все доходы.</li>
        </ul>
        <p className="info-card-text">
          Эти режимы можно выбрать вручную при расчёте, если они применимы к вам.
        </p>
        <p className="info-source">
          Источник:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_28165/ed7f95b14e4d8f9a0c99df5dd911f292482e62cb/"
            target="_blank"
            rel="noopener noreferrer"
            className="info-source-link"
          >
            статья 224 Налогового кодекса РФ
          </a>
          .
        </p>
      </div>

      <div className="info-button-container">
        <button
          onClick={() => navigate("/")}
          className="info-button"
        >
          Вернуться к калькулятору
        </button>
      </div>
    </div>
  );
}