'use strict';

window.startSession = () => {
  console.log("hello Albin");

  let ws = new WebSocket("ws://lenovo.local:5000/ws", "jsonrpc-1.0");

  ws.onopen = (event) => {
    console.log("open" + event);

    let rpc = {method: "Arith.Multiply", params: [{A: 1, B: 2}], id: 5};
    console.log(rpc);
    ws.send(JSON.stringify(rpc));

    setTimeout(() => {
      console.log("closing");
      ws.close();
    }, 500);
  };

  ws.onmessage = (event) => {
    console.log("data:" + event.data);
  }
}
