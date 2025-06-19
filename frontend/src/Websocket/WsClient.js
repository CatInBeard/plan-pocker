class WsClient {

    constructor(setConncetionStatus, pushMessage) {
        if (WsClient.instance) {
            return WsClient.instance;
        }
        WsClient.instance = this;
        this.setConncetionStatus = setConncetionStatus
        this.pushMessage = pushMessage
        this.socketConnect()
    }

    send(message){
        if (this.socket.readyState > 1) {
            console.log("Connection closed, try to reconnect")
            this.socketConnect()
        }
        this.socket.send(JSON.stringify(message))
    }

    socketConnect() {

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.host;
        const websocketUrl = `${protocol}//${host}/api/websocket`;

        this.socket = new WebSocket(websocketUrl)

        this.socket.addEventListener("open", event => {
            this.setConncetionStatus("established")
        });

        this.socket.addEventListener("message", event => {
            try {
            var message = JSON.parse(event.data)
            }  catch (error){
                console.error("Error :" + error + "\nmessage:" + event.data)
                return
            }
            this.pushMessage(message)
        });
        this.socket.addEventListener("close", (event) => {
            this.setConncetionStatus("closed")
            setTimeout(() => {
                this.socketConnect()
                console.log("Try to reconnect")
            }, 1000)
        });

        this.socket.onerror = function(event) {
            console.error("Ws error:", event);
        };
    }

}

export default WsClient;