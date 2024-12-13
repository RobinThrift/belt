import { Loader } from "@/components/Loader"
import { Memo, type PartialMemoUpdate } from "@/components/Memo"
import type { Tag } from "@/domain/Tag"
import type { Filter } from "@/state/memolist"
import { useSetting } from "@/state/settings"
import React, { useCallback, useEffect } from "react"
import { useSingleMemoPageState } from "./state"

export interface SingleMemoPageProps {
    memoID: string
    onChangeFilters?: (filter: Filter) => void
}

export function SingleMemoPage(props: SingleMemoPageProps) {
    let [doubleClickToEdit] = useSetting("controls.doubleClickToEdit")

    let { state, actions } = useSingleMemoPageState(props.memoID)

    useEffect(() => {
        if (!state?.memo) {
            actions.load()
        }
    }, [state?.memo, actions.load])

    let onClickTag = useCallback(
        (tag: string) => {
            props.onChangeFilters?.({ tag })
        },
        [props.onChangeFilters],
    )

    let updateMemoCallback = useCallback(
        (memo: PartialMemoUpdate) => {
            actions.update(memo)
        },
        [actions.update],
    )

    let tags: Tag[] = []

    return (
        <div className="container mx-auto max-w-4xl">
            {(!state || state?.isLoading) && (
                <div className="flex justify-center items-center min-h-[200px]">
                    <Loader />
                </div>
            )}

            {state?.memo && (
                <Memo
                    memo={state?.memo}
                    tags={tags}
                    actions={{
                        link: false,
                        edit: !state.memo.isArchived && !state.memo.isDeleted,
                    }}
                    className="animate-in slide-in-from-bottom fade-in"
                    onClickTag={onClickTag}
                    updateMemo={updateMemoCallback}
                    doubleClickToEdit={doubleClickToEdit}
                />
            )}
        </div>
    )
}
