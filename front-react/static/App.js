"use client";

import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-calendar/dist/Calendar.css';
import 'react-tabs/style/react-tabs.css';

import 'bootstrap-icons/font/bootstrap-icons.css';

import * as React from "react";
import ReactDOM from "react-dom/client";
import {PaneProjectList} from './components/project_list'
import {WhoIAm} from "./components/whoiam";

import * as api from "./components/api";
import {PaneProjectDetails} from "./components/project_details";
import {TaskDetails} from "./components/task_details";
import {Col, Row} from "react-bootstrap";
import {ErrorArea} from "./components/error_area";

const App = () => {

    return <>
        <div className="d-flex">
            <WhoIAm/>
        </div>

        <div className="red-banner" id="serverError">
            <ErrorArea saveUpdater={(updater) => api.assignServerErrorUpdater(updater)} />
        </div>

        <Row className="g-3 w-auto h-auto">
            <Col className={"col-auto"}>
                <PaneProjectList/>
            </Col>
            <Col className={"col-auto"} id="projectDetails" style={{display: 'none'}}>
                <PaneProjectDetails/>
            </Col>
            <Col className={"col-auto"} id="taskForm" style={{display: 'none'}}>
                <TaskDetails/>
            </Col>
        </Row>
    </>
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    // <React.StrictMode>
    <App/>
    // </React.StrictMode>
);
