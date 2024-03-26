<script setup>

import {login} from "../service/auth";
import {sendMessage} from "../service/chat";
import {ref} from "vue";


const loggedIn = ref(true);
const phoneNumber = ref('');
const messageBody = ref('');
let messages = ref([]);
const friends = [
  {id: 1, name: '张三'},
  {id: 2, name: '李四'},
  {id: 3, name: '王五'}
]
let selectedFriend = null


function doLogin(phoneNumber) {
  login(phoneNumber)
  loggedIn.value = true;
}

function doSendMessage() {
  sendMessage(messageBody);
  messageBody.value = '';
}

function handleSelectFriend(index) {
  const friendId = parseInt(index);
  selectedFriend = friends.find(friend => friend.id === friendId);
  // TODO: 根据选择的好友加载聊天消息
  loadMessagesForFriend(friendId);
}

function loadMessagesForFriend(friendId) {
}

</script>

<template>
  <h1>Chat Server</h1>

  <div v-if="!loggedIn">
    <h2>登录</h2>
    <input type="text" v-model="phoneNumber" placeholder="用户名">
    <button @click="doLogin(phoneNumber)">登录</button>
  </div>
  <div v-else>
    <div class="chat-container">
      <!-- 好友列表部分 -->
      <div class="friend-list">
        <el-menu
            default-active="1"
            class="el-menu-vertical-demo"
            @select="handleSelectFriend"
        >
          <el-menu-item-group title="好友列表">
            <el-menu-item index="1" v-for="friend in friends" :key="friend.id">
              {{ friend.name }}
            </el-menu-item>
          </el-menu-item-group>
        </el-menu>
      </div>

      <!-- 聊天框部分 -->
      <div class="chat-box">
        <el-card v-if="selectedFriend" class="box-card">
          <div slot="header" class="clearfix">
            <span>与 {{ selectedFriend.name }} 的聊天</span>
          </div>
          <div>
            <!-- 这里可以放置聊天内容 -->
            <p v-for="message in messages" :key="message.id">
              {{ message.content }}
            </p>
          </div>
        </el-card>

        <!-- 欢迎页部分 -->
        <el-card v-else class="box-card">
          <div slot="header" class="clearfix">
            <span>欢迎来到聊天室</span>
          </div>
          <div>
            <p>请选择一个好友开始聊天</p>
          </div>
        </el-card>
      </div>
    </div>
  </div>
  <input type="text" v-model="messageBody" @keyup.enter="doSendMessage" placeholder="输入消息">
</template>

<style scoped>
.chat-box {
  height: 300px;
  overflow-y: scroll;
  border: 1px solid #ccc;
  padding: 10px;
}

.chat-container {
  display: flex;
}

.friend-list {
  width: 30%;
}

.chat-box {
  flex-grow: 1;
  padding: 20px;
}

.box-card {
  margin-bottom: 20px;
}
</style>
