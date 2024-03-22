<script setup>

import {login} from "../service/auth";
import {sendMessage} from "../service/chat";
import {ref} from "vue";

const loggedIn = ref(false);
const phoneNumber = ref('');
const messageBody = ref('');
const messages = ref([]);

function doLogin(phoneNumber) {
  login(phoneNumber)
  loggedIn.value = true;
}

function doSendMessage() {
  sendMessage(messageBody);
  messageBody.value = '';
}

</script>

<template>
  <h1>Chat Server</h1>
  <div>
    <div v-if="!loggedIn">
      <h2>登录</h2>
      <input type="text" v-model="phoneNumber" placeholder="用户名">
      <button @click="doLogin(phoneNumber)">登录</button>
    </div>
    <div v-else>
      <h2>聊天室</h2>
      <div class="chat-box">
        <div v-for="(message, index) in messages" :key="index">
          {{ message }}
        </div>
      </div>
      <input type="text" v-model="messageBody" @keyup.enter="doSendMessage" placeholder="输入消息">
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
