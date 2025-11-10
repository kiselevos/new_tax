document.addEventListener("DOMContentLoaded", function() {
    const form = document.querySelector("form");
    const salaryInput = document.getElementById("grossSalary");

    if (!form || !salaryInput) return;

    const MIN_SALARY = window.APP_CONFIG.minSalary;
    const MAX_SALARY = 100000000;
    const salaryContainer = salaryInput.closest(".salary-field-container");

    // === Сообщения об ошибках ===
    let messageEl = salaryContainer.querySelector(".field-message");
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
        salaryInput.style.animation = "shake 0.3s ease";
        setTimeout(() => (salaryInput.style.animation = ""), 300);
    }

    function clearMessages() {
        messageEl.textContent = "";
        messageEl.className = "field-message";
        salaryInput.classList.remove("field-error");
    }

    // === Разрешаем только цифры, запятую и точку ===
    salaryInput.addEventListener("input", function () {
        let value = this.value;

        value = value.replace(/[^0-9.,]/g, ""); // только цифры и разделители

        const firstSeparator = value.match(/[.,]/);
        if (firstSeparator) {
            const sep = firstSeparator[0];
            const pos = firstSeparator.index;
            value = value.replace(/[.,]/g, ""); // удаляем все разделители
            value = value.slice(0, pos) + sep + value.slice(pos); // вставляем первый
        }

        // максимум 2 цифры после разделителя
        value = value.replace(/^(\d+)([.,])(\d{0,2}).*$/, (_, int, sep, frac) =>
            frac !== undefined ? int + sep + frac : int
        );

        this.value = value;
    });

    // === Форматирование при blur (всегда 2 копейки) ===
    salaryInput.addEventListener("blur", function () {
        let value = this.value.trim().replace(/\s/g, "").replace(",", ".");
        const number = parseFloat(value);

        if (!isNaN(number)) {
            this.value = number.toLocaleString("ru-RU", {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2,
            });
        }
    });

    // === При фокусе убираем форматирование ===
    salaryInput.addEventListener("focus", function () {
        this.value = this.value.replace(/\s/g, "").replace(",", ".");
    });

    // === Проверка при отправке формы ===
    form.addEventListener("submit", function (e) {
        e.preventDefault();

        const rawValue = salaryInput.value
            .trim()
            .replace(/\s/g, "")
            .replace(",", ".");
        const salary = parseFloat(rawValue);

        clearMessages();
        let isValid = true;

        if (!rawValue || isNaN(salary)) {
            showError("Введите корректный оклад (например, 50000)");
            isValid = false;
        } else if (salary < MIN_SALARY) {
            showError(`Минимальный оклад - ${MIN_SALARY.toLocaleString("ru-RU")} ₽`);
            isValid = false;
        } else if (salary > MAX_SALARY) {
            showError(
                `Сумма слишком велика. Максимум - ${MAX_SALARY.toLocaleString("ru-RU")} ₽`
            );
            isValid = false;
        }

        if (isValid) form.submit();
    });

    // === Подсказки ===
    function initTooltips() {
        const tooltips = document.querySelectorAll(".field-tooltip");

        tooltips.forEach(tooltip => {
            const icon = tooltip.querySelector(".tooltip-icon");
            const text = tooltip.querySelector(".tooltip-text");
            if (!icon || !text) return;

            const isTouch = window.matchMedia("(hover: none)").matches;

            function adjustTooltipPosition() {
                // сбрасываем стили
                text.classList.remove("tooltip-above", "tooltip-below");
                text.style.left = "50%";
                text.style.right = "auto";
                text.style.top = "auto";
                text.style.bottom = "auto";
                text.style.transform = "translateX(-50%)";

                // форсируем рефлоу
                void text.offsetWidth;

                const iconRect = icon.getBoundingClientRect();
                const spaceAbove = iconRect.top;
                const spaceBelow = window.innerHeight - iconRect.bottom;
                const prefersAbove = spaceAbove > spaceBelow;

                if (prefersAbove) {
                    text.classList.add("tooltip-above");
                    text.style.bottom = "calc(100% + 10px)";
                } else {
                    text.classList.add("tooltip-below");
                    text.style.top = "calc(100% + 10px)";
                }

                // проверяем выход за края
                const rect = text.getBoundingClientRect();
                if (rect.left < 8) {
                    text.style.left = "0";
                    text.style.transform = "none";
                } else if (rect.right > window.innerWidth - 8) {
                    text.style.left = "auto";
                    text.style.right = "0";
                    text.style.transform = "none";
                }
            }

            // 🖱️ Hover (desktop)
            if (!isTouch) {
                tooltip.addEventListener("mouseenter", () => {
                    text.classList.add("visible");
                    requestAnimationFrame(() => adjustTooltipPosition());
                });
                tooltip.addEventListener("mouseleave", () => {
                    text.classList.remove("visible");
                });
            }

            // 📱 Touch (mobile)
            else {
                icon.addEventListener("click", e => {
                    e.stopPropagation();
                    const isVisible = text.classList.contains("visible");

                    document
                        .querySelectorAll(".tooltip-text.visible")
                        .forEach(el => el.classList.remove("visible"));

                    if (!isVisible) {
                        text.classList.add("visible");
                        requestAnimationFrame(() => adjustTooltipPosition());
                    }
                });

                document.addEventListener("click", () => {
                    document
                        .querySelectorAll(".tooltip-text.visible")
                        .forEach(el => el.classList.remove("visible"));
                });
            }
        });
    }

    // === Взаимоисключающие чекбоксы ===
    function initExclusiveCheckboxes() {
        const options = document.querySelectorAll(".exclusive-option");
        options.forEach(option => {
            const checkbox = option.querySelector('input[type="checkbox"]');
            if (!checkbox) return;

            checkbox.addEventListener("change", function () {
                if (this.checked) {
                    option.classList.add("selected");

                    options.forEach(otherOption => {
                        if (otherOption !== option) {
                            const otherCheckbox = otherOption.querySelector(
                                'input[type="checkbox"]'
                            );
                            if (otherCheckbox) {
                                otherCheckbox.checked = false;
                                otherOption.classList.remove("selected");
                            }
                        }
                    });
                    showCheckboxHint(this);
                } else {
                    option.classList.remove("selected");
                    removeCheckboxHint();
                }
            });
        });
    }

    function showCheckboxHint(checkedCheckbox) {
        removeCheckboxHint();
        const hint = document.createElement("div");
        hint.className = "checkbox-hint";

        if (checkedCheckbox.name === "hasTaxPrivilege") {
            hint.textContent = "✓ Выбраны льготы для силовых структур";
        } else if (checkedCheckbox.name === "isNotResident") {
            hint.textContent = "✓ Выбран статус налогового нерезидента";
        }

        const container = checkedCheckbox.closest(".exclusive-checkboxes");
        if (container) container.appendChild(hint);
    }

    function removeCheckboxHint() {
        document.querySelectorAll(".checkbox-hint").forEach(hint => hint.remove());
    }

    initExclusiveCheckboxes();
    initTooltips();
});
