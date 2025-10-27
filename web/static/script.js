document.addEventListener("DOMContentLoaded", function() {
    const form = document.querySelector("form");
    const salaryInput = document.getElementById("grossSalary");
    
    if (!form || !salaryInput) return;

    const MIN_SALARY = window.APP_CONFIG.minSalary;
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

    // Валидация при вводе
    salaryInput.addEventListener('input', function() {
        const value = this.value.trim().replace(/\s/g, '').replace(',', '.');
        if (value && !isNaN(parseFloat(value)) && parseFloat(value) >= MIN_SALARY) {
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

        // Проверка на пустое или некорректное значение
        if (!value || isNaN(salary) || salary <= 0) {
            showError("Введите корректный оклад (например, 10000.00)");
            isValid = false;
        }
        // Проверка минимального оклада
        else if (salary < MIN_SALARY) {
            showError(`Минимальный оклад — ${MIN_SALARY.toLocaleString('ru-RU')} ₽`);
            isValid = false;
        }

        // Если валидация прошла успешно, отправляем форму
        if (isValid) {
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

    function initExclusiveCheckboxes() {
        const options = document.querySelectorAll('.exclusive-option');
        console.log('🔍 Found exclusive options:', options.length); // Отладочный вывод
        
        options.forEach(option => {
            const checkbox = option.querySelector('input[type="checkbox"]');
            console.log('🔍 Checkbox:', checkbox); // Отладочный вывод
            
            checkbox.addEventListener('change', function() {
                console.log('🔍 Checkbox changed:', this.name, this.checked); // Отладочный вывод
                
                if (this.checked) {
                    // Добавляем класс selected
                    option.classList.add('selected');
                    
                    // Снимаем выбор с других
                    options.forEach(otherOption => {
                        if (otherOption !== option) {
                            const otherCheckbox = otherOption.querySelector('input[type="checkbox"]');
                            otherCheckbox.checked = false;
                            otherOption.classList.remove('selected');
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
        if (container) {
            container.appendChild(hint);
        }
    }

    function removeCheckboxHint() {
        const hints = document.querySelectorAll('.checkbox-hint');
        hints.forEach(hint => hint.remove());
    }

    initExclusiveCheckboxes();
    initTooltips();
});