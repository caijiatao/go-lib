import {ref} from "vue";
import {authService, axiosInstance} from "../util/axios";

let websocket = null;
const messages = ref([]);

class Message {
    to_user: Number;
    content: string;

    constructor(to_user: Number, content: string) {
        this.to_user = to_user;
        this.content = content;
    }
}

export const initWebSocket = () => {
    websocket = new WebSocket('ws://localhost:13177/ws/chat?user_id=' + authService.getUserId());
    websocket.onmessage = handleMessage;
    // TODO 定时发送心跳
    setInterval(() => {
        // websocket.send('ping');
    }, 10000);
};

const handleMessage = (event) => {
    const message = JSON.parse(event.data);
    messages.value.push(`${message.from_user}: ${message.message_body}`);
};

export function sendMessage(messageBody: any) {
    if (messageBody.value.trim() !== '') {
        const message = {
            to_user: 111,
            content: messageBody.value,
        };
        console.log(message)
        websocket.send(JSON.stringify(message));
        messageBody.value = '';
    }
};

