import { useNavigate } from "react-router-dom";
import "./InfoPages.css";

export default function SpecialTaxModes() {
  const navigate = useNavigate();

  return (
    <div className="info-page">
      <h1 className="info-page-title">
        Особые режимы налогообложения НДФЛ
      </h1>

      <p className="info-page-subtitle">
        В отдельных случаях налог на доходы физических лиц (НДФЛ) рассчитывается по
        специальным правилам. К ним относятся: <strong>льготы для сотрудников силовых структур</strong>
        и <strong>налогообложение нерезидентов Российской Федерации</strong>.
      </p>

      {/* Льготы для силовых структур */}
      <div className="info-card">
        <h2 className="info-card-title">💼 Льготы для сотрудников силовых структур</h2>

        <p className="info-card-text">
          Для военнослужащих, сотрудников МВД, ФСБ, Росгвардии, прокуратуры, Следственного комитета
          и других органов, установлены <strong>особые правила расчёта НДФЛ</strong>.
        </p>

        <ul className="info-bullet-list">
          <li>Весь доход облагается по <strong>упрощённой шкале</strong> — 13% до 5 млн ₽ и 15% свыше.</li>
          <li>Не разделяется на базу «оклад» и «надбавки» — налог считается с общей суммы дохода.</li>
          <li>Льгота применяется автоматически, если гражданин получает доход в рамках служебного контракта.</li>
        </ul>

        <p className="info-card-text">
          Эти категории освобождаются от применения полной прогрессивной шкалы (13–22%),
          принятой для гражданских лиц, и пользуются более мягким режимом.
        </p>

        <div className="info-example">
          ℹ️ Пример: если военный получает 120 000 ₽ в месяц, то налог удерживается по ставке 13%.
          Если его годовой доход превысит 5 млн ₽, ставка на сумму превышения составит 15%.
        </div>

        <p className="info-source">
          📚 Источники:
        </p>
        <ul className="info-source-list">
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_28165/ed7f95b14e4d8f9a0c99df5dd911f292482e62cb/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-source-link"
            >
              Статья 224 Налогового кодекса РФ — налоговые ставки
            </a>
          </li>
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_41812/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-source-link"
            >
              Федеральный закон «О статусе военнослужащих»
            </a>
          </li>
          <li>
            <a
              href="https://www.nalog.gov.ru/rn77/news/activities_fts/15018377/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-source-link"
            >
              Разъяснение ФНС России о налогообложении доходов силовых структур
            </a>
          </li>
        </ul>
      </div>

      {/* Налог для нерезидентов */}
      <div className="info-card">
        <h2 className="info-card-title">🌍 Налогообложение нерезидентов РФ</h2>

        <p className="info-card-text">
          Нерезидентом РФ считается физическое лицо, которое находится на территории России
          менее <strong>183 дней в течение последних 12 месяцев</strong>.
        </p>

        <p className="info-card-text">
          Для таких граждан применяется <strong>единая ставка 30%</strong> на все виды доходов,
          полученные из источников в России.
        </p>

        <div className="info-example">
          ℹ️ Пример: если нерезидент получает 100 000 ₽ в месяц, налог составит 30 000 ₽.
          Прогрессивная шкала (13–22%) к нему не применяется.
        </div>

        <p className="info-card-text">
          Исключение составляют некоторые категории иностранных специалистов, работающих
          по трудовому договору в России — для них ставка может быть снижена до 13%.
        </p>

        <p className="info-source">
          📚 Источники:
        </p>
        <ul className="info-source-list">
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_28165/37e18e3b7a1fa90f31c19b17f2ff74b9a8e883b5/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-source-link"
            >
              Статья 207 Налогового кодекса РФ — налоговое резидентство
            </a>
          </li>
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_28165/ed7f95b14e4d8f9a0c99df5dd911f292482e62cb/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-source-link"
            >
              Статья 224 Налогового кодекса РФ — налоговые ставки для нерезидентов
            </a>
          </li>
          <li>
            <a
              href="https://www.nalog.gov.ru/rn77/news/activities_fts/15338315/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-source-link"
            >
              Разъяснение ФНС России о налогообложении нерезидентов
            </a>
          </li>
        </ul>
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