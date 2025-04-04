import React, { useCallback } from "react"

import type { PlaintextPrivateKey } from "@/lib/crypto"
import { Alert } from "@/ui/components/Alert"
import { BuildInfo } from "@/ui/components/BuildInfo"
import { Button } from "@/ui/components/Button"
import { ConveyorBeltText } from "@/ui/components/ConveyorBeltText/ConveyorBeltText"
import { Form } from "@/ui/components/Form"
import { PasswordIcon } from "@/ui/components/Icons"
import { Input } from "@/ui/components/Input"
import { SelectMode } from "@/ui/components/ThemeSwitcher"
import { useT } from "@/ui/i18n"

import { Checkbox } from "@/ui/components/Input/Checkbox"
import { useUnlockScreenState } from "./useUnlockScreenState"

export function UnlockScreen() {
    let t = useT("screens/Unlock")

    let { unlock, error } = useUnlockScreenState()

    let onSubmit = useCallback(
        (e: React.FormEvent) => {
            e.preventDefault()
            e.preventDefault()
            e.stopPropagation()

            let target = e.target as HTMLFormElement

            let plaintextPrivateKey = target.querySelector(
                "#password",
            ) as HTMLInputElement

            let storeKey = target.querySelector(
                "#store_key",
            ) as HTMLInputElement

            if (plaintextPrivateKey.value) {
                unlock({
                    plaintextKeyData:
                        plaintextPrivateKey.value as PlaintextPrivateKey,
                    storeKey: storeKey.dataset.state === "checked",
                })
            }
        },
        [unlock],
    )

    return (
        <div className="unlock-screen">
            <div className="unlock-screen-bg" aria-hidden>
                <div className="spot-3" />
                <div className="spot-2" />
                <div className="spot-1" />
                <div className="noise" />
            </div>

            <ConveyorBeltText
                className="logo"
                start="Co"
                middle="n"
                end="veyor"
                aria-hidden
            />

            <header className="unlock-screen-header">
                <SelectMode className="mode-select" />
            </header>

            <div className="unlock-window-positioner">
                <div className="unlock-window">
                    <h1>{t.Title}</h1>

                    <Form
                        action="#"
                        onSubmit={onSubmit}
                        className="p-4 space-y-4"
                    >
                        <Input
                            name="password"
                            icon={<PasswordIcon />}
                            label={t.PrivateKeyLabel}
                            type="password"
                            autoComplete="password"
                            required={true}
                            placeholder={t.PrivateKeyLabel}
                            messages={t}
                        />

                        <Checkbox name="store_key" label={t.StoreKeyLabel} />

                        {error && (
                            <Alert variant="danger">
                                {error.name}: {error.message}
                                {error.stack && (
                                    <pre>
                                        <code>{error.stack}</code>
                                    </pre>
                                )}
                            </Alert>
                        )}

                        <div className="flex items-center justify-end mt-4">
                            <Button
                                type="submit"
                                variant="primary"
                                size="lg"
                                className="w-full"
                            >
                                {t.UnlockButton}
                            </Button>
                        </div>
                    </Form>
                </div>
            </div>
            <footer className="unlock-footer">
                <BuildInfo />
            </footer>
        </div>
    )
}
