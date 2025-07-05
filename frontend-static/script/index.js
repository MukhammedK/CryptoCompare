fetch('/cryptos')
    .then(res => res.json())
    .then(data => {
        const container = document.getElementById('crypto-list');
        data.forEach(c => {
            const div = document.createElement('div');
            div.className = 'p-4 border rounded shadow';
            div.innerHTML = `
        <h2 class="font-semibold">${c.name} (${c.symbol})</h2>
        <a class="text-blue-600" href="crypto.html?symbol=${c.symbol}">Подробнее</a>
      `;
            container.appendChild(div);
        });
    });
