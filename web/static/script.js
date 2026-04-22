document.addEventListener("DOMContentLoaded", function() {
    // === Запускаем логику страницы результата (если она открыта) ===
    initBonuses();
    initEditParams();
    initDeductions();

    // === Ниже - только страница с формой расчёта ===
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
            text.removeAttribute('style');
            text.classList.remove("tooltip-above", "tooltip-below", "adjust-left", "adjust-right");

            const viewportWidth = window.innerWidth;
            const viewportHeight = window.innerHeight;
            const safeMargin = 12;
            const iconRect = icon.getBoundingClientRect();

            // fixed — все координаты относительно viewport, независимо от позиции родителя
            text.style.position = 'fixed';
            text.style.zIndex = '10000';
            text.style.boxSizing = 'border-box';
            text.style.transform = 'none';
            text.style.right = 'auto';
            text.style.bottom = 'auto';

            // Ширина: 300px или меньше если экран маленький
            const tooltipWidth = Math.min(300, viewportWidth - safeMargin * 2);
            text.style.width = tooltipWidth + 'px';

            // Горизонталь: центр на иконке, зажато в границы экрана
            const iconCenterX = iconRect.left + iconRect.width / 2;
            const left = Math.max(
                safeMargin,
                Math.min(iconCenterX - tooltipWidth / 2, viewportWidth - tooltipWidth - safeMargin)
            );
            text.style.left = left + 'px';

            // Измеряем высоту при позиции снизу
            text.style.top = (iconRect.bottom + 8) + 'px';
            void text.offsetWidth;

            const tooltipHeight = text.getBoundingClientRect().height;
            const spaceBelow = viewportHeight - iconRect.bottom;

            if (spaceBelow >= tooltipHeight + 20 || iconRect.top < 100) {
                text.classList.add('tooltip-below');
                text.style.top = (iconRect.bottom + 8) + 'px';
            } else {
                text.classList.add('tooltip-above');
                text.style.top = 'auto';
                text.style.bottom = (viewportHeight - iconRect.top + 8) + 'px';
            }
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

    // === Тип занятости: блок НПД и блокировка несовместимых опций ===
    function initEmploymentType() {
        var radios = document.querySelectorAll('input[name="employmentType"]');
        var privilegeCheckbox = document.querySelector('input[name="hasTaxPrivilege"]');
        var residentCheckbox = document.querySelector('input[name="isNotResident"]');
        var npdParams = document.getElementById("npd-params");
        if (!radios.length) return;

        function setDisabled(checkbox, disabled) {
            if (!checkbox) return;
            if (disabled) checkbox.checked = false;
            checkbox.disabled = disabled;
            var optionEl = checkbox.closest(".exclusive-option");
            if (optionEl) {
                optionEl.classList.toggle("exclusive-option--disabled", disabled);
                if (disabled) optionEl.classList.remove("selected");
            }
        }

        function updateCompatibility() {
            var selected = document.querySelector('input[name="employmentType"]:checked');
            var value = selected ? selected.value : "TD";
            var isGPH = value === "GPH";
            var isSelfEmployed = value === "SELF_EMPLOYED";

            // Блок НПД
            if (npdParams) {
                npdParams.classList.toggle("hidden", !isSelfEmployed);
            }

            // Льготы: недоступны при ГПХ и самозанятости
            setDisabled(privilegeCheckbox, isGPH || isSelfEmployed);

            // Нерезидент: недоступен при самозанятости (НПД только для резидентов РФ)
            setDisabled(residentCheckbox, isSelfEmployed);

            // Территориальный коэффициент и северная надбавка: только для ТД
            var hideCoeffs = isGPH || isSelfEmployed;
            ["territorialMultiplier", "northernCoefficient"].forEach(function(name) {
                var el = document.querySelector('[name="' + name + '"]');
                if (!el) return;
                var row = el.closest(".field-row");
                if (row) row.style.display = hideCoeffs ? "none" : "";
                if (hideCoeffs) el.value = "100";
            });
        }

        radios.forEach(function(radio) {
            radio.addEventListener("change", updateCompatibility);
        });

        updateCompatibility();
    }

    // Инициализация всех компонентов
    initExclusiveCheckboxes();
    initTooltips();
    initAccordion();
    initEmploymentType();
});

// ===================================================================
// СТРАНИЦА РЕЗУЛЬТАТА - помесячный аккордеон и ввод премий
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

    // --- Форматирование полей премий ---
    document.querySelectorAll(".bonus-input").forEach(initAmountInput);

    // --- Снимок начального состояния полей (чтобы знать, изменилось ли что-то) ---
    var initialBonusValues = {};
    document.querySelectorAll(".bonus-input").forEach(function(input) {
        initialBonusValues[input.name] = input.value;
    });

    // --- Обновление липкой панели ---
    function updateStickyBar() {
        var count = 0;
        var changed = false;
        document.querySelectorAll(".bonus-input").forEach(function(input) {
            var raw = input.value.replace(/\D/g, "");
            var val = parseInt(raw, 10);
            if (!isNaN(val) && val > 0) count++;
            if (input.value !== (initialBonusValues[input.name] || "")) changed = true;
        });

        var panelOpen = document.getElementById("edit-params-panel") &&
                        document.getElementById("edit-params-panel").classList.contains("open");

        if (changed) {
            var text = count > 0
                ? "Пересчитать с " + count + "\u00A0" + declension(count, ["премией", "премиями", "премиями"])
                : "Пересчитать без премий";
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
        // если panelOpen и бонусов нет - не трогаем (за это отвечает initEditParams)
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
            var v = parseInt(input.value.replace(/\D/g, ""), 10);
            if (!isNaN(v) && v > 0) hasBonuses = true;
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

    // Скрываем территориальный/северный при самозанятом в панели редактирования
    var editRadios = panel.querySelectorAll('input[name="employmentType"]');
    function updateEditPanelFields() {
        var sel = panel.querySelector('input[name="employmentType"]:checked');
        var isSE = sel && sel.value === "SELF_EMPLOYED";
        var isGPH = sel && sel.value === "GPH";

        // Территориальный и северный — только для ТД
        var hideCoeffs = isSE || isGPH;
        ["territorialMultiplier", "northernCoefficient"].forEach(function(name) {
            var el = panel.querySelector('[name="' + name + '"]');
            if (!el) return;
            var field = el.closest(".edit-field");
            if (field) field.style.display = hideCoeffs ? "none" : "";
        });

        // Особые условия (льготы и нерезидент) — скрываем всё поле при НПД,
        // при ГПХ оставляем видимым но отключаем только льготы
        var privilegeCb = panel.querySelector('input[name="hasTaxPrivilege"]');
        var residentCb  = panel.querySelector('input[name="isNotResident"]');

        if (isSE) {
            // Самозанятый: несовместимо ни с льготами, ни с нерезидентством
            [privilegeCb, residentCb].forEach(function(cb) {
                if (!cb) return;
                cb.checked = false;
                cb.disabled = true;
                var item = cb.closest(".edit-checkbox-item");
                if (item) item.style.opacity = "0.4";
            });
        } else if (isGPH) {
            // ГПХ: только льготы недоступны
            if (privilegeCb) {
                privilegeCb.checked = false;
                privilegeCb.disabled = true;
                var item = privilegeCb.closest(".edit-checkbox-item");
                if (item) item.style.opacity = "0.4";
            }
            if (residentCb) {
                residentCb.disabled = false;
                var item2 = residentCb.closest(".edit-checkbox-item");
                if (item2) item2.style.opacity = "";
            }
        } else {
            // ТД: всё доступно
            [privilegeCb, residentCb].forEach(function(cb) {
                if (!cb) return;
                cb.disabled = false;
                var item = cb.closest(".edit-checkbox-item");
                if (item) item.style.opacity = "";
            });
        }
    }
    editRadios.forEach(function(r) {
        r.addEventListener("change", updateEditPanelFields);
    });
    updateEditPanelFields();

    // Если после пересчёта были изменены параметры - открываем панель
    // (определяем по наличию query-параметра, который не нужен - панель закрыта по умолчанию)
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

    if (section.dataset.notResident) return;

    // --- Счётчик детей ---
    makeCounter("children-minus", "children-plus", "children-count");

    // --- Счётчик детей-инвалидов ---
    makeCounter("disabled-minus", "disabled-plus", "disabled-count");

    // --- Форматирование полей сумм ---
    section.querySelectorAll(".deduction-amount-input").forEach(initAmountInput);
}

// Форматирование числового поля: пробелы-разделители при blur, сырое число при focus.
function initAmountInput(input) {
    function formatValue(raw) {
        var n = parseInt(raw.replace(/\D/g, ""), 10);
        if (isNaN(n) || n === 0) return "";
        return n.toLocaleString("ru-RU");
    }

    function rawValue() {
        return input.value.replace(/\s/g, "").replace(/\u00A0/g, "");
    }

    // Форматируем предзаполненное значение сразу
    if (input.value) input.value = formatValue(input.value);

    input.addEventListener("input", function() {
        var pos = input.selectionStart;
        var before = input.value.length;
        // Разрешаем только цифры
        input.value = input.value.replace(/[^\d]/g, "");
        var diff = input.value.length - before;
        input.setSelectionRange(pos + diff, pos + diff);
    });

    input.addEventListener("blur", function() {
        input.value = formatValue(rawValue());
    });

    input.addEventListener("focus", function() {
        var raw = rawValue();
        input.value = raw === "0" ? "" : raw;
    });
}

// Вспомогательная функция: создаёт логику кнопок +/− для числового поля.
function makeCounter(minusBtnId, plusBtnId, inputId) {
    var minusBtn = document.getElementById(minusBtnId);
    var plusBtn  = document.getElementById(plusBtnId);
    var input    = document.getElementById(inputId);
    if (!input) return;

    if (minusBtn) {
        minusBtn.addEventListener("click", function() {
            var v = parseInt(input.value, 10) || 0;
            if (v > 0) input.value = v - 1;
        });
    }
    if (plusBtn) {
        plusBtn.addEventListener("click", function() {
            var v = parseInt(input.value, 10) || 0;
            var max = parseInt(input.max, 10) || 10;
            if (v < max) input.value = v + 1;
        });
    }
}