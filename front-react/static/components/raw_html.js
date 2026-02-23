"use client";

import * as React from "react";

export const RawHtml = ({rawHtml}) => {
    const __html = html => ({
        // === panedrone: wrap it just to get rid of ide warnings
        __html: html, // https://stackoverflow.com/questions/29044518/safe-alternative-to-dangerouslysetinnerhtml
    });
    return (
        <span dangerouslySetInnerHTML={__html(rawHtml)}></span>
    )
}
