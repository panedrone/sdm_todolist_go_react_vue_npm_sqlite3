"use client";

import * as React from "react";
import {StringField} from "./form_components";
import fire from './event_bus.js'
import * as api from "./api"
import {Button, Card, Table} from "react-bootstrap";

let _updateProjects = (_) => {
}

const ProjectList = () => {

    function _handleClick(project) {
        fire.fetchCurrentProject(project.p_id)
        fire.fetchProjectTasks(project.p_id)
        fire.setVisibleTaskForm(false)
    }

    React.useEffect(() => {
        fetchProjects()
    }, [])

    const [data, setData] = React.useState([])

    _updateProjects = setData

    return <Table className="section mb-3">
        <thead>
        <tr>
            <th>Project</th>
            <th className="text-nowrap" width={1}>Tasks Count</th>
        </tr>
        </thead>
        <tbody>
        {
            data.map((project, index) => {
                    return (
                        <tr key={index}>
                            <td onClick={() => _handleClick(project)}>
                                <a>{project.p_name}</a>
                            </td>
                            <td className="center w1">
                                {project.p_tasks_count}
                            </td>
                        </tr>
                    )
                }
            )
        }
        </tbody>
    </Table>
}

export function fetchProjects() {
    api.getJsonArray("api/projects", (arr) => {
        if (arr) {
            _updateProjects(arr)
        }
    })
}

const ProjectCreateButton = () => {

    function projectCreate() {
        if (_newProjectName.length === 0) {
            _newProjectName = '?'
        }
        let json = JSON.stringify({"p_name": _newProjectName})
        api.postJson201("api/projects", json, () => {
            fire.fetchProjects();
        })
    }

    return (
        <Button variant={"secondary"} className={"opacity-50"} onClick={() => projectCreate()}>
            +
        </Button>
    )
}

let _newProjectName = ""

const ProjectCreateField = () => {
    return (
        <StringField onChange={v => _newProjectName = v}/>
    )
}

export const PaneProjectList = () => {
    return (
        <Card className={"bg-white rounded-3 mb-3"}>
            <Card.Body>
                <ProjectList/>
                <div className="d-inline-flex align-items-center gap-2 w-100">
                    <ProjectCreateField/>
                    <ProjectCreateButton/>
                </div>
           </Card.Body>
        </Card>
    )
}