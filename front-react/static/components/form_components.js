"use client";

import * as React from "react";
import Form from "react-bootstrap/Form";
import {Button, ButtonGroup, InputGroup} from "react-bootstrap";

export const StringField = ({onChange, saveUpdater}) => {

    const refDom = React.useRef()

    function _setValue(value) {
        if (onChange) {
            onChange(value)
        }
        refDom.current.value = value
    }

    if (saveUpdater) {
        saveUpdater((value) => _setValue(value))
    }

    function handleChange(event) {
        _setValue(event.target.value)
    }

    return (
        <label>
            <Form.Control
                type="text"
                ref={refDom}
                onChange={handleChange}
            />
        </label>
    )
}

export const IntegerField = ({onChange, saveUpdater, min = 1, max = 10}) => {

    const refDom = React.useRef()
    const refValue = React.useRef(min)

    function _setValue(value) {
        let parsed = parseInt(value.toString())
        if (parsed < min || parsed > max) {
            parsed = refValue.current // === panedrone: reset to the old one
        }
        refValue.current = parsed
        value = parsed.toString()
        if (onChange) {
            onChange(value)
        }
        refDom.current.value = value
    }

    if (saveUpdater) {
        saveUpdater(_setValue)
    }

    function handleChange(event) {
        let value = event.target.value
        if (!value) {
            return // === panedrone: allow typing from scratch
        }
        _setValue(value);
    }

    function handleUp() {
        let current = refValue.current
        if (current >= max) {
            return
        }
        current = current + 1
        _setValue(current);
    }

    function handleDown() {
        let current = refValue.current
        if (current <= min) {
            return
        }
        current = current - 1
        _setValue(current);
    }

    // === panedrone: <input type="number" is buggy:
    //      - impossible to disable typing of not-numbers
    //      - with invalid values, "onChange" is not fired

    // return (
    //     <label>
    //         <input type="number" min="1" max="10" pattern="[0-9\s]" value={value)} onChange={handleChange}/>
    //     </label>
    // )

    const spinStyle = {padding: "0 6px", lineHeight: 1, fontSize: "0.7em"}

    return (
        <InputGroup size="sm">
            <Form.Control type="text" pattern="[0-9\s]"
                          ref={refDom} placeholder="Type a number..." onChange={handleChange}/>
            <ButtonGroup vertical>
                <Button type="button" variant="outline-secondary" style={spinStyle} onClick={handleUp}>
                    <i className="bi bi-caret-up-fill"/>
                </Button>
                <Button type="button" variant="outline-secondary" style={spinStyle} onClick={handleDown}>
                    <i className="bi bi-caret-down-fill"/>
                </Button>
            </ButtonGroup>
        </InputGroup>
    )
}

export const TextAreaField = ({onChange, saveUpdater}) => {

    const refTextArea = React.useRef()

    function _setValue(value) {
        if (onChange) {
            onChange(value)
        }
        refTextArea.current.value = value
    }

    if (saveUpdater) {
        saveUpdater((value) => _setValue(value))
    }

    function handleChange(event) {
        _setValue(event.target.value)
    }

    return (
        <label>
            <textarea cols="40" rows="10" ref={refTextArea} onChange={handleChange}></textarea>
        </label>
    )
}
