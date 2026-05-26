!function () {
    var t = localStorage.getItem("theme");
    if (!t) {
        t = window.matchMedia("(prefers-color-scheme:dark)").matches ? "dark" : "light";
    }
    document.documentElement.setAttribute("data-bs-theme", t);
}()

// dark and light mode toggler
!function () {
    var key = 'theme';
    var stored = localStorage.getItem(key);
    
    // https://stackoverflow.com/a/57795495
    var prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    var theme = stored || (prefersDark ? 'dark' : 'light');
    document.documentElement.setAttribute('data-bs-theme', theme);

    var toggle = document.getElementById('themeToggle');
    var icon = document.getElementById('themeIcon');
    if (toggle) {
        icon.textContent = theme === 'dark' ? '☀️' : '🌖';
        toggle.addEventListener('click', function () {
            var current = document.documentElement.getAttribute('data-bs-theme');
            var next = current === 'dark' ? 'light' : 'dark';
            document.documentElement.setAttribute('data-bs-theme', next);
            localStorage.setItem(key, next);
            icon.textContent = next === 'dark' ? '☀️' : '🌖';
        });
    }
}()

// batch delete checkbox management
!function () {
    document.querySelectorAll('.week-summary').forEach(function (weekCard) {
        var selectAll = weekCard.querySelector('.select-all');
        var checkboxes = weekCard.querySelectorAll('.entry-checkbox');
        var deleteBtn = weekCard.querySelector('.delete-selected-btn');

        if (!selectAll || !checkboxes.length || !deleteBtn) return;

        function updateState() {
            var allChecked = true;
            var anyChecked = false;
            checkboxes.forEach(function (cb) {
                if (!cb.checked) allChecked = false;
                if (cb.checked) anyChecked = true;
            });
            selectAll.checked = allChecked;
            deleteBtn.disabled = !anyChecked;
        }

        selectAll.addEventListener('change', function () {
            checkboxes.forEach(function (cb) {
                cb.checked = selectAll.checked;
            });
            deleteBtn.disabled = !selectAll.checked;
        });

        checkboxes.forEach(function (cb) {
            cb.addEventListener('change', updateState);
        });
    });
}()