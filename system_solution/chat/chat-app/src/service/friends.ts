import {ref} from "vue";
import {axiosInstance} from "../util/axios";

const friends = ref([]);

export function getFriends() {
    axiosInstance.get('/api/friends').then((response) => {
        console.log(response.data);
    }).catch((error) => {
        console.log(error);
    });
}

export function getFriendMessages(friendId: number) {
    axiosInstance.get(`/api/friends/${friendId}/messages`).then((response) => {
        console.log(response.data);
    }).catch((error) => {
        console.log(error);
    });
}
