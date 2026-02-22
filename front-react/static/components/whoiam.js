"use client";

import * as React from "react";
import * as api from "./api"
import {RawHtml} from "./raw_html"


export const WhoIAm = () => {

    const [who, setWho] = React.useState('')

    React.useEffect(() => {
        fetchWhoIAm()
    }, [])
    
    function fetchWhoIAm() {
        api.getText('api/whoiam', (text) => {
            if (!text) {
                return
            }
            if (text.includes('sqlx')) {
                text += ', <a target="_blank" href="swagger/index.html">swagger</a>'
            }
            text += ", npm, react " + React.version

            setWho(text)
        })
    }

    return (
        <RawHtml rawHtml={who}/>
    )
}