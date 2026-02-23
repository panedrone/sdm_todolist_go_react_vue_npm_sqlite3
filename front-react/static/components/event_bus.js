export default {

    fetchProjects: () => {
        console.log('fetchProjects')
    },

    fetchCurrentProject: (p_id) => {
        console.error('fetchCurrentProject', p_id)
    },

    fetchProjectTasks: (p_id) => {
        console.error('fetchProjectTasks', p_id)
    },

    fetchTask: (t_id) => {
        console.error('fetchTask', t_id)
    },

    setVisibleProjectDetails: (yes) => {
        console.error('setVisibleProjectDetails', yes)
    },

    setVisibleTaskForm: (yes) => {
        console.error('setVisibleTaskForm', yes)
    },

    showServerError: (msg) => {
        console.error('showServerError', msg)
    },
}
