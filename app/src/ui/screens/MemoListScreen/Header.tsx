import clsx from "clsx"
import React from "react"

import type { ListMemosQuery as Filter } from "@/domain/Memo"
import { DateTime } from "@/ui/components/DateTime"
import { Greeting } from "@/ui/components/Greeting"
import { useT } from "@/ui/i18n"

export function Header({
    className,
    filter,
}: { className?: string; filter: Filter }) {
    let t = useT("components/MemoListHeader")

    if (!filter) {
        return (
            <header className="memo-list-header">
                <Greeting key="greeting" />
            </header>
        )
    }

    let children: React.ReactNode[] = []

    if (filter.isDeleted) {
        children.push(t.Deleted, <br key="br-is-deleted" />)
    }

    if (filter.isArchived) {
        children.push(t.Archived, <br key="br-is-archived" />)
    }

    if (filter.tag) {
        children.push(t.MemosForTag, <em key="tag">{`#${filter.tag}`}</em>)
    }

    if (filter.exactDate) {
        children.push(
            children.length === 0
                ? t.MemosForExactDateStandalone
                : t.MemosForExactDate,
            <em key="exactDate">
                <DateTime
                    date={filter.exactDate}
                    opts={{ dateStyle: "medium" }}
                />
            </em>,
        )
    }

    if (filter.query) {
        children.push(
            children.length === 0 ? t.MemosForQueryStandalone : t.MemosForQuery,
            <em key="query">{`"${filter.query}"`}</em>,
        )
    }

    if (children.length === 0) {
        children = [<Greeting key="greeting" />]
    }

    return (
        <header className={clsx("memo-list-header", className)}>
            {children}
        </header>
    )
}
