import type { Tag } from "@/domain/Tag"
import { CaretRight, Hash } from "@phosphor-icons/react"
import clsx from "clsx"
import React, { useCallback, useMemo } from "react"
import TreeView, {
    type INode,
    type ITreeViewOnNodeSelectProps,
} from "react-accessible-treeview"

export interface TagTreeProps {
    className?: string
    tags: Tag[]
    selected?: string
    onSelect: (selected?: string) => void
}

export function TagTree({
    className,
    tags,
    selected,
    onSelect: propagateSelectionChange,
}: TagTreeProps) {
    let tagTree = useMemo(() => {
        let tree: Record<string, INode> = {
            root: { name: "", id: "root", children: [], parent: null },
        }

        tags.forEach((tag) => {
            let segments = tag.tag.replace("#", "").split("/")
            let id = ""
            segments.forEach((segment) => {
                let parent = id
                id = id === "" ? segment : `${id}/${segment}`
                let metadata = { count: id === tag.tag ? tag.count : 0 }
                if (!tree[id]) {
                    tree[id] = {
                        id,
                        name: segment,
                        parent: parent || "root",
                        children: [],
                        metadata,
                    }
                } else {
                    tree[id].metadata = metadata
                }
            })
        })

        let nodes: INode[] = []

        Object.values(tree).forEach((node) => {
            if (node.parent) {
                tree[node.parent]?.children.push(node.id)
            }
            nodes.push(node)
        })

        return nodes
    }, [tags])

    let onSelect = useCallback(
        ({ element }: ITreeViewOnNodeSelectProps) => {
            if (element.id === selected) {
                propagateSelectionChange(undefined)
            } else {
                propagateSelectionChange(element.id as string)
            }
        },
        [propagateSelectionChange, selected],
    )

    let expandedIDs = useMemo(() => {
        if (!selected || !tagTree.find((n) => n.id === selected)) {
            return []
        }

        return selected.split("/").reduce((ids, segment) => {
            if (ids.length === 0) {
                ids.push(segment)
            } else {
                ids.push(`${ids.at(-1)}/${segment}`)
            }
            return ids
        }, [] as string[])
    }, [selected, tagTree])

    if (tags.length === 0) {
        return null
    }

    return (
        <div className={clsx("tag-tree", className)}>
            <TreeView
                data={tagTree}
                className=""
                aria-label="tag tree"
                onNodeSelect={onSelect}
                selectedIds={selected ? [selected] : []}
                defaultExpandedIds={expandedIDs}
                togglableSelect={true}
                nodeRenderer={({
                    element,
                    getNodeProps,
                    level,
                    isExpanded,
                    isSelected,
                    handleExpand,
                }) => {
                    let { className, onClick } = getNodeProps()
                    return (
                        <div
                            className={className}
                            style={{ marginLeft: 20 * (level - 1) }}
                        >
                            <div
                                className={clsx(
                                    "flex gap-1 items-center cursor-pointer rounded-sm xs:hover:bg-subtle active:bg-subtle ps-1 pe-2",
                                    {
                                        "bg-primary-light! text-primary-contrast!":
                                            isSelected,
                                    },
                                )}
                            >
                                <Hash
                                    className={clsx({
                                        "text-subtle-dark": !isSelected,
                                        "text-primary-contrast": isSelected,
                                    })}
                                />
                                <button
                                    className="flex-1 appearance-none text-left"
                                    onClick={
                                        element.metadata?.count
                                            ? onClick
                                            : handleExpand
                                    }
                                    type="button"
                                >
                                    {element.name}
                                    {element.metadata?.count
                                        ? ` (${element.metadata.count})`
                                        : null}
                                </button>

                                {element.children.length !== 0 && (
                                    <button
                                        type="button"
                                        className="appearance-none"
                                        onClick={handleExpand}
                                    >
                                        <CaretRight
                                            className={clsx("transition", {
                                                "rotate-90": isExpanded,
                                            })}
                                        />
                                    </button>
                                )}
                            </div>
                        </div>
                    )
                }}
            />
        </div>
    )
}
