const urlParams = new URLSearchParams(window.location.search);
const symbol = urlParams.get('symbol');

document.getElementById('title').innerText = `Криптовалюта: ${symbol}`;

// Цена
fetch(`/cryptos/${symbol}/price`)
    .then(res => res.json())
    .then(data => {
        document.getElementById('price').innerText = `Цена: $${data.price}`;
    })
    .catch(() => {
        document.getElementById('price').innerText = 'Цена недоступна';
    });

// Сравнение
document.getElementById('compare-btn').addEventListener('click', () => {
    fetch(`/cryptos/${symbol}/compare`)
        .then(res => res.json())
        .then(data => {
            const div = document.getElementById('compare-result');
            div.innerHTML = '<h2 class="font-semibold">Результат:</h2>';
            data.forEach(row => {
                const p = document.createElement('p');
                p.textContent = `${row.exchange}: $${row.price}`;
                div.appendChild(p);
            });
        });
});
