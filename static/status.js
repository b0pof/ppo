function load() {
    fetch('/status/raw')
        .then(r => r.text())
        .then(t => {
            document.getElementById('data').innerText = t;
        });
}

load();
setInterval(load, 2000);
