import { faker } from "@faker-js/faker"
import type { Meta, StoryObj } from "@storybook/react"
import { DateTime } from "./DateTime"

import "@/index.css"

const meta: Meta<typeof DateTime> = {
    title: "Components/DateTime",
    component: DateTime,
}

export default meta

type Story = StoryObj<typeof DateTime>

export const Basic: Story = {
    name: "DateTime",
    args: {
        date: faker.date.past({ years: 1 }),
    },
}

export const Relative: Story = {
    args: {
        date: faker.date.recent(),
        relative: true,
    },
}
