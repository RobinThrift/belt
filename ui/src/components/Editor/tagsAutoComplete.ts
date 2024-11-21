import type { Tag } from "@/domain/Tag"
import {
    type Completion,
    type CompletionContext,
    autocompletion,
} from "@codemirror/autocomplete"
import type { Extension } from "@codemirror/state"

export function tagsAutoComplete(tags: Tag[] = []): Extension {
    let completions: Completion[] = tags.map((tag) => ({ label: tag.tag }))

    return autocompletion({
        override: [
            (context: CompletionContext) => {
                let word = context.matchBefore(/#([\w/]+)?/)
                if (!word) return null
                if (word && word.from === word.to && !context.explicit) {
                    return null
                }
                return {
                    from: word?.from,
                    options: completions,
                }
            },
        ],
    })
}
