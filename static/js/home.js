var message;
var ws;
window.onload = async function() {
  //update geolocation
  console.log("Loaded page");
  let links = document.querySelectorAll("a");
  let navbar = document.querySelector(".navbar-burger");
  let navmenu = document.querySelector(".navbar-menu");
  let HOST = "ws://" + document.location.host + "/ws/" + localStorage.token;
  ws = new WebSocket(HOST);
  function addEventListenerList(list, event, fn) {
    for (var i = 0, len = list.length; i < len; i++) {
      list[i].addEventListener(event, fn, false);
    }
  }
  navbar.onclick = function() {
    if (navbar.classList.contains("is-active")) {
      navbar.classList.remove("is-active");
      navmenu.classList.remove("is-active");
    } else {
      navbar.classList.add("is-active");
      navmenu.classList.add("is-active");
    }
  }
  ws.onopen = () => {
    ws.send("Connection opened from this side")
  }
  ws.onmessage = function(event) {
    console.log("Got here!");
    message = event.data;
    //res = JSON.parse(res);
    //if (res.status) {
    //  resolveMessage();
    //}
    var msg = JSON.parse(message);
    console.log(msg.column);
    executeComponent(msg.column);
    console.log("function closed");
  }
  function executeComponent(column) {
    if (column == "rightColumn") {
      resolveMessage();
    }
  }
  addEventListenerList(links, "click", function(target){
    let link = target.target.innerHTML.trim()
    let command = {
      "type": "command",
      "commandType": "link",
      "component": link
    }
    ws.send(JSON.stringify(command))
  })
  async function getIp() {
    response = await fetch('https://ipinfo.io/json')
    json = response;
    return json.json();
  }
  function getLocation() {
    if (!navigator.geolocation)
      return false;
    let promise = new Promise((resolve, reject) => {
      let position = navigator.geolocation.getCurrentPosition((pos) => {
        return resolve(pos);
      });
    });
    return promise;
    //if (navigator.geolocation) {
      //let position = navigator.geolocation.getCurrentPosition((pos) => {
        //return pos;
      //})
      //result = wait position;
      //return result;
    //}
  }
  //get location from navigator.geolocation api
  pos = await getLocation();
  pos = pos["coords"];
  pos = [pos.longitude, pos.latitude];
  //get location from ip address
  json = await getIp();
  json = json["loc"].split(",").reverse();
  json[0] = Number(json[0]);
  json[1] = Number(json[1]);
  console.log(json);
  console.log(pos);
}
