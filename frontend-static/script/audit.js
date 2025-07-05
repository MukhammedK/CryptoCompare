fetch('/audit')
    .then(res => res.json())
    .then(data => {
        const container = document.getElementById('audit-log');

        data.forEach(entry => {
            const div = document.createElement('div');
            div.className = 'border p-4 rounded shadow';

            const date = new Date(entry.timestamp).toLocaleString();

            div.innerHTML = `<h2 class="font-semibold mb-2">${entry.symbol} — ${date}</h2>`;
            entry.results.forEach(r => {
                const p = document.createElement('p');
                p.textContent = `${r.exchange}: $${r.price}`;
                div.appendChild(p);
            });

            container.appendChild(div);
        });
    });
