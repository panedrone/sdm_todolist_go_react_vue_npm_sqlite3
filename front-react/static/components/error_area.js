"use client";

import * as React from "react";

export const ErrorArea = ({saveUpdater}) => {

    const [value, setValue] = React.useState("")

    if (saveUpdater) {
        saveUpdater(setValue)
    }

    return (
        <div>
            {
                value.length > 0
                &&
                <div>
                    <button onClick={() => setValue("")}>&#x2713;</button>
                    &nbsp;
                    <strong>Error:</strong>&nbsp;{value}
                </div>}
        </div>
    );
}
