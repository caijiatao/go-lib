import axios from 'axios';
import {authService} from "../util/axios";
import {initWebSocket} from "./chat";


// 创建一个 Axios 实例
const axiosInstance = axios.create({
    baseURL: 'http://localhost:13177',
});

// 登录方法，接收用户名和密码，返回登录成功后的 Token
export const login = async (phoneNumber: string) => {
    try {
        const response = await axiosInstance.post('/api/user/login', {phoneNumber});
        if (response.status === 200) {
            authService.setToken(response.data.token);
            authService.setUserId(response.data.data.id);
            initWebSocket();
            return true;
        } else {
            throw new Error('登录失败');
        }
    } catch (error) {
        console.log(error)
        throw new Error('登录请求失败');
    }
};

