"use client";

import * as React from "react";

import * as shared from "./shared";
import {ErrorArea} from "./error_area";
import {IntegerField, StringField, TextAreaField} from "./form_components";
import fire from './event_bus.js'
import * as api from "./api"
import {Button, Card} from "react-bootstrap";

let _updateSubject = (_) => {
}
let _updateDate = (_) => {
}
let _updatePriority = (_) => {
}
let _updateComments = (_) => {
}

let _currentTask = {
    "t_id": 0,
    "p_id": 0,
    "t_priority": 0,
    "t_date": "2024-02-12 03:34:16",
    "t_subject": "",
    "t_comments": ""
}

let _updateTaskTitle = (_) => {
}

const TaskTitle = ({initial}) => {
    const [taskTitle, setTaskTitle] = React.useState(initial);
    _updateTaskTitle = setTaskTitle
    return <span>{taskTitle}</span>
}

fire.fetchTask = (t_id) => {
    api.getJson(`api/tasks/${t_id}`, (json) => {
        fire.setVisibleTaskForm(true)
        _currentTask = json
        if (!_currentTask) {
            fire.showServerError("failed to get task data")
            return
        }
        _updateTaskTitle(_currentTask.t_subject)
        _updateSubject(_currentTask.t_subject)
        _updateDate(_currentTask.t_date)
        _updatePriority(_currentTask.t_priority)
        _updateComments(_currentTask.t_comments)
    })
}

const TaskButtons = () => {

    function _taskUpdate() {
        if (!isNaN(_currentTask.t_priority)) {
            _currentTask.t_priority = parseInt(_currentTask.t_priority.toString());
            if (!_currentTask.t_priority) {
                _currentTask.t_priority = 1
            }
        }
        let json = JSON.stringify(_currentTask)
        let t_id = _currentTask.t_id
        api.putJson(`api/tasks/${t_id}`, json, async (resp) => {
            if (resp.status === 200) {
                fire.fetchProjectTasks(_currentTask.p_id);
                fire.fetchTask(_currentTask.t_id);
                _hideTaskError()
                return
            }
            await _showTaskError(resp)
        })
    }

    function _taskDelete() {
        let t_id = _currentTask.t_id
        api.delete204(`api/tasks/${t_id}`, () => {
            fire.fetchProjects(); // update tasks count
            fire.fetchProjectTasks(_currentTask.p_id);
            fire.setVisibleTaskForm(false)
        })
    }

    return <>
        <table className="controls">
            <tbody>
            <tr>
                <td className="w100">
                    <div className="title" id="taskTitle">
                        <TaskTitle/>
                    </div>
                </td>
                <td className="w1">
                    <Button
                        variant={"secondary"} className={"opacity-50"}
                        onClick={() => _taskUpdate()}
                    >
                        &#x2713;
                    </Button>
                </td>
                <td className="w1">
                    <Button
                        variant={"secondary"} className={"opacity-50"}
                        onClick={() => _taskDelete()}
                    >
                        x
                    </Button>
                </td>
                <td className="w1">
                    <Button
                        variant={"secondary"} className={"opacity-50"}
                        onClick={() => fire.setVisibleTaskForm(false)}
                    >
                        &lt;
                    </Button>
                </td>
            </tr>
            </tbody>
        </table>
    </>
}

function _hideTaskError() {
    _updateTaskError("")
}

async function _showTaskError(resp) {
    let msg = await resp.text()
    msg = api.unicodeToChar(msg);
    // https://stackoverflow.com/questions/6640382/how-to-remove-backslash-escaping-from-a-javascript-var
    msg = msg.replace(/\\\\"/g, '"');
    msg = msg.replace(/\\"/g, '"');
    msg = resp.status.toString() + " ==> " + msg
    _updateTaskError(msg)
}

let _updateTaskError = (_) => {
}

export const TaskDetails = () => {
    React.useEffect(() => {
    }, [])
    return <>
        <Card className={"bg-white rounded-3 mb-3"}>
            <Card.Body>
                <div id="taskActions">
                    <TaskButtons/>
                </div>

                <table className="edit-form">
                    <tbody>
                    <tr>
                        <td className="form-label">Date</td>
                        <td id="t_date">
                            <StringField onChange={v => {
                                _currentTask.t_date = v
                            }} saveUpdater={(updater) => {
                                _updateDate = updater
                            }}/>
                        </td>
                    </tr>
                    <tr>
                        <td className="form-label">Subject</td>
                        <td className="w100" id="t_subject">
                            <StringField onChange={v => {
                                _currentTask.t_subject = v
                            }} saveUpdater={(updater) => {
                                _updateSubject = updater
                            }}/>
                        </td>
                    </tr>
                    <tr>
                        <td className="form-label w1">Priority 1..10</td>
                        <td id="t_priority">
                            <IntegerField onChange={v => {
                                _currentTask.t_priority = v
                            }} saveUpdater={(updater) => {
                                _updatePriority = updater
                            }}/>
                        </td>
                    </tr>
                    <tr>
                        <td colSpan="2" id="t_comments">
                            <TextAreaField onChange={v => {
                                _currentTask.t_comments = v
                            }} saveUpdater={(updater) => {
                                _updateComments = updater
                            }}/>
                        </td>
                    </tr>
                    </tbody>
                </table>
                <div className="task-error" id="taskError">
                    <ErrorArea saveUpdater={(updater) => _updateTaskError = updater}/>
                </div>

            </Card.Body>
        </Card>
    </>
}