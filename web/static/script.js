// static/script.js
document.addEventListener("DOMContentLoaded", function() {
    const MIN_SALARY = window.APP_CONFIG?.minSalary;
    const LIVING_WAGE = window.APP_CONFIG?.livingWage;

    const form = document.querySelector("form");
    const salaryInput = document.getElementById("grossSalary");
    
    if (!form || !salaryInput) return;

    const salaryField = salaryInput.closest(".field-row");

    // Создаем элемент для сообщений
    const messageEl = document.createElement("div");
    messageEl.classList.add("field-message");
    salaryField.appendChild(messageEl);

    // Функции для показа сообщений
    function showError(msg) {
        messageEl.textContent = msg;
        messageEl.className = "field-message error";
        salaryInput.classList.add("field-error");
        salaryInput.focus();
    }

    function showWarning(msg) {
        messageEl.textContent = msg;
        messageEl.className = "field-message warning";
        salaryInput.classList.add("field-warning");
    }

    function clearMessages() {
        messageEl.textContent = "";
        messageEl.className = "field-message";
        salaryInput.classList.remove("field-error", "field-warning");
    }

    // Валидация при вводе
    salaryInput.addEventListener('input', function() {
        const value = this.value.trim().replace(/\s/g, '').replace(',', '.');
        if (value && !isNaN(parseFloat(value))) {
            clearMessages();
        }
    });

    // Валидация при отправке формы
    form.addEventListener("submit", function(e) {
        e.preventDefault();
        
        const value = salaryInput.value.trim().replace(/\s/g, '').replace(',', '.');
        const salary = parseFloat(value);
        
        clearMessages();

        let isValid = true;
        let warningMessage = "";

        // Проверка на пустое или некорректное значение
        if (!value || isNaN(salary) || salary <= 0) {
            showError("Введите корректный оклад (например, 20000.00)");
            isValid = false;
        }
        // Проверка минимального оклада
        else if (salary < MIN_SALARY) {
            showError(`Минимальный оклад — ${MIN_SALARY.toLocaleString('ru-RU')} ₽`);
            isValid = false;
        }
        // Проверка на прожиточный минимум (предупреждение)
        else if (salary < LIVING_WAGE) {
            warningMessage = `Сумма меньше прожиточного минимума (${LIVING_WAGE.toLocaleString('ru-RU')} ₽). Проверьте данные.`;
        }

        // Если есть предупреждение, показываем его но все равно отправляем
        if (warningMessage && isValid) {
            showWarning(warningMessage);
            setTimeout(() => {
                form.submit();
            }, 1000);
        } 
        // Если валидация прошла успешно, отправляем форму
        else if (isValid) {
            form.submit();
        }
    });

    // Автоматическое форматирование оклада
    salaryInput.addEventListener('blur', function() {
        let value = this.value.trim().replace(/\s/g, '').replace(',', '.');
        if (value && !isNaN(parseFloat(value))) {
            const number = parseFloat(value);
            this.value = number.toLocaleString('ru-RU', {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            });
        }
    });
    
    salaryInput.addEventListener('focus', function() {
        let value = this.value.replace(/\s/g, '').replace(',', '.');
        if (value && !isNaN(parseFloat(value))) {
            this.value = parseFloat(value);
        }
    });

    // Обработка тултипов
    initTooltips();
});