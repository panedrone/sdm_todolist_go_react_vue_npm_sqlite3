"use client";

import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-calendar/dist/Calendar.css';
import 'react-tabs/style/react-tabs.css';

import 'bootstrap-icons/font/bootstrap-icons.css';

import * as React from "react";
import ReactDOM from "react-dom/client";
import {ProjectCreateButton, ProjectCreateField, ProjectList} from './components/project_list'
import {WhoIAm} from "./components/whoiam";


// project_list.renderComponents()
// project_details.renderComponents()
// task_details.renderComponents()

// api.renderComponents()

const App = () => {

    return <>
        <table className="bg">
            <tbody>
            <tr>
                <td>
                    <div className="card">
                        <h2 id="whoiam">
                            <WhoIAm/>
                        </h2>
                    </div>
                </td>
            </tr>
            <tr>
                <td>
                    <div className="red-banner" id="serverError"></div>
                </td>
            </tr>
            </tbody>
        </table>

        <table className="bg">
            <tbody>
            <tr>
                <td>
                    <div className="card">
                        <table>
                            <tbody>
                            <tr>
                                <td id="projects">
                                    <ProjectList/>
                                </td>
                            </tr>
                            <tr>
                                <td>
                                    <table className="controls">
                                        <tbody>
                                        <tr>
                                            <td id="newProjectName">
                                                <ProjectCreateField/>
                                            </td>
                                            <td className="w1" id="projectCreate">
                                                <ProjectCreateButton/>
                                            </td>
                                        </tr>
                                        </tbody>
                                    </table>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </td>
                <td id="projectDetails" style={{display: 'none'}}>
                    <div className="card">
                        <table>
                            <tbody>
                            <tr>
                                <td id="projectActions">
                                </td>
                            </tr>
                            <tr>
                                <td id="tasks">
                                </td>
                            </tr>
                            <tr>
                                <td>
                                    <table className="controls">
                                        <tbody>
                                        <tr>
                                            <td id="newTaskSubject">
                                            </td>
                                            <td className="w1" id="taskCreate">
                                            </td>
                                        </tr>
                                        </tbody>
                                    </table>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </td>
                <td id="taskForm" style={{padding: '2px', display: 'none'}}>
                    <div className="card">
                        <div id="taskActions"></div>

                        <table className="edit-form">
                            <tbody>
                            <tr>
                                <td className="form-label">Date</td>
                                <td id="t_date">
                                </td>
                            </tr>
                            <tr>
                                <td className="form-label">Subject</td>
                                <td className="w100" id="t_subject">
                                </td>
                            </tr>
                            <tr>
                                <td className="form-label w1">Priority 1..10</td>
                                <td id="t_priority">
                                </td>
                            </tr>
                            <tr>
                                <td colSpan="2" id="t_comments">
                                </td>
                            </tr>
                            </tbody>
                        </table>
                        <div className="task-error" id="taskError">
                        </div>
                    </div>
                </td>
            </tr>
            </tbody>
        </table>

    </>
}

// async function windowOnLoad() {
//     whoiam.fetchWhoIAm()
//     project_list.fetchProjects()
// }
//
// windowOnLoad().then(() => console.log('== windowOnLoad() completed =='))

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    // <React.StrictMode>
    <App/>
    // </React.StrictMode>
);
