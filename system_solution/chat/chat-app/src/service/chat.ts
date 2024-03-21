import {ref} from "vue";
import {authService} from "../util/axios";

let websocket = null;
const messages = ref([]);

export const initWebSocket = () => {
    websocket = new WebSocket('ws://localhost:13177/ws/chat?user_id=' + authService.getUserId());
    websocket.onmessage = handleMessage;
};

const handleMessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('Received message:', message);
    messages.value.push(`${message.from_user}: ${message.message_body}`);
};

export const sendMessage = (messageBody: any) => {
    if (messageBody.value.trim() !== '') {
        const message = {
            // to_user: 111, TODO 选择发送对象
            content: messageBody.value,
        };
        console.log(message)
        websocket.send(JSON.stringify(message));
        messageBody.value = '';
    }
};