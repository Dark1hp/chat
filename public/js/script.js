var ws = new WebSocket("ws://" + location.host + "/ws");

ws.onopen = function() {
    console.log("Connected");
};
ws.onclose = function(event) {
    if (event.wasClean) {
        console.log('Connection closed');
    } else {
        console.log('ERROR: Connection reset');
        console.log('Code: ' + event.code + 'reason: ' + event.reason);
    }
};
ws.onerror = function(error) {
    console.log("Error: " + error.message);
};

ws.onmessage = function(event) {
    var msg = JSON.parse(event.data);

    switch (msg.Type) {
        case "Text":
            textMsg(msg);
            break;
        case "Image":
            imgMsg(msg);
            break;
    }
}

function readUrl(input) {

    if (input.files && input.files[0]) {
        var reader = new FileReader();

        console.log(input.files[0]);

        reader.onload = function(e) {
/*            var blockView = document.querySelector('.js-view');
            var blockMsg = document.createElement('div');
            var blockCont = document.createElement('div');
            var blockText = document.createElement('p');
            var blockImg = document.createElement('img');
            blockImg.setAttribute('src', e.target.result);
            blockImg.setAttribute('width', '250');
            blockMsg.className = 'js-view__msg';
            blockCont.className = 'js-view__cont'
            blockText.className = 'js-view__text app-chat-history__text--img';
            blockText.appendChild(blockImg);
            blockCont.appendChild(blockText);
            blockMsg.appendChild(blockCont);
            blockView.insertBefore(blockMsg, blockView.firstChild);*/
            return "Ok";
        }
        sendMessage(input.files[0]);

        reader.readAsDataURL(input.files[0]);
    }
}

function sendMessage(msg) {
    ws.send(msg);
}

function handleSubmit() {
    var msgArea = document.querySelector(".js-msg-area");
    console.log(msgArea.value)
    sendMessage(msgArea.value);
    msgArea.value = '';
    return false;
}

function textMsg(msg) {
    console.log(msg);
    var blockView = document.querySelector('.js-view');
    var blockMsg = document.createElement('div');
    var blockCont = document.createElement('div');
    var blockText = document.createElement('p');
    var blockId = document.createElement('span');
    blockMsg.className = 'app-chat-history__msg';
    blockCont.className = 'app-chat-history__cont'
    blockText.className = 'app-chat-history__text';
    blockId.className = 'app-chat-history__name';
    blockText.innerHTML = msg.Msg;
    blockId.innerHTML = msg.Id;
    blockCont.appendChild(blockId);
    blockCont.appendChild(blockText);
    blockMsg.innerHTML = '<a class="app-chat-history__img"><img src="img/aomine.jpg" alt=""></a>';
    blockMsg.appendChild(blockCont);
    blockView.insertBefore(blockMsg, blockView.firstChild);
    return "Ok";
}

function imgMsg(msg) {
    console.log(msg);
    var blockView = document.querySelector('.js-view');
    var blockMsg = document.createElement('div');
    var blockCont = document.createElement('div');
    var blockText = document.createElement('p');
    var blockImg = document.createElement('img');
    var blockId = document.createElement('span');
    blockId.innerHTML = msg.Id;
    blockImg.setAttribute('src', msg.Msg);
    blockImg.setAttribute('width', '250');
    blockMsg.className = 'app-chat-history__msg';
    blockCont.className = 'app-chat-history__cont'
    blockText.className = 'app-chat-history__text app-chat-history__text--img';
    blockId.className = 'app-chat-history__name';
    blockText.appendChild(blockImg);
    blockCont.appendChild(blockId);
    blockCont.appendChild(blockText);
    blockMsg.innerHTML = '<a class="app-chat-history__img"><img src="img/aomine.jpg" alt=""></a>';
    blockMsg.appendChild(blockCont);
    blockView.insertBefore(blockMsg, blockView.firstChild);
    return "Ok";
}
