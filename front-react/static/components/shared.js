// import * as ReactDOM from "react-dom";
import ReactDOM from "react-dom/client";

import fire from "./event_bus";

export function render(component, containerID) {
    // ReactDOM.render(component, document.getElementById(containerID))
    const root = ReactDOM.createRoot(document.getElementById(containerID));
    root.render(
        // <React.StrictMode>
        component
        // </React.StrictMode>
    );
}

fire.setVisibleProjectDetails = (yes) => {
    let el = document.getElementById('projectDetails')
    if (yes) {
        el.style.display = "table-cell"; // to show
    } else {
        el.style.display = "none"; // to hide
        fire.setVisibleTaskForm(false)
    }
}

fire.setVisibleTaskForm = (yes) => {
    let el = document.getElementById('taskForm')
    if (yes) {
        el.style.display = "table-cell"; // to show
        // el.style.display = "block"; // to show
    } else {
        el.style.display = "none"; // to hide
    }
}
