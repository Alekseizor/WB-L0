document.querySelector('button').onclick = myClick;

function myClick() {
    event.preventDefault(); // Предотвращение действия по умолчанию
    var xhr = new XMLHttpRequest();
    var inputValue = document.getElementById("IDInput").value;
    xhr.open("GET", "http://localhost:8080/api/orders/" + inputValue, true);
    xhr.setRequestHeader("Access-Control-Allow-Origin", "http://localhost:8080");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            var response = JSON.parse(xhr.responseText);

            // Найдите таблицу в HTML
            var table = document.querySelector('table');

// Очистить существующие строки таблицы (кроме заголовка)
            while (table.rows.length > 1) {
                table.deleteRow(1);
            }

// Создать новую строку
            var row = table.insertRow();

// Добавить ячейки с данными в строку
            for (var key in response) {
                if (response.hasOwnProperty(key)) {
                    var cell = row.insertCell();
                    cell.innerHTML = response[key];
                }
            }
        } else {
            console.error('Произошла ошибка при выполнении запроса.');
        }
    };

    xhr.onerror = function () {
        console.error('Произошла ошибка при выполнении запроса.');
    };

    xhr.send();
}