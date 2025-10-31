document.addEventListener("DOMContentLoaded", function() {
    const form = document.querySelector("form");
    const salaryInput = document.getElementById("grossSalary");
    
    if (!form || !salaryInput) return;

    const MIN_SALARY = window.APP_CONFIG.minSalary;
    const MAX_SALARY = 100000000; // 💰 Максимальный оклад — 100 млн
    const salaryContainer = salaryInput.closest(".salary-field-container");
    
    // Создаем элемент для сообщений
    let messageEl = salaryContainer.querySelector('.field-message');
    if (!messageEl) {
        messageEl = document.createElement("div");
        messageEl.classList.add("field-message");
        salaryContainer.appendChild(messageEl);
    }

    function showError(msg) {
        messageEl.textContent = msg;
        messageEl.className = "field-message error";
        salaryInput.classList.add("field-error");
        salaryInput.focus();
        
        // Анимация для привлечения внимания
        salaryInput.style.animation = 'shake 0.3s ease';
        setTimeout(() => {
            salaryInput.style.animation = '';
        }, 300);
    }

    function clearMessages() {
        messageEl.textContent = "";
        messageEl.className = "field-message";
        salaryInput.classList.remove("field-error");
    }

    // 💡 Фильтрация: разрешаем цифры, запятую и точку
    salaryInput.addEventListener('input', function () {
  let value = this.value;

  // Убираем всё кроме цифр, точки и запятой
  value = value.replace(/[^0-9.,]/g, '');

  // Если несколько точек/запятых — оставляем только первую
  const firstSeparator = value.match(/[.,]/);
  if (firstSeparator) {
    const sep = firstSeparator[0];
    const pos = firstSeparator.index;
    // удаляем все остальные разделители
    value = value.replace(/[.,]/g, '');
    // вставляем только первую найденную
    value = value.slice(0, pos) + sep + value.slice(pos);
  }

  // Ограничиваем два знака после разделителя
  value = value.replace(/^(\d+)([.,])(\d{0,2}).*$/, (_, int, sep, frac) =>
    frac !== undefined ? int + sep + frac : int
  );

  this.value = value;
});

    // Валидация при отправке формы
    form.addEventListener("submit", function(e) {
        e.preventDefault();
        
        const rawValue = salaryInput.value.trim().replace(/\s/g, '').replace(',', '.');
        const salary = parseFloat(rawValue);
        
        clearMessages();
        let isValid = true;

        // Проверка на пустое или некорректное значение
        if (!rawValue || isNaN(salary)) {
            showError("Введите корректный оклад (например, 10000.00)");
            isValid = false;
        }
        // Проверка минимального оклада
        else if (salary < MIN_SALARY) {
            showError(`Минимальный оклад — ${MIN_SALARY.toLocaleString('ru-RU')} ₽`);
            isValid = false;
        }
        // Проверка максимального оклада
        else if (salary > MAX_SALARY) {
            showError(`Сумма слишком велика. Максимум — ${MAX_SALARY.toLocaleString('ru-RU')} ₽`);
            isValid = false;
        }

        // Если валидация прошла успешно — отправляем форму
        if (isValid) {
            form.submit();
        }
    });

    // Автоматическое форматирование при blur
    salaryInput.addEventListener('blur', function() {
        let value = this.value.trim().replace(/\s/g, '').replace(',', '.');
        if (value && !isNaN(parseFloat(value))) {
            const number = parseFloat(value);
            this.value = number.toLocaleString('ru-RU', {
                minimumFractionDigits: value.includes('.') || value.includes(',') ? 2 : 0,
                maximumFractionDigits: 2
            });
        }
    });
    
    salaryInput.addEventListener('focus', function() {
        this.value = this.value.replace(/\s/g, '').replace(',', '.');
    });

    // ======== Exclusive checkboxes ========
    function initExclusiveCheckboxes() {
        const options = document.querySelectorAll('.exclusive-option');
        console.log('🔍 Found exclusive options:', options.length);
        
        options.forEach(option => {
            const checkbox = option.querySelector('input[type="checkbox"]');
            if (!checkbox) return;
            
            checkbox.addEventListener('change', function() {
                if (this.checked) {
                    option.classList.add('selected');
                    
                    options.forEach(otherOption => {
                        if (otherOption !== option) {
                            const otherCheckbox = otherOption.querySelector('input[type="checkbox"]');
                            if (otherCheckbox) {
                                otherCheckbox.checked = false;
                                otherOption.classList.remove('selected');
                            }
                        }
                    });
                    showCheckboxHint(this);
                } else {
                    option.classList.remove('selected');
                    removeCheckboxHint();
                }
            });
        });
    }

    function showCheckboxHint(checkedCheckbox) {
        removeCheckboxHint();
        const hint = document.createElement('div');
        hint.className = 'checkbox-hint';
        
        if (checkedCheckbox.name === 'hasTaxPrivilege') {
            hint.textContent = '✓ Выбраны льготы для силовых структур';
        } else if (checkedCheckbox.name === 'isNotResident') {
            hint.textContent = '✓ Выбран статус налогового нерезидента';
        }
        
        const container = checkedCheckbox.closest('.exclusive-checkboxes');
        if (container) container.appendChild(hint);
    }

    function removeCheckboxHint() {
        const hints = document.querySelectorAll('.checkbox-hint');
        hints.forEach(hint => hint.remove());
    }

    initExclusiveCheckboxes();
    initTooltips();
});
