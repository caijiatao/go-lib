// api.js

import axios from 'axios';

// 创建一个 Axios 实例
const axiosInstance = axios.create({
    baseURL: 'http://localhost:8080',
});

// 登录方法，接收用户名和密码，返回登录成功后的 Token
export const login = async (phone_number) => {
    try {
        const response = await axiosInstance.post('/api/user/login', { phone_number: phone_number });
        if (response.status === 200) {
            return response.data.token;
        } else {
            throw new Error('登录失败');
        }
    } catch (error) {
        throw new Error('登录请求失败');
    }
};

// 创建一个带有 Token 的 Axios 实例
export const createAuthenticatedInstance = (token) => {
    const authenticatedInstance = axios.create({
        baseURL: 'http://localhost:8080',
        headers: {
            Authorization: `Bearer ${token}`, // 在请求头中添加 Token
        },
    });
    return authenticatedInstance;
};
