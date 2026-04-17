document.addEventListener("DOMContentLoaded", function() {
    // === Запускаем логику страницы результата (если она открыта) ===
    initBonuses();
    initEditParams();
    initDeductions();

    // === Ниже — только страница с формой расчёта ===
    const salaryInput = document.getElementById("grossSalary");
    if (!salaryInput) return;

    const form = document.querySelector("form.form");

    const MIN_SALARY = window.APP_CONFIG?.minSalary || 0;
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

    salaryInput.addEventListener("input", clearMessages);

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

        if (isValid) {
            this.submit();
        }
    });

    // === Подсказки ===
function initTooltips() {
    const tooltips = document.querySelectorAll(".field-tooltip");

    tooltips.forEach(tooltip => {
        const icon = tooltip.querySelector(".tooltip-icon");
        const text = tooltip.querySelector(".tooltip-text");
        if (!icon || !text) return;

        function adjustTooltipPosition() {
            // ПЕРЕСЧИТЫВАЕМ isMobile ПРИ КАЖДОМ ПОКАЗЕ
            const isMobile = window.matchMedia("(max-width: 768px)").matches;
            
            if (isMobile) {
                adjustMobileTooltipPosition(text, icon);
            } else {
                adjustDesktopTooltipPosition(text, icon);
            }
        }

        function adjustMobileTooltipPosition(text, icon) {
            // Сбрасываем ВСЕ стили
            text.removeAttribute('style');
            text.classList.remove("tooltip-above", "tooltip-below", "adjust-left", "adjust-right");

            // Применяем только нужные стили
            text.style.position = 'absolute';
            text.style.zIndex = '10000';
            text.style.width = '300px';
            text.style.maxWidth = '300px';
            text.style.boxSizing = 'border-box';
            text.style.left = '50%';
            text.style.transform = 'translateX(-50%)';

            void text.offsetWidth; // форсируем рефлоу

            const viewportWidth = window.innerWidth;
            const viewportHeight = window.innerHeight;

            // ПОЛУЧАЕМ РАЗМЕРЫ ЭЛЕМЕНТОВ
            const iconRect = icon.getBoundingClientRect();
            const tooltipRect = text.getBoundingClientRect();

            // Определяем доступное пространство
            const spaceBelow = viewportHeight - iconRect.bottom;
            const spaceAbove = iconRect.top;
            const tooltipHeight = tooltipRect.height;

            // Выбираем позицию (сверху или снизу)
            if (spaceBelow >= tooltipHeight + 20 || spaceAbove < 100) {
                // Показываем снизу
                text.classList.add('tooltip-below');
                text.style.top = 'calc(100% + 8px)';
            } else {
                // Показываем сверху
                text.classList.add('tooltip-above');
                text.style.bottom = 'calc(100% + 8px)';
            }

            // РАСЧЕТ ГОРИЗОНТАЛЬНОЙ ПОЗИЦИИ С УЧЕТОМ ГРАНИЦ ЭКРАНА
            const iconCenterX = iconRect.left + (iconRect.width / 2);
            const tooltipWidth = 300; // Фиксированная ширина
            let desiredLeft = iconCenterX - (tooltipWidth / 2);

            // Корректируем позицию чтобы не выходить за экран
            const safeMargin = 15;

            // Если тултип не помещается по ширине, уменьшаем его
            if (tooltipWidth > viewportWidth - safeMargin * 2) {
                const newWidth = viewportWidth - safeMargin * 2;
                text.style.width = newWidth + 'px';
                text.style.maxWidth = newWidth + 'px';
            }

            if (desiredLeft < safeMargin) {
                // Выравниваем по левому краю с отступом
                text.classList.add('adjust-left');
                text.style.left = safeMargin + 'px';
                text.style.transform = 'none';
            } else if (desiredLeft + tooltipWidth > viewportWidth - safeMargin) {
                // Выравниваем по правому краю с отступом
                text.classList.add('adjust-right');
                text.style.left = 'auto';
                text.style.right = safeMargin + 'px';
                text.style.transform = 'none';
            } else {
                // Центрируем относительно иконки
                text.style.left = '50%';
                text.style.transform = 'translateX(-50%)';
            }

            // ФИНАЛЬНАЯ ПРОВЕРКА
            requestAnimationFrame(() => {
                const finalRect = text.getBoundingClientRect();

                // Если все равно выходит за край, принудительно корректируем ширину
                if (finalRect.right > viewportWidth - 5) {
                    const overflow = finalRect.right - viewportWidth + 5;
                    const newWidth = tooltipWidth - overflow;
                    text.style.width = newWidth + 'px';
                }

                if (finalRect.left < 5) {
                    const overflow = 5 - finalRect.left;
                    const newWidth = tooltipWidth - overflow;
                    text.style.width = newWidth + 'px';
                    if (text.classList.contains('adjust-left')) {
                        text.style.left = '5px';
                    }
                }
            });
        }

        function adjustDesktopTooltipPosition(text, icon) {
            text.classList.remove("tooltip-above", "tooltip-below");
            text.style.left = "50%";
            text.style.right = "auto";
            text.style.top = "auto";
            text.style.bottom = "auto";
            text.style.transform = "translateX(-50%)";
            text.style.position = "absolute";
            text.style.width = "280px";
            text.style.maxWidth = "280px";
            text.style.zIndex = "1001";

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

        const isTouch = window.matchMedia("(hover: none)").matches;

        // 🖱️ Hover (desktop)
        if (!isTouch) {
            tooltip.addEventListener("mouseenter", () => {
                text.classList.add("visible");
                requestAnimationFrame(() => adjustTooltipPosition());
            });
            tooltip.addEventListener("mouseleave", () => {
                text.classList.remove("visible");
                document.body.style.overflowX = '';
            });
        }

        // 📱 Touch (mobile) - ПРОСТОЙ ВАРИАНТ
else {
    icon.addEventListener("click", e => {
        e.stopPropagation();
        const isVisible = text.classList.contains("visible");

        document
            .querySelectorAll(".tooltip-text.visible")
            .forEach(el => {
                el.classList.remove("visible");
                document.body.classList.remove("tooltip-mobile-open");
            });

        if (!isVisible) {
            text.classList.add("visible");
            document.body.classList.add("tooltip-mobile-open");
            requestAnimationFrame(() => adjustTooltipPosition());
        } else {
            document.body.classList.remove("tooltip-mobile-open");
        }
    });

    // ЗАКРЫВАЕМ ТУЛТИПЫ ПРИ ЛЮБОМ ДВИЖЕНИИ ПАЛЬЦЕМ
    document.addEventListener('touchmove', function() {
        if (document.body.classList.contains('tooltip-mobile-open')) {
            document
                .querySelectorAll(".tooltip-text.visible")
                .forEach(el => {
                    el.classList.remove("visible");
                    document.body.classList.remove("tooltip-mobile-open");
                });
        }
    }, { passive: true });

    document.addEventListener("click", () => {
        document
            .querySelectorAll(".tooltip-text.visible")
            .forEach(el => {
                el.classList.remove("visible");
                document.body.classList.remove("tooltip-mobile-open");
            });
    });
}

        // Пересчитываем при изменении размера
        window.addEventListener('resize', () => {
            if (text.classList.contains('visible')) {
                requestAnimationFrame(() => adjustTooltipPosition());
            }
        });

        window.addEventListener('orientationchange', () => {
            setTimeout(() => {
                if (text.classList.contains('visible')) {
                    requestAnimationFrame(() => adjustTooltipPosition());
                }
            }, 100);
        });
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
            } else {
                option.classList.remove("selected");
            }
        });
    });
}


    // === Аккордеон ===
    function initAccordion() {
        const toggleBtn = document.querySelector(".additional-toggle");
        const additionalBlock = document.getElementById("additional-params");

        if (toggleBtn && additionalBlock) {
            toggleBtn.addEventListener("click", () => {
                const isOpen = additionalBlock.classList.toggle("open");
                toggleBtn.classList.toggle("open");

                if (isOpen) {
                    additionalBlock.style.maxHeight = additionalBlock.scrollHeight + "px";
                    toggleBtn.style.setProperty('--rotation', '180deg');
                } else {
                    additionalBlock.style.maxHeight = null;
                    toggleBtn.style.setProperty('--rotation', '0deg');
                }
            });
        }
    }

    // Инициализация всех компонентов
    initExclusiveCheckboxes();
    initTooltips();
    initAccordion();
});

// ===================================================================
// СТРАНИЦА РЕЗУЛЬТАТА — помесячный аккордеон и ввод премий
// ===================================================================
function initBonuses() {
    const stickyBar = document.getElementById("bonus-sticky-bar");
    if (!stickyBar) return; // не страница результата

    const stickyLabel = document.getElementById("bonus-sticky-label");
    const inlineRecalcBtn = document.getElementById("bonus-recalc-btn");

    // --- Аккордеон помесячной детализации ---
    document.querySelectorAll(".accordion-header").forEach(function(header) {
        header.addEventListener("click", function(e) {
            // Игнорируем клик по кнопкам внутри заголовка
            if (e.target.closest("button")) return;
            header.parentElement.classList.toggle("active");
        });
    });

    // --- Обработчик кнопки «Добавить премию» ---
    function onAddBtnClick() {
        var section = this.closest(".bonus-section");
        var inputRow = section.querySelector(".bonus-input-row");
        this.style.display = "none";
        inputRow.classList.remove("bonus-input-row--hidden");
        var input = inputRow.querySelector(".bonus-input");
        if (input) input.focus();
    }

    document.querySelectorAll(".bonus-add-btn").forEach(function(btn) {
        btn.addEventListener("click", onAddBtnClick);
    });

    // --- Обработчик кнопки × (убрать премию) ---
    document.querySelectorAll(".bonus-clear-btn").forEach(function(btn) {
        btn.addEventListener("click", function() {
            var section = btn.closest(".bonus-section");
            var inputRow = btn.closest(".bonus-input-row");
            var input = inputRow.querySelector(".bonus-input");

            if (input) input.value = "";
            inputRow.classList.add("bonus-input-row--hidden");

            // Показываем кнопку «Добавить премию» или создаём её (если была преднаполненная строка)
            var addBtn = section.querySelector(".bonus-add-btn");
            if (!addBtn) {
                addBtn = document.createElement("button");
                addBtn.type = "button";
                addBtn.className = "bonus-add-btn";
                addBtn.textContent = "＋ Добавить премию";
                addBtn.addEventListener("click", onAddBtnClick);
                section.insertBefore(addBtn, inputRow);
            }
            addBtn.style.display = "";

            updateStickyBar();
        });
    });

    // --- Слушаем ввод во всех полях премий ---
    document.addEventListener("input", function(e) {
        if (e.target.classList.contains("bonus-input")) {
            updateStickyBar();
        }
    });

    // --- Обновление липкой панели ---
    function updateStickyBar() {
        var count = 0;
        document.querySelectorAll(".bonus-input").forEach(function(input) {
            var val = parseFloat(input.value);
            if (!isNaN(val) && val > 0) count++;
        });

        var panelOpen = document.getElementById("edit-params-panel") &&
                        document.getElementById("edit-params-panel").classList.contains("open");

        if (count > 0) {
            var label = declension(count, ["премией", "премиями", "премиями"]);
            var text = "Пересчитать с " + count + "\u00A0" + label;
            stickyBar.classList.add("visible");
            stickyLabel.textContent = text;
            if (inlineRecalcBtn) {
                inlineRecalcBtn.style.display = "inline-block";
                inlineRecalcBtn.textContent = text + " →";
            }
        } else if (!panelOpen) {
            stickyBar.classList.remove("visible");
            if (inlineRecalcBtn) inlineRecalcBtn.style.display = "";
        }
        // если panelOpen и бонусов нет — не трогаем (за это отвечает initEditParams)
    }

    // --- Склонение числительных ---
    function declension(n, forms) {
        var mod10 = n % 10;
        var mod100 = n % 100;
        if (mod100 >= 11 && mod100 <= 19) return forms[2];
        if (mod10 === 1) return forms[0];
        if (mod10 >= 2 && mod10 <= 4) return forms[1];
        return forms[2];
    }

    // Инициализация: при пересчёте с бонусами поля уже заполнены
    updateStickyBar();
}

// ===================================================================
// ПАНЕЛЬ РЕДАКТИРОВАНИЯ ПАРАМЕТРОВ
// ===================================================================
function initEditParams() {
    var toggle = document.getElementById("edit-params-toggle");
    var panel  = document.getElementById("edit-params-panel");
    var cancel = document.getElementById("edit-params-cancel");
    if (!toggle || !panel) return;

    function openPanel() {
        panel.classList.add("open");
        toggle.classList.add("active");
        toggle.textContent = "✏️ Параметры открыты";
        showRecalcButtons("Пересчитать");
    }

    function closePanel() {
        panel.classList.remove("open");
        toggle.classList.remove("active");
        toggle.textContent = "✏️ Изменить параметры";
        // Скрываем кнопки только если нет заполненных бонусов
        var hasBonuses = false;
        document.querySelectorAll(".bonus-input").forEach(function(input) {
            if (!isNaN(parseFloat(input.value)) && parseFloat(input.value) > 0) hasBonuses = true;
        });
        if (!hasBonuses) hideRecalcButtons();
    }

    function showRecalcButtons(label) {
        var stickyBar   = document.getElementById("bonus-sticky-bar");
        var stickyLabel = document.getElementById("bonus-sticky-label");
        var inlineBtn   = document.getElementById("bonus-recalc-btn");
        if (stickyBar)   stickyBar.classList.add("visible");
        if (stickyLabel) stickyLabel.textContent = label;
        if (inlineBtn) {
            inlineBtn.style.display = "inline-block";
            inlineBtn.textContent = label + " →";
        }
    }

    function hideRecalcButtons() {
        var stickyBar = document.getElementById("bonus-sticky-bar");
        var inlineBtn = document.getElementById("bonus-recalc-btn");
        if (stickyBar) stickyBar.classList.remove("visible");
        if (inlineBtn) inlineBtn.style.display = "";
    }

    toggle.addEventListener("click", function() {
        if (panel.classList.contains("open")) {
            closePanel();
        } else {
            openPanel();
        }
    });

    if (cancel) {
        cancel.addEventListener("click", closePanel);
    }

    // Взаимоисключающие чекбоксы внутри панели
    var exclusives = panel.querySelectorAll(".edit-exclusive");
    exclusives.forEach(function(cb) {
        cb.addEventListener("change", function() {
            if (this.checked) {
                exclusives.forEach(function(other) {
                    if (other !== cb) other.checked = false;
                });
            }
        });
    });

    // Если после пересчёта были изменены параметры — открываем панель
    // (определяем по наличию query-параметра, который не нужен — панель закрыта по умолчанию)
}

// ===================================================================
// БЛОК НАЛОГОВЫХ ВЫЧЕТОВ
// ===================================================================
function initDeductions() {
    var section = document.querySelector(".deductions-section");
    if (!section) return;

    var toggle  = document.getElementById("toggle-deductions");
    var content = document.getElementById("deductions-content");
    var arrow   = toggle ? toggle.querySelector(".deductions-arrow") : null;

    if (toggle && content) {
        toggle.addEventListener("click", function() {
            var isOpen = content.classList.toggle("open");
            if (arrow) arrow.style.transform = isOpen ? "rotate(180deg)" : "";
            toggle.classList.toggle("open", isOpen);
        });
    }

    // Если нерезидент — интерактивная часть не нужна
    if (section.dataset.notResident) return;

    // Данные из атрибутов (копейки)
    var annualTaxKop     = parseInt(section.dataset.annualTax,     10) || 0;
    var monthlyGrossKop  = parseInt(section.dataset.monthlyGross,  10) || 0;

    // --- Элементы ---
    var childrenCountEl  = document.getElementById("children-count");
    var disabledChildEl  = document.getElementById("disabled-child");
    var childrenMinusEl  = document.getElementById("children-minus");
    var childrenPlusEl   = document.getElementById("children-plus");

    var housingEl        = document.getElementById("housing-amount");
    var mortgageEl       = document.getElementById("mortgage-amount");
    var socialEl         = document.getElementById("social-amount");
    var childEduEl       = document.getElementById("child-edu-amount");

    var childrenResultEl = document.getElementById("children-result");
    var propertyResultEl = document.getElementById("property-result");
    var socialResultEl   = document.getElementById("social-result");
    var totalEl          = document.getElementById("deductions-total");
    var totalValueEl     = document.getElementById("total-return-value");
    var totalNoteEl      = document.getElementById("total-return-note");

    // --- Счётчик детей ---
    if (childrenMinusEl && childrenCountEl) {
        childrenMinusEl.addEventListener("click", function() {
            var v = parseInt(childrenCountEl.value, 10) || 0;
            if (v > 0) childrenCountEl.value = v - 1;
            recalc();
        });
    }
    if (childrenPlusEl && childrenCountEl) {
        childrenPlusEl.addEventListener("click", function() {
            var v = parseInt(childrenCountEl.value, 10) || 0;
            if (v < 10) childrenCountEl.value = v + 1;
            recalc();
        });
    }

    // --- Слушаем все поля ---
    [childrenCountEl, disabledChildEl, housingEl, mortgageEl, socialEl, childEduEl].forEach(function(el) {
        if (el) el.addEventListener("change", recalc);
        if (el && el.type === "number") el.addEventListener("input", recalc);
    });

    // --- Форматирование ---
    function fmtRub(rubles) {
        return Math.round(rubles).toLocaleString("ru-RU") + "\u00A0\u20BD";
    }

    function renderResultBlock(container, rows, noteText) {
        if (!container) return;
        if (!rows || rows.length === 0) {
            container.innerHTML = "";
            return;
        }
        var html = rows.map(function(row) {
            return '<div class="deduction-result-row">' +
                '<span class="deduction-result-label">' + row.label + '</span>' +
                '<span class="deduction-result-amount">' + row.value + '</span>' +
                '</div>';
        }).join("");
        if (noteText) {
            html += '<div class="deduction-result-note">' + noteText + '</div>';
        }
        container.innerHTML = html;
    }

    // --- Расчёт вычета на детей (ст. 218) ---
    function calcChildren() {
        var count      = parseInt((childrenCountEl && childrenCountEl.value) || 0, 10);
        var hasDisabled = disabledChildEl && disabledChildEl.checked;
        if (count === 0 && !hasDisabled) {
            renderResultBlock(childrenResultEl, []);
            return 0;
        }

        // Ежемесячная сумма вычета (рубли)
        var monthlyDed = 0;
        for (var i = 1; i <= count; i++) {
            if      (i === 1) monthlyDed += 1400;
            else if (i === 2) monthlyDed += 2800;
            else              monthlyDed += 6000;
        }
        if (hasDisabled && count > 0) monthlyDed += 12000;

        // Месяцев до превышения лимита 450 000 ₽
        var monthlyGrossRub = monthlyGrossKop / 100;
        var months = monthlyGrossRub > 0 ? Math.min(12, Math.floor(450000 / monthlyGrossRub)) : 12;
        if (months <= 0) months = 0;

        var totalDed  = monthlyDed * months;
        var taxReturn = Math.min(Math.round(totalDed * 0.13), annualTaxKop / 100);

        var note = months < 12
            ? "Вычет действует " + months + "\u00A0мес. (до превышения лимита дохода 450\u00A0000\u00A0\u20BD)"
            : null;

        renderResultBlock(childrenResultEl, [
            { label: "Ежемесячный вычет:", value: fmtRub(monthlyDed) },
            { label: "Возврат налога за период:", value: fmtRub(taxReturn) }
        ], note);

        return taxReturn;
    }

    // --- Расчёт имущественного вычета (ст. 220) ---
    function calcProperty() {
        var housing  = parseFloat((housingEl  && housingEl.value)  || 0) || 0;
        var mortgage = parseFloat((mortgageEl && mortgageEl.value) || 0) || 0;
        if (housing <= 0 && mortgage <= 0) {
            renderResultBlock(propertyResultEl, []);
            return 0;
        }

        var housingDed  = Math.min(housing,  2000000);
        var mortgageDed = Math.min(mortgage, 3000000);
        var maxReturn   = Math.round((housingDed + mortgageDed) * 0.13);

        var annualTaxRub  = annualTaxKop / 100;
        var thisYearReturn = Math.min(maxReturn, annualTaxRub);

        var rows = [];
        if (housing > 0)  rows.push({ label: "Возврат за жильё (до 260\u00A0000\u00A0\u20BD):", value: fmtRub(Math.round(housingDed * 0.13)) });
        if (mortgage > 0) rows.push({ label: "Возврат за ипотеку (до 390\u00A0000\u00A0\u20BD):", value: fmtRub(Math.round(mortgageDed * 0.13)) });

        var note = null;
        if (maxReturn > annualTaxRub) {
            note = "Остаток " + fmtRub(maxReturn - thisYearReturn) + " можно перенести на следующие годы";
        }

        renderResultBlock(propertyResultEl, rows, note);
        return thisYearReturn;
    }

    // --- Расчёт социального вычета (ст. 219) ---
    function calcSocial() {
        var social   = parseFloat((socialEl   && socialEl.value)   || 0) || 0;
        var childEdu = parseFloat((childEduEl && childEduEl.value) || 0) || 0;
        if (social <= 0 && childEdu <= 0) {
            renderResultBlock(socialResultEl, []);
            return 0;
        }

        var socialDed   = Math.min(social,   150000);
        var childEduDed = Math.min(childEdu, 110000);
        var totalReturn = Math.round((socialDed + childEduDed) * 0.13);
        var capped      = Math.min(totalReturn, annualTaxKop / 100);

        var rows = [];
        if (social > 0)   rows.push({ label: "Лечение и своё обучение (до 19\u00A0500\u00A0\u20BD):", value: fmtRub(Math.round(socialDed * 0.13)) });
        if (childEdu > 0) rows.push({ label: "Обучение ребёнка (до 14\u00A0300\u00A0\u20BD):",       value: fmtRub(Math.round(childEduDed * 0.13)) });

        renderResultBlock(socialResultEl, rows, null);
        return capped;
    }

    // --- Пересчёт итога ---
    function recalc() {
        var r1 = calcChildren();
        var r2 = calcProperty();
        var r3 = calcSocial();

        var total = r1 + r2 + r3;
        if (total <= 0) {
            if (totalValueEl) totalValueEl.textContent = "\u2014";
            if (totalNoteEl)  totalNoteEl.textContent  = "";
            return;
        }

        // Итог не может превышать уплаченный НДФЛ
        var annualTaxRub = annualTaxKop / 100;
        var capped = Math.min(total, annualTaxRub);

        if (totalValueEl) totalValueEl.textContent = fmtRub(capped);
        if (totalNoteEl) {
            totalNoteEl.textContent = total > annualTaxRub
                ? "Ограничено уплаченным НДФЛ. Имущественный вычет можно перенести на следующие годы."
                : "";
        }
    }
}