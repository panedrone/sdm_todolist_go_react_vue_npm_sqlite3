"use client";

import * as React from "react";
import * as api from "./api"
import {RawHtml} from "./raw_html"
import {Card} from "react-bootstrap";


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
        <Card className={"bg-white rounded-3 mb-2"}>
            <Card.Body>
                <h2 id="whoiam">
                    <RawHtml rawHtml={who}/>
                </h2>
            </Card.Body>
        </Card>
    )
}