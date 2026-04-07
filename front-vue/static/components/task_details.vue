<script>
import fire from "./event_bus";
import {delete204, getJson, putJson, unicodeToChar} from "./api";

const NO_TASK = {"t_id": -1, "p_id": -1, "t_date": null, "t_subject": null, "t_priority": -1, "t_comments": null}

export default {
  data() {
    return {
      current_subject: null,
      current_task: NO_TASK,
      task_error: null,
      show_task_details: false,
    }
  },
  created() {
    fire.renderTaskDetails = this.renderTaskDetails
    fire.hideTaskDetails = this.hideTaskDetails
  },
  methods: {
    hideTaskDetails() {
      this.show_task_details = false
    },
    renderTaskDetails(t_id) {
      getJson("api/tasks/" + t_id, (json) => {
        this.current_subject = json.t_subject;
        this.current_task = json;
        this.task_error = null;
        this.show_task_details = true
      })
    },
    taskUpdate() {
      if (!isNaN(this.current_task.t_priority)) {
        this.current_task.t_priority = parseInt(this.current_task.t_priority);
      }
      let json = JSON.stringify(this.current_task)
      let p_id = this.current_task.p_id
      let t_id = this.current_task.t_id
      putJson("api/tasks/" + t_id, json, async (resp) => {
        if (resp.status === 200) {
          fire.renderProjectTasks(p_id);
          this.renderTaskDetails(t_id);
          return
        }
        await this._showTaskError(resp);
      })
    },
    async _showTaskError(resp) {
      let msg = await resp.text()
      msg = unicodeToChar(msg);
      // https://stackoverflow.com/questions/6640382/how-to-remove-backslash-escaping-from-a-javascript-var
      msg = msg.replace(/\\\\"/g, '"');
      msg = msg.replace(/\\"/g, '"');
      msg = resp.status.toString() + " ==> " + msg
      this.task_error = msg
    },
    taskDelete() {
      let p_id = this.current_task.p_id
      let t_id = this.current_task.t_id
      delete204("api/tasks/" + t_id, () => {
        this.hideTaskDetails()
        fire.renderProjects(); // update tasks count
        fire.renderProjectTasks(p_id);
      })
    },
    hideTaskError() {
      this.task_error = null
    }
  }
}
</script>

<template>
  <div class="card" v-if="show_task_details">
    <div class="title" id="subj">
      {{ current_subject }}
    </div>
    <table>
      <tr>
        <td class="form-label">Date</td>
        <td>
          <label>
            <input type="text" v-model="current_task.t_date"/>
          </label>
        </td>
      </tr>
      <tr>
        <td class="form-label">Subject</td>
        <td class="w100">
          <label>
            <input type="text" v-model="current_task.t_subject"/>
          </label>
        </td>
      </tr>
      <tr>
        <td class="form-label">Priority</td>
        <td>
          <label>
            <input type="text" v-model="current_task.t_priority"/>
          </label>
        </td>
      </tr>
      <tr>
        <td colspan="2">
          <label>
            <textarea cols="40" rows="10" v-model="current_task.t_comments"></textarea>
          </label>
        </td>
      </tr>
      <tr>
        <td>
          <table class="controls">
            <tr>
              <td class="w1">
                <a href="#" @click="taskUpdate()">
                  <input type="button" value="✓"/>
                </a>
              </td>
              <td class="w1">
                <a href="#" @click="taskDelete()">
                  <input type="button" value="x"/>
                </a>
              </td>
              <td class="w1">
                <a href="#" @click="hideTaskDetails()">
                  <input type="button" value="&lt;"/>
                </a>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
    <div style="color: crimson;max-width: 310px;"  v-if="task_error !== null">
      <div>
        <button @click="hideTaskError()">&#x2713;</button>
        &nbsp;
        <strong>Error:</strong>&nbsp;{{ task_error }}
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>