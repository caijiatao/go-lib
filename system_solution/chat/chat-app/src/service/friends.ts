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
