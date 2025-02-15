import type { Preview } from "@storybook/react"
import { initialize, mswLoader } from "msw-storybook-addon"
import React, { useEffect } from "react"
import { mockAPI } from "./mockapi"
import type { ServerData } from "../src/App/ServerData"
import { Theme } from "../src/components/Theme"
import { history } from "../src/external/history"
import { configureRootStore } from "../src/state"
import { faker } from "@faker-js/faker"

// Initialize MSW
initialize({
    onUnhandledRequest: "bypass",
})

let _mockRootStore: ReturnType<typeof configureRootStore> | undefined

let serverData: ServerData = {
    account: {
        username: "user",
        displayName: faker.person.firstName(),
    },

    settings: {
        locale: {
            language:
                localStorage.getItem("belt.settings.locale.language") ?? "en",
            region: localStorage.getItem("belt.settings.locale.region") ?? "gb",
        },

        theme: {
            colourScheme:
                localStorage.getItem("belt.settings.theme.colourScheme") ??
                "default",
            mode: localStorage.getItem("belt.settings.theme.mode") ?? "auto",
            icon: localStorage.getItem("belt.settings.theme.icon") ?? "default",
            listLayout: "masonry",
        },
        controls: {
            vim: JSON.parse(
                localStorage.getItem("belt.settings.controls.vim") ?? "true",
            ),
            doubleClickToEdit: JSON.parse(
                localStorage.getItem(
                    "belt.settings.controls.doubleClickToEdit",
                ) ?? "true",
            ),
        },
    },

    components: {
        LoginPage: {},
        LoginChangePasswordPage: {},
        SettingsPage: { validationErrors: {} },
    },

    buildInfo: {
        version: "vDEV",
        commitHash: "b23121",
        commitDate: "2024-12-07T17:26:43Z",
        projectLink: "https://github.com/RobinThrift/belt",
        goVersion: "1.23.4",
    },
}

const preview: Preview = {
    parameters: {
        actions: { argTypesRegex: "^on[A-Z].*" },
        msw: {
            handlers: mockAPI,
        },

        viewport: {
            viewports: {
                phone: {
                    name: "Phone",
                    styles: {
                        height: "852px",
                        width: "393px",
                    },
                    type: "mobile",
                },
                xs: {
                    name: "Small Window",
                    styles: {
                        width: "768px",
                        height: "900px",
                    },
                    type: "desktop",
                },
                tablet: {
                    name: "Tablet",
                    styles: {
                        height: "1024px",
                        width: "768px",
                    },
                    type: "desktop",
                },
                sm: {
                    name: "Small",
                    styles: {
                        width: "1440px",
                        height: "900px",
                    },
                    type: "desktop",
                },
                md: {
                    name: "Medium",
                    styles: {
                        width: "1680px",
                        height: "1050px",
                    },
                    type: "desktop",
                },
                lg: {
                    name: "Large",
                    styles: {
                        width: "2100px",
                        height: "1280px",
                    },
                    type: "desktop",
                },
            },
        },
    },

    loaders: [mswLoader],

    decorators: [
        (Story, { globals: { themeMode, themeColours, serverData } }) => {
            useEffect(() => {
                localStorage.setItem(
                    "belt.settings.theme.colourScheme",
                    themeColours,
                )
            }, [themeColours])

            useEffect(() => {
                localStorage.setItem("belt.settings.theme.mode", themeMode)
            }, [themeMode])

            if (!document.getElementById("__belt_ui_data__")) {
                let uiElement = document.createElement("script")
                uiElement.type = "belt_ui/data"
                uiElement.id = "__belt_ui_data__"
                uiElement.innerHTML = JSON.stringify(serverData)
                document.body.prepend(uiElement)
            }

            if (!document.querySelector("meta[name=theme-color]")) {
                let metaThemeEl = document.createElement("meta")
                metaThemeEl.name = "theme-color"
                metaThemeEl.content = ""
                document.head.prepend(metaThemeEl)
            }

            return (
                <Theme colourScheme={themeColours} mode={themeMode}>
                    <Story />
                </Theme>
            )
        },
    ],

    initialGlobals: {
        serverData,
        configureMockRootStore: () => {
            if (_mockRootStore) {
                return _mockRootStore
            }

            _mockRootStore = configureRootStore({
                router: { href: history.current },
                ...serverData,
            })
            return _mockRootStore
        },
    },

    globalTypes: {
        themeColours: {
            description: "Colour Scheme",
            defaultValue:
                localStorage.getItem("belt.settings.theme.colourScheme") ??
                "default",
            toolbar: {
                title: "Default",
                icon: "contrast",
                items: [
                    { value: "default", title: "Default" },
                    { value: "warm", title: "Warm" },
                    { value: "rosepine", title: "Rosé Pine" },
                ],
                dynamicTitle: true,
            },
        },
        themeMode: {
            description: "Mode",
            defaultValue:
                localStorage.getItem("belt.settings.theme.mode") ?? "auto",
            toolbar: {
                title: "Auto",
                icon: "lightning",
                items: [
                    { value: "auto", title: "Auto" },
                    { value: "dark", title: "Dark" },
                    { value: "light", title: "Light" },
                ],
                dynamicTitle: true,
            },
        },
    },
}

export default preview
