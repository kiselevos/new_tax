import { useNavigate } from "react-router-dom";

export default function HowItWorks() {
  const navigate = useNavigate();

  return (
    <div className="max-w-3xl mx-auto py-10 px-4 text-gray-800">
      <h1 className="text-3xl font-bold mb-4 text-center">
        Как работает этот калькулятор
      </h1>

      <p className="text-lg text-gray-600 mb-8 text-center">
        Этот калькулятор помогает рассчитать подоходный налог (НДФЛ) в России с учётом
        вашего оклада, региональных коэффициентов и северных надбавок.
        Всё просто и прозрачно — вы можете увидеть, как меняется доход по месяцам
        и сколько налога удерживается.
      </p>

      {/* Основная идея */}
      <div className="border border-gray-200 rounded-xl p-5 mb-6 shadow-sm bg-white">
        <h2 className="text-xl font-semibold mb-3">🔹 Основная идея</h2>
        <p className="mb-2 leading-relaxed">
          Калькулятор рассчитывает ваш годовой доход и налог помесячно —
          так же, как это делает бухгалтерия работодателя. Все расчёты ведутся
          по календарному году — с января по декабрь.
        </p>
        <p className="mb-3 leading-relaxed">
          С 1 января 2025 года в России действует <strong>пятиступенчатая прогрессивная шкала НДФЛ</strong>.
          Чем выше годовой доход, тем выше ставка:
        </p>
        <ul className="list-disc pl-6 mb-3 leading-relaxed">
          <li>до 2,4 млн ₽ — 13%</li>
          <li>от 2,4 до 5 млн ₽ — 15%</li>
          <li>от 5 до 20 млн ₽ — 18%</li>
          <li>от 20 до 50 млн ₽ — 20%</li>
          <li>свыше 50 млн ₽ — 22%</li>
        </ul>
        <p className="leading-relaxed mb-3">
          То есть налог рассчитывается постепенно: к каждому диапазону применяется своя ставка.
          Это и есть прогрессивное налогообложение — чем больше зарабатываешь, тем выше доля налога
          именно на превышающую часть дохода.
        </p>
        <p className="text-sm text-gray-500">
          Источник:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_480709/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline"
          >
            Федеральный закон № 389-ФЗ от 25.07.2024 г.
          </a>{" "}
          (вступает в силу с 1 января 2025 года);{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_28165/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline"
          >
            статья 224 Налогового кодекса РФ
          </a>
          .
        </p>
      </div>

      {/* Коэффициенты */}
      <div className="border border-gray-200 rounded-xl p-5 mb-6 shadow-sm bg-white">
        <h2 className="text-xl font-semibold mb-3">
          📍 Районный коэффициент и северная надбавка
        </h2>
        <p className="mb-2 leading-relaxed">
          В некоторых регионах России оклады повышаются с помощью{" "}
          <strong>районных коэффициентов</strong> и{" "}
          <strong>северных надбавок</strong>. Эти надбавки компенсируют
          сложные климатические условия и высокие расходы на жизнь.
        </p>
        <p className="mb-2 leading-relaxed">
          Калькулятор сначала увеличивает оклад на районный коэффициент —
          это ваш «региональный оклад». Затем отдельно рассчитывается северная
          надбавка — с неё налог считается по упрощённой ставке 13% до 5 млн руб.
        </p>
        <p className="leading-relaxed mb-3">
          Подробнее о региональных коэффициентах можно узнать{" "}
          <a
            href="/regional-info"
            className="text-blue-600 hover:underline"
          >
            здесь
          </a>
          .
        </p>
        <p className="text-sm text-gray-500">
          Источники:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_118861/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline"
          >
            Закон РФ «О районных коэффициентах и надбавках»
          </a>
          ,{" "}
          <a
            href="https://www.nalog.gov.ru/rn29/news/activities_fts/15907099/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline"
          >
            Разъяснение ФНС России о северных надбавках
          </a>
          .
        </p>
      </div>

      {/* Как считается налог */}
      <div className="border border-gray-200 rounded-xl p-5 mb-6 shadow-sm bg-white">
        <h2 className="text-xl font-semibold mb-3">💰 Как считается налог</h2>
        <p className="mb-2 leading-relaxed">
          Каждый месяц к сумме вашего дохода добавляется налог, рассчитанный по накопительной
          системе. Это значит, что ставка может постепенно меняться в течение года,
          если ваш доход переходит в следующий диапазон.
        </p>
        <p className="mb-3 leading-relaxed">
          Налог округляется до целого рубля по правилам Налогового кодекса:
          суммы меньше 50 копеек отбрасываются, а 50 копеек и больше округляются вверх.
        </p>
        <p className="text-sm text-gray-500">
          Источник:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_28165/27e8ceefb17f7ef8f4c7e8db918fe86b3543ab2b/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline"
          >
            п. 6 ст. 52 Налогового кодекса РФ
          </a>
          .
        </p>
      </div>

      {/* Особые случаи */}
      <div className="border border-gray-200 rounded-xl p-5 mb-6 shadow-sm bg-white">
        <h2 className="text-xl font-semibold mb-3">⚙️ Особые случаи</h2>
        <p className="mb-2 leading-relaxed">
          В калькуляторе предусмотрены отдельные режимы для:
        </p>
        <ul className="list-disc pl-6 mb-2">
          <li>💼 сотрудников силовых структур — с льготной шкалой 13–15%;</li>
          <li>🌍 налоговых нерезидентов — ставка 30% на все доходы.</li>
        </ul>
        <p className="leading-relaxed mb-3">
          Эти режимы можно выбрать вручную при расчёте, если они применимы к вам.
        </p>
        <p className="text-sm text-gray-500">
          Источник:{" "}
          <a
            href="https://www.consultant.ru/document/cons_doc_LAW_28165/ed7f95b14e4d8f9a0c99df5dd911f292482e62cb/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline"
          >
            статья 224 Налогового кодекса РФ
          </a>
          .
        </p>
      </div>

      {/* Прозрачность */}
      <div className="border border-dashed border-gray-300 rounded-xl p-5 mb-8 bg-gray-50">
        <p className="text-gray-700 text-sm leading-relaxed">
          🔍 Все расчёты выполняются по открытым формулам и не сохраняются на сервере.
          Это инструмент для справочных целей, чтобы вы могли оценить свой доход
          «на руки» и понять, как именно формируется сумма налога.
        </p>
      </div>

      {/* Кнопка */}
      <div className="flex justify-center">
        <button
          onClick={() => navigate("/")}
          className="bg-blue-600 hover:bg-blue-700 text-white text-lg font-medium px-6 py-2 rounded-lg transition-colors"
        >
          🔙 Вернуться к калькулятору
        </button>
      </div>
    </div>
  );
}
