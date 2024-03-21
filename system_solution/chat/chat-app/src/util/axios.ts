import axios, {AxiosInstance} from 'axios';

class AuthService {
    private token: string | null = null;
    private userId: Number | null = null;

    // 设置登录后的 Token
    setToken(token: string) {
        this.token = token;
    }

    // 获取 Token
    getToken(): string | null {
        return this.token;
    }

    setUserId(userId: Number) {
        this.userId = userId;
    }

    getUserId(): Number | null {
        return this.userId;
    }

    // 检查是否已登录
    isLoggedIn(): boolean {
        return !!this.token;
    }
}

export let authService = new AuthService();


// 创建 Axios 实例
const axiosInstance: AxiosInstance = axios.create({
    baseURL: 'http://localhost:8080', // 替换为您的 API 地址
});

// 请求拦截器
axiosInstance.interceptors.request.use((config) => {
    // 如果已登录，则在请求头中添加 Token
    if (authService.isLoggedIn()) {
        config.headers.Authorization = `Bearer ${authService.getToken()}`;
    }
    return config;
}, (error) => {
});


