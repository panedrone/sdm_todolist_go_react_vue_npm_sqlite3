"use client";

import * as React from "react";
import {StringField} from "./form_components";
import fire from './event_bus'
import * as api from "./api"
import {Button, Card} from "react-bootstrap";

let _updateCurrentProjectName = (_) => {
}

const ProjectButtons = () => {

    function _projectUpdate() {
        let p_id = _currentProject.p_id
        let json = JSON.stringify(vars.currentProject)
        api.putJson200("api/projects/" + p_id, json, () => {
            fire.fetchProjects();
        })
    }

    function _projectDelete() {
        let p_id = _currentProject.p_id
        api.delete204(`api/projects/${p_id}`, () => {
            fire.setVisibleProjectDetails(false)
            fire.fetchProjects();
        })
    }

    return (
        <table className="controls">
            <tbody>
            <tr>
                <td id="currentProjectName">
                    <StringField onChange={v => {
                        _currentProject.p_name = v
                    }} saveUpdater={(updater) => {
                        _updateCurrentProjectName = updater
                    }}/>
                </td>
                <td className="w1">
                    <Button variant={"secondary"} className={"opacity-50"}
                            onClick={() => _projectUpdate()}
                    >
                        &#x2713;
                    </Button>
                </td>
                <td className="w1">
                    <Button variant={"secondary"} className={"opacity-50"}
                            onClick={() => _projectDelete()}
                    >
                        x
                    </Button>
                </td>
                <td className="w1">
                    <Button variant={"secondary"} className={"opacity-50"}
                            onClick={() => fire.setVisibleProjectDetails(false)}
                    >
                        &lt;
                    </Button>
                </td>
            </tr>
            </tbody>
        </table>
    )
}

let _updateProjectTasks = (_) => {
}

const ProjectTasks = () => {

    const [data, setData] = React.useState([])

    _updateProjectTasks = setData

    return <table className="section">
        <thead>
        <tr>
            <th className="w1">Date</th>
            <th>Subject</th>
            <th className="w1">Priority</th>
        </tr>
        </thead>
        <tbody>
        {
            data.map((task, index) => {
                    return (
                        <tr key={index}>
                            <td className="w1 nowrap">
                                {task.t_date}
                            </td>
                            <td onClick={() => fire.fetchTask(task.t_id)}>
                                <a>{task.t_subject}</a>
                            </td>
                            <td className="center">
                                {task.t_priority}
                            </td>
                        </tr>
                    )
                }
            )
        }
        </tbody>
    </table>
}

const TaskCreateButton = () => {

    function _taskCreate() {
        let p_id = _currentProject.p_id
        let json = JSON.stringify({"t_subject": _newTaskSubject})
        api.postJson201(`api/projects/${p_id}/tasks`, json, () => {
            fire.fetchProjects();
            fire.fetchProjectTasks(p_id); // update tasks count
        })
    }

    return (
        <Button variant={"secondary"} className={"opacity-50"} onClick={() => _taskCreate()}>
            +
        </Button>
    )
}

let _currentProject = {p_id: -1, p_name: ""};

fire.fetchCurrentProject = (p_id) => {
    api.getJson(`api/projects/${p_id}`, (json) => {
        if (!json) {
            fire.showServerError("failed to get project data")
            return
        }
        _currentProject = json
        _updateCurrentProjectName(_currentProject.p_name)
    })
}

fire.fetchProjectTasks = (p_id) => {
    api.getJsonArray(`api/projects/${p_id}/tasks`, (arr) => {
        fire.setVisibleProjectDetails(true)
        if (arr) {
            _updateProjectTasks(arr)
        }
    })
}

let _newTaskSubject = ""

export const PaneProjectDetails = () => {

    React.useEffect(() => {

    }, [])

    return <>
        <Card className={"bg-white rounded-3 mb-3"}>
            <Card.Body>

                <table>
                    <tbody>
                    <tr>
                        <td id="projectActions">
                            <ProjectButtons/>
                        </td>
                    </tr>
                    <tr>
                        <td id="tasks">
                            <ProjectTasks/>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <table className="controls">
                                <tbody>
                                <tr>
                                    <td id="newTaskSubject">
                                        <StringField onChange={v =>
                                            _newTaskSubject = v
                                        }/>
                                    </td>
                                    <td className="w1" id="taskCreate">
                                        <TaskCreateButton/>
                                    </td>
                                </tr>
                                </tbody>
                            </table>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </Card.Body>
        </Card>

    </>
}
