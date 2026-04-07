import {createApp} from 'vue';
import whoiam from "./components/whoiam.vue";
import server_error from "./components/server_error.vue";
import projects from "./components/project_list.vue";
import project_details from "./components/project_details.vue";
import task_details from "./components/task_details.vue";
import fire from "./components/event_bus";

const app = createApp({
    components: {
        whoiam,
        server_error,
        projects,
        project_details,
        task_details
    },
    mounted() {
        fire.renderWhoIAm();
        fire.renderProjects();
    },
})

app.mount("#app")
