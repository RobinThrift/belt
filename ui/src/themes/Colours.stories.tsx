import { Button } from "@/components/Button"
import type { Meta, StoryObj } from "@storybook/react"
import React from "react"

import defaultCSS from "./default.css?inline"

const meta: Meta = {
    title: "Belt/Colours",
}

export default meta
type Story = StoryObj

export const Colours: Story = {
    render: () => {
        let start = defaultCSS.indexOf("{")
        let end = defaultCSS.indexOf("--btn-")
        let varnames = defaultCSS
            .substring(start + 2, end)
            .split("\n")
            .filter((l) => l.length !== 0)
            .map((l) => {
                let start = l.indexOf("-")
                let end = l.indexOf(":")
                return l.substring(start, end)
            })
        console.log("varnames", varnames)

        let swatches = varnames.map((varname) => (
            <div
                key={varname}
                className="flex justify-start items-start rounded-lg size-20 p-2"
                style={{
                    backgroundColor: `rgb(var(${varname}))`,
                }}
            >
                <span
                    className="text-xs"
                    style={{
                        color: `rgb(var(${varname}-contrast))`,
                    }}
                >
                    {varname}
                </span>
            </div>
        ))

        return (
            <div className="grid grid-cols-2 gap-2">
                <div className="p-2 border border-subtle rounded-lg flex gap-2 flex-wrap w-fit">
                    {swatches}
                </div>

                <div className="p-2 border border-subtle rounded-lg grid grid-cols-4 gap-2 w-full">
                    <Button variant="regular">Regular</Button>
                    <Button variant="primary">Primary</Button>
                    <Button variant="danger">Danger</Button>
                    <Button variant="success">Success</Button>

                    <Button variant="regular" outline={true}>
                        Outline Regular
                    </Button>
                    <Button variant="primary" outline={true}>
                        Outline Primary
                    </Button>
                    <Button variant="danger" outline={true}>
                        Outline Danger
                    </Button>
                    <Button variant="success" outline={true}>
                        Outline Success
                    </Button>

                    <Button variant="regular" plain={true}>
                        Regular (Plain)
                    </Button>
                    <Button variant="primary" plain={true}>
                        Primary (Plain)
                    </Button>
                    <Button variant="danger" plain={true}>
                        Danger (Plain)
                    </Button>
                    <Button variant="success" plain={true}>
                        Success (Plain)
                    </Button>
                </div>
            </div>
        )
    },
}