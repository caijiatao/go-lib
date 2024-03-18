<script setup>
import {ref} from "vue";
import axios from "axios";

const loggedIn = ref(false);
const phoneNumber = ref('');
const messageBody = ref('');
const messages = ref([]);
let userId = ref('');
let websocket = null;

const login = async () => {
  if (phoneNumber.value.trim() !== '') {
    try {
      const response = await axios.post('http://localhost:8080/api/user/login', { phoneNumber: phoneNumber.value });
      if (response.status === 200) {
        loggedIn.value = true;
        userId = response.data.data.id;
        initWebSocket();
      } else {
        alert('登录失败');
      }
    } catch (error) {
      console.error('登录请求失败:', error);
      alert('登录请求失败');
    }
  } else {
    alert('请输入用户名');
  }
};

const initWebSocket = () => {
  websocket = new WebSocket('ws://localhost:8080/ws/chat?user_id=' + userId);
  websocket.onmessage = handleMessage;
};

const handleMessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received message:', message);
  messages.value.push(`${message.from_user}: ${message.message_body}`);
};

const sendMessage = () => {
  if (messageBody.value.trim() !== '') {
    const message = {
      // to_user: 111, TODO 选择发送对象
      message_body: messageBody.value,
    };
    websocket.send(JSON.stringify(message));
    messageBody.value = '';
  }
};

</script>

<template>
  <h1>Chat Server</h1>
  <div>
    <div v-if="!loggedIn">
      <h2>登录</h2>
      <input type="text" v-model="phoneNumber" placeholder="用户名">
      <button @click="login">登录</button>
    </div>
    <div v-else>
      <h2>聊天室</h2>
      <div class="chat-box">
        <div v-for="(message, index) in messages" :key="index">
          {{ message }}
        </div>
      </div>
      <input type="text" v-model="messageBody" @keyup.enter="sendMessage" placeholder="输入消息">
    </div>
  </div>
</template>

<style scoped>
.chat-box {
  height: 300px;
  overflow-y: scroll;
  border: 1px solid #ccc;
  padding: 10px;
}
</style>