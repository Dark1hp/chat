window.socket = new WebSocket("ws://" + location.host + "/ws");

/*document.getElementById("clip-file").onchange(function(event) {
  // console.log(this.files[0]);
  readUrl(this);
});*/

function readUrl(input) {

  if (input.files && input.files[0]) {
    var reader = new FileReader();

    reader.onload = function (e) {
     var blockView = document.querySelector('.js-view');
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
     blockView.insertBefore(blockMsg, blockView.firstChild);
     return "Ok";
   }
     sendMessage(input.files[0]);

   reader.readAsDataURL(input.files[0]);
 }
}

function sendMessage(msg) {
  socket.send(msg);
}

function handleSubmit() {
  var msgArea = document.querySelector(".js-msg-area");
  sendMessage(msgArea.value);
  msgArea.value = '';
  return false;
}

function setUpSocket(msg) {
  socket.onopen = function() {
    console.log("Connected");
  };
  socket.onclose = function(event) {
    if (event.wasClean) {
      console.log('Connection closed');
    } else {
      console.log('ERROR: Connection reset');
      console.log('Code: ' + event.code + 'reason: ' + event.reason);
    }
  };
  socket.onmessage = msg;
  socket.onerror = function(error) {
    console.log("Error: " + error.message);
  };
}
function displayMessage(msg) {
  var blockView = document.querySelector('.js-view');
  var blockMsg = document.createElement('div');
  var blockCont = document.createElement('div');
  var blockText = document.createElement('p');
  blockMsg.className = 'js-view__msg';
  blockCont.className = 'js-view__cont'
  blockText.className = 'js-view__text';
  blockText.innerHTML = msg.data;
  blockCont.appendChild(blockText);
  blockMsg.appendChild(blockCont);
  blockView.insertBefore(blockMsg, blockView.firstChild);
  return "Ok";
}
setUpSocket(displayMessage);