if ("serviceWorker" in navigator) {
    window.addEventListener("load", function () {
        this.navigator.serviceWorker
            .register("/js/worker.js")
            .then(res => console.log("service worker registered"))
            .catch(err => console.log("Error registering service worker: ", err))
    })
}

document.getElementById('add-item')
    .addEventListener('click', addItem);

function addItem() {
    var shoppingList = document.getElementById('shopping-list');
    const item = newItem();
    shoppingList.insertBefore(item, shoppingList.lastElementChild);
    item.querySelector('input[type=text]').focus();
}

function newItem() {
    const itemId = document.createElement('input');
    itemId.type = 'hidden';
    itemId.classList.add('id');
    const itemInput = document.createElement('input');
    itemInput.type = 'text';
    itemInput.addEventListener('focusout', finishEditItem);
    itemInput.addEventListener('focusin', startEditItem);
    const itemQuantity = document.createElement('input');
    itemQuantity.type = 'text';
    const item = document.createElement('article');
    item.classList.add('shopping-item');
    item.addEventListener('click', function(e) {
        if (e.target.tagName !== 'INPUT') {
            e.target.querySelector('.item input').focus();
        }
    });
    item.appendChild(itemId);
    item.appendChild(label('Item', 'item', itemInput));
    item.appendChild(label('Quantity', 'quantity', itemQuantity));
    return item;
}

function label(text, classes, child) {
    const label = document.createElement('label');
    label.textContent = text;
    label.classList.add(classes);
    label.appendChild(child);
    return label;
}

function startEditItem() {
    document.getElementById('add-item').disabled = false;
}

function finishEditItem(e) {
    document.getElementById('add-item').disabled = true;
    if(e.target.value && Array.from(document.querySelectorAll(".shopping-item > .item > input")).pop().value) {
        document.getElementById('add-item').disabled = false;
        const article = e.target.closest('article');
        const item = {
            name: article.querySelector('.item > input').value,
            quantity: article.querySelector('.quantity > input').value
        }
        idInput = article.querySelector('.id');
        
        if (idInput.value) {
            item.id = Number(idInput.value);
        }
        
        httpPut('/api/items', item, (r) => { idInput.value = JSON.parse(r.target.responseText).id; });
    }
}

function generateId() {
    return Math.random().toString(36).slice(2);
}

function httpPut(path, body, callback) {
    var r = new XMLHttpRequest();
    r.open("PUT", path, true);
    r.onreadystatechange = function () {
        if (r.readyState != 4 || r.status != 200 || r.status != 201) return;
    };
    r.onload = callback;
    r.send(JSON.stringify(body));
}