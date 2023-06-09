if ("serviceWorker" in navigator) {
    window.addEventListener("load", function () {
        this.navigator.serviceWorker
            .register("/js/worker.js")
            .then(res => console.log("service worker registered"))
            .catch(err => console.log("Error registering service worker: ", err))
    })
}

window.addEventListener('load', function() {
    const shoppingList = document.getElementById('shopping-list');
    httpGet('/api/items', (r) => { 
        JSON.parse(r.target.responseText).forEach((i) => {
            const item = newItem(i.id, i.name, i.quantity, new Date(i.dateAdded));
            shoppingList.insertBefore(item, shoppingList.lastElementChild);
        });
    })
});

document.getElementById('add-item')
    .addEventListener('click', addItem);

function addItem() {
    const shoppingList = document.getElementById('shopping-list');
    const item = newItem();
    shoppingList.insertBefore(item, shoppingList.lastElementChild);
    item.querySelector('input[type=text]').focus();
}

function newItem(id, name, quantity, dateAdded) {
    const itemId = document.createElement('input');
    itemId.type = 'hidden';
    itemId.classList.add('id');
    if (id) {
        itemId.value = id;
    }
    const nameInput = document.createElement('input');
    nameInput.type = 'text';
    if (name) {
        nameInput.value = name;
    }
    nameInput.addEventListener('focusout', finishEditItem);
    nameInput.addEventListener('focusin', startEditItem);
    const quantityInput = document.createElement('input');
    quantityInput.type = 'text';
    if (quantity) {
        quantityInput.value = quantity;
    }
    const item = document.createElement('article');
    item.classList.add('shopping-item');
    item.addEventListener('click', function(e) {
        if (e.target.tagName !== 'INPUT' && e.target.tagName !== 'BUTTON') {
            e.target.querySelector('.item input').focus();
        }
    });
    item.appendChild(itemId);
    item.appendChild(label('Item', 'item', nameInput));
    item.appendChild(label('Quantity', 'quantity', quantityInput));

    if (id) {
        appendDeleteButton(item, id);
    } else {
        const spacer = document.createElement('div')
        spacer.classList.add('action-spacer');
        item.appendChild(spacer);
    }

    if (dateAdded) {
        appendDate(item, dateAdded);
    }

    return item;
}

function appendDeleteButton(article, itemId) {
    if (article.querySelector('.delete-item')) {
        return;
    }

    const spacer = article.querySelector('.action-spacer');
    if (spacer) {
        spacer.remove();
    }
    const deleteBtn = document.createElement('button');
    deleteBtn.addEventListener('click', () => deleteItem(itemId), false);
    deleteBtn.innerHTML = 'X';
    deleteBtn.classList.add('item-action', 'delete-item');
    article.appendChild(deleteBtn);
}

function appendDate(article, dateAdded) {
    if (article.querySelector('.dateAdded')) {
        return;
    }

    const br = document.createElement('div');
    br.classList.add('flex-break');
    article.appendChild(br);
    const date = document.createElement('span');
    date.classList.add('dateAdded');
    date.innerHTML = `Added on: <time datetime="${moment(dateAdded).format()}">${moment(dateAdded).format('YYYY-MM-DD HH:mm')}</time>`;
    article.appendChild(date);
}

function deleteItem(id) {
    httpDelete(`/api/items/${id}`, () => {
        const item = document.querySelector(`input[type=hidden][value="${id}"]`);
        item.closest('article').remove();
    });
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
    if(e.target.value && Array.from(document.querySelectorAll('.shopping-item > .item > input')).pop().value) {
        document.getElementById('add-item').disabled = false;
        const article = e.target.closest('article');
        const item = {
            name: article.querySelector('.item > input').value,
            quantity: article.querySelector('.quantity > input').value,
            dateAdded: new Date().getTime()
        }
        idInput = article.querySelector('.id');
        
        if (idInput.value) {
            item.id = Number(idInput.value);
        }
        
        httpPut('/api/items', item, (r) => { 
            idInput.value = JSON.parse(r.target.responseText).id;
            appendDeleteButton(article, idInput.value);
            appendDate(article, item.dateAdded);
        });
    }
}

function generateId() {
    return Math.random().toString(36).slice(2);
}

function httpGet(path, callback) {
    const r = new XMLHttpRequest();
    r.open('GET', path, true);
    r.onload = callback;
    r.send();
}

function httpPut(path, body, callback) {
    const r = new XMLHttpRequest();
    r.open('PUT', path, true);
    r.onreadystatechange = function () {
        if (r.readyState != 4 || r.status != 200 || r.status != 201) return;
    };
    r.onload = callback;
    r.send(JSON.stringify(body));
}

function httpDelete(path, callback) {
    const r = new XMLHttpRequest();
    r.open('DELETE', path, true);
    r.onload = callback;
    r.send();
}