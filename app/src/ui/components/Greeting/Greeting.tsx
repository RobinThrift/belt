import clsx from "clsx"
import React from "react"

import { useT } from "@/ui/i18n"
import { useSetting } from "@/ui/settings"

export interface GreetingProps {
    className?: string
}

export function Greeting(props: GreetingProps) {
    let t = useT("components/Greeting")
    let [displayName] = useSetting("account.displayName")
    let now = new Date()
    let greeting = t.Evening

    if (now.getHours() < 12) {
        greeting = t.Morning
    }

    if (now.getHours() < 18) {
        greeting = t.Afternoon
    }

    return (
        <div className={clsx("greeting", props.className)}>
            {greeting}
            <em>{displayName}</em>
        </div>
    )
}
