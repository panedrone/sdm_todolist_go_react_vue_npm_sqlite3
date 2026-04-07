<script>
import fire from "./event_bus";
import {delete204, getJson, getJsonArray, postJson201, putJson200} from "./api";

const NO_PROJECT = {"p_id": -1, "p_name": null, "p_tasks_count": -1}
const TASK = {"t_id": -1, "t_date": null, "t_subject": null, "t_priority": -1}

export default {
  data() {
    return {
      current_project: NO_PROJECT,
      tasks: [TASK],
      t_subject: null,
      show_project_details: false,
    }
  },
  created() {
    fire.hideProjectDetails = this.hideProjectDetails;
    fire.renderProjectDetails = this.renderProjectDetails;
    fire.renderProjectTasks = this.renderProjectTasks;
  },
  methods: {
    hideProjectDetails() {
      fire.hideTaskDetails()
      this.show_project_details = false
    },
    renderProjectDetails(p_id) {
      this.renderCurrentProject(p_id)
      this.renderProjectTasks(p_id)
      this.show_project_details = true
    },
    renderCurrentProject(p_id) {
      getJson("api/projects/" + p_id, (json) => {
        this.current_project = json
      })
    },
    renderProjectTasks(p_id) {
      getJsonArray("api/projects/" + p_id + "/tasks", (arr) => {
        this.tasks = arr
      })
    },
    projectUpdate() {
      let p_id = this.current_project.p_id
      let json = JSON.stringify(this.current_project)
      putJson200("api/projects/" + p_id, json, () => {
        fire.renderProjects();
      })
    },
    projectDelete() {
      let p_id = this.current_project.p_id
      delete204("api/projects/" + p_id, () => {
        this.hideProjectDetails()
        fire.renderProjects();
      })
    },
    taskCreate() {
      let p_id = this.current_project.p_id
      let json = JSON.stringify({"t_subject": this.t_subject})
      postJson201("api/projects/" + p_id + "/tasks", json, () => {
        fire.renderProjects(); // update tasks count
        this.renderProjectDetails(p_id);
      })
    },
    renderTaskDetails(t_id) {
      // this.$root.renderTaskDetails(t_id);
      fire.renderTaskDetails(t_id)
    }
  }
}
</script>

<template>
  <div class="card" v-if="show_project_details">
  <table>
    <tr>
      <td>
        <table class="controls">
          <tr>
            <td>
              <label>
                <input type="text" class="project-header" v-model="current_project.p_name"/>
              </label>
            </td>
            <td class="w1">
              <a href="#" @click="projectUpdate()">
                <input type="button" value="✓"/>
              </a>
            </td>
            <td class="w1">
              <a href="#" @click="projectDelete()">
                <input type="button" value="x"/>
              </a>
            </td>
            <td class="w1">
              <a href="#" @click="hideProjectDetails()">
                <input type="button" value="&lt;"/>
              </a>
            </td>
          </tr>
        </table>
      </td>
    </tr>
    <tr>
      <td>
        <table class="section">
          <tr>
            <th>Date</th>
            <th>Subject</th>
            <th>Priority</th>
          </tr>
          <tr v-for="t in tasks">
            <td>
              {{ t.t_date }}
            </td>
            <td>
              <a href="#" @click="renderTaskDetails(t.t_id)">
                {{ t.t_subject }}
              </a>
            </td>
            <td style="text-align: center;">
              {{ t.t_priority }}
            </td>
          </tr>
        </table>
      </td>
    </tr>
    <tr>
      <td>
        <table class="controls">
          <tr>
            <td>
              <label>
                <input type="text" v-model="t_subject"/>
              </label>
            </td>
            <td class="w1">
              <a href="#" @click="taskCreate()">
                <input type="button" value="+"/>
              </a>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
  </div>
</template>

<style scoped>

</style>