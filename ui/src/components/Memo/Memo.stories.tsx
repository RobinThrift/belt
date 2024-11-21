import type { Tag } from "@/domain/Tag"
import { faker } from "@faker-js/faker"
import type { Meta, StoryObj } from "@storybook/react"
import React from "react"
import { Memo } from "./Memo"

import "@/index.css"

const meta: Meta<typeof Memo> = {
    title: "Components/Memo",
    component: Memo,
}

export default meta
type Story = StoryObj<typeof Memo>

export const Basic: Story = {
    name: "Memo",
    argTypes: {
        updateMemo: { action: "updateMemo" },
    },
    args: {
        doubleClickToEdit: true,
        memo: {
            id: "12345",
            content: `# Markdown Content (Heading 1)
${faker.lorem.lines({ min: 1, max: 10 })}

## Paragraphs with Tags (Heading 2)
${faker.lorem.lines({ min: 1, max: 10 })}

#tag-1 #tag/nesting

${faker.lorem.lines({ min: 1, max: 10 })}

### Lists (Heading 3)

- this is an
    - unordered list
    - with some
        - with nested items

- back to top level

1. orderd list
1. items
    - with nested regular
    - list
    - items
1. and ordered
    1. nested list
    1. items

#### Blockquote (Heading 4)
> ${faker.lorem.lines({ min: 1, max: 10 })}

And make note of the footnote[^fn1]

[^fn1]: Very Important content here

##### Text Styling (Heading 5)

This is some *emphasized* content and some **strong** content.
This text will be \`monospaced\`... hopefully. ~~Scratch this.~~

###### Heading 6
This is a [link](${faker.internet.url()}), an auto link (http://example.com) and an image will follow:

![image caption](${faker.image.url({ height: 200, width: 200 })})

***

### More Tests

The following  
will be on separate  
lines.


## Code Content
\`\`\`typescript
function parseMarkdown(raw: string): React.ReactNode | React.ReactNode[] {
    let ast = fromMarkdown(raw)
    return ast.children.map((c) => astNodeToJSX(c))
}
\`\`\`

## Tables

| left aligned | centre aligned | right aligned |
| :----------- | :------------: | ------------: |
| Cell A1      | Cell B1        | Cell C1       |
| Cell A2      | Cell B2        | Cell C2       |
| Cell A3      | Cell B3        | Cell C3       |

`,
            isArchived: false,
            createdAt: new Date(),
            updatedAt: new Date(),
        },
        tags: (() => {
            let tags: Tag[] = []

            for (let i = 0; i < 100; i++) {
                tags.push({
                    tag: `#${faker.word.noun()}/${faker.word.noun()}`,
                    count: 0,
                })
            }

            tags.sort()
            return tags
        })(),
    },
    render: (args) => (
        <div className="container mx-auto">
            <Memo {...args} />
        </div>
    ),
}
