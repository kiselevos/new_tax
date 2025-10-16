import { ExternalLink } from "lucide-react";
import { useNavigate } from "react-router-dom";
import "./InfoPages.css";

export default function RegionalInfo() {
  const navigate = useNavigate();

  return (
    <div className="info-page">
      <h1 className="info-page-title">
        Районные коэффициенты и северные надбавки
      </h1>

      <p className="info-page-subtitle">
        В некоторых регионах России работники получают повышенную зарплату — 
        это компенсация за работу в особых климатических условиях.  
        Размер таких надбавок и коэффициентов установлен федеральными и региональными законами.
      </p>

      {/* Районный коэффициент */}
      <div className="info-card">
        <h2 className="info-card-title">📍 Районный коэффициент</h2>
        <p className="info-card-text">
          Районный коэффициент — это повышающий коэффициент к зарплате, 
          установленный для работников в районах с суровыми климатическими условиями.  
          Он применяется к окладу с первого дня работы и не зависит от стажа.
        </p>
        <ul className="info-link-list">
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_256207/660c6bc7428f5eba28690e62cb42131f66d3b514/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-link"
            >
              <ExternalLink size={16} className="info-link-icon" />
              Федеральная таблица районных коэффициентов (КонсультантПлюс)
            </a>
          </li>
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_118861/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-link"
            >
              <ExternalLink size={16} className="info-link-icon" />
              Закон РФ «О районных коэффициентах и надбавках»
            </a>
          </li>
        </ul>
      </div>

      {/* Северная надбавка */}
      <div className="info-card">
        <h2 className="info-card-title">❄️ Северная надбавка</h2>
        <p className="info-card-text">
          Северная надбавка — это дополнительный процент к заработной плате за работу 
          в районах Крайнего Севера и приравненных местностях.  
          Её размер зависит от стажа работы в этих условиях и может достигать 100%.
        </p>
        <ul className="info-link-list">
          <li>
            <a
              href="https://www.nalog.gov.ru/rn29/news/activities_fts/15907099/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-link"
            >
              <ExternalLink size={16} className="info-link-icon" />
              Разъяснение ФНС России о порядке налогообложения северных надбавок
            </a>
          </li>
          <li>
            <a
              href="https://www.consultant.ru/document/cons_doc_LAW_118861/"
              target="_blank"
              rel="noopener noreferrer"
              className="info-link"
            >
              <ExternalLink size={16} className="info-link-icon" />
              Обзор законодательства о северных надбавках (КонсультантПлюс)
            </a>
          </li>
          <li>
            <a
              href="https://www.buhsoft.ru/article/4542-severnaya-nadbavka-v-2025-godu-razmer-tablitsa"
              target="_blank"
              rel="noopener noreferrer"
              className="info-link"
            >
              <ExternalLink size={16} className="info-link-icon" />
              Таблица размеров северных надбавок по регионам (БухСофт)
            </a>
          </li>
        </ul>
      </div>

      {/* Примечание */}
      <div className="info-note">
        ⚠️ Обратите внимание: даже внутри одного субъекта Российской Федерации 
        коэффициенты и надбавки могут отличаться по районам и населенным пунктам.  
        Для точных данных рекомендуется проверять постановления местных органов власти 
        или кадровую документацию вашего работодателя.
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