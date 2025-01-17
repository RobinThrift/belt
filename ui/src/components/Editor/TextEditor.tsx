import type { Tag } from "@/domain/Tag"
import * as eventbus from "@/eventbus"
import { useAttachmentUploader } from "@/state/attachments"
import { useSetting, useTheme } from "@/state/settings"
import { markdown, markdownLanguage } from "@codemirror/lang-markdown"
import { languages } from "@codemirror/language-data"
import { SearchCursor } from "@codemirror/search"
import { Vim, vim } from "@replit/codemirror-vim"
import CodeMirror, { EditorView, type Extension } from "@uiw/react-codemirror"
import clsx from "clsx"
import React, { useCallback, useEffect, useMemo } from "react"
import { fileDropHandler } from "./cmFileDropHandler"
import { tagsAutoComplete } from "./tagsAutoComplete"

Vim.defineEx("write", "w", (cm: { cm6: EditorView }) => {
    eventbus.emit(`vim:write:${cm.cm6.dom.parentElement?.id ?? "global"}`)
})

Vim.defineEx("wquit", "wq", (cm: { cm6: EditorView }) => {
    eventbus.emit(`vim:write:${cm.cm6.dom.parentElement?.id ?? "global"}`)
})

Vim.defineEx("quit", "q", (cm: { cm6: EditorView }) => {
    eventbus.emit(`vim:quit:${cm.cm6.dom.parentElement?.id ?? "global"}`)
})

Vim.defineEx("cquit", "cq", (cm: { cm6: EditorView }) => {
    eventbus.emit(`vim:quit:${cm.cm6.dom.parentElement?.id ?? "global"}`)
})

Vim.map("A", "g$a")
Vim.map("I", "g0i")

Vim.map("j", "gj")
Vim.map("j", "gj", "visual")

Vim.map("k", "gk")
Vim.map("k", "gk", "visual")

export interface TextEditorProps {
    id: string
    tags: Tag[]
    content: string
    onChange: (content: string) => void
    onSave: () => void
    onCancel: () => void
    autoFocus?: boolean
    placeholder?: string
    placeCursorAt?: { x: number; y: number; snippet?: string }
}

export function TextEditor(props: TextEditorProps) {
    let { colours } = useTheme()

    useAttachmentUploader()

    let [enableVimKeybindings] = useSetting("controls.vim")

    let extensions: Extension[] = useMemo(() => {
        return [
            fileDropHandler(),
            enableVimKeybindings ? vim() : [],
            markdown({ base: markdownLanguage, codeLanguages: languages }),
            tagsAutoComplete(props.tags),
            EditorView.lineWrapping,
        ]
    }, [props.tags, enableVimKeybindings])

    useEffect(() => {
        return eventbus.on(`vim:write:${props.id}`, props.onSave)
    }, [props.id, props.onSave])

    useEffect(() => {
        return eventbus.on(`vim:quit:${props.id}`, props.onCancel)
    }, [props.id, props.onCancel])

    let onCreateEditor = useCallback(
        (view: EditorView) => {
            if (!props.placeCursorAt) {
                return
            }

            let pos: number | null = null
            if (props.placeCursorAt.snippet) {
                let cursor = new SearchCursor(
                    view.state.doc,
                    props.placeCursorAt.snippet,
                    view.posAtCoords(props.placeCursorAt, false) ?? 0,
                    view.state.doc.length,
                    (x) => x.toLowerCase(),
                )

                pos = cursor.next().value?.from
            }

            if (!pos || pos === -1) {
                pos = view.posAtCoords(props.placeCursorAt, false)
            }

            if (pos) {
                view.dispatch({
                    selection: {
                        anchor: pos,
                    },
                    scrollIntoView: true,
                })
            }
        },
        [props.placeCursorAt],
    )

    return (
        <CodeMirror
            id={props.id}
            className={clsx("text-editor", {
                "vim-enabled ": enableVimKeybindings,
            })}
            value={props.content}
            extensions={[extensions]}
            onChange={props.onChange}
            theme={colours.cm}
            autoFocus={props.autoFocus}
            placeholder={props.placeholder}
            onCreateEditor={onCreateEditor}
            basicSetup={{
                highlightActiveLine: false,
                lineNumbers: false,
                bracketMatching: true,
                indentOnInput: true,
                closeBrackets: true,
                foldGutter: false,
                highlightSelectionMatches: true,
                drawSelection: true,
                tabSize: 4,
                dropCursor: true,
            }}
        />
    )
}
