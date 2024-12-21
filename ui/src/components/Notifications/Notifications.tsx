import { Button } from "@/components/Button"
import { useT } from "@/i18n"
import type { Notification as NotificationT } from "@/notifications"
import { useNotifications } from "@/state/notifications"
import { Info, Warning, X } from "@phosphor-icons/react"
import * as Toast from "@radix-ui/react-toast"
import clsx from "clsx"
import React, { useMemo } from "react"

export interface NotificationsProps {
    className?: string
    durationMs?: number
}

export function Notifications(props: NotificationsProps) {
    let t = useT("components/Notifications")
    let { notifications, remove } = useNotifications()
    let durationMs = props.durationMs || 10000

    let toasts = useMemo(
        () =>
            notifications.map((n, i) =>
                n.requiresAction ? (
                    <PermanentNotification
                        key={n.title}
                        n={n}
                        ariaLabelCloseBtn={t.Dismiss}
                        onOpenChange={(open) =>
                            !open && setTimeout(() => remove(i), 100)
                        }
                    />
                ) : (
                    <Notification
                        key={n.title}
                        n={n}
                        ariaLabelCloseBtn={t.Dismiss}
                        onOpenChange={(open) =>
                            !open && setTimeout(() => remove(i), 100)
                        }
                    />
                ),
            ),
        [notifications, remove, t.Dismiss],
    )

    return (
        <Toast.Provider
            swipeDirection="right"
            duration={durationMs}
            label={t.Label}
        >
            <Toast.Viewport
                className={clsx("notifications-viewport", props.className)}
            >
                {toasts}
            </Toast.Viewport>
        </Toast.Provider>
    )
}

function Notification({
    n,
    ariaLabelCloseBtn,
    onOpenChange,
}: {
    n: NotificationT
    ariaLabelCloseBtn: string
    onOpenChange: (open: boolean) => void
}) {
    return (
        <Toast.Root
            className={clsx("notification", n.type)}
            onOpenChange={onOpenChange}
            duration={n.durationMs}
        >
            <Toast.Title className="notification-title">
                <div className="icon">
                    {n.type === "info" ? <Info /> : <Warning />}
                </div>
                {n.title}
            </Toast.Title>
            <Toast.Close
                className="notification-close"
                aria-label={ariaLabelCloseBtn}
            >
                <X aria-hidden />
            </Toast.Close>
            <Toast.Description className="notification-message">
                {n.message}
            </Toast.Description>

            <div className="notification-buttons">
                {n.buttons?.map((btn) => (
                    <Toast.Action
                        asChild
                        key={btn.ariaLabel}
                        altText={btn.ariaLabel || ""}
                    >
                        <Button {...btn} />
                    </Toast.Action>
                ))}
            </div>
        </Toast.Root>
    )
}

function PermanentNotification({
    n,
    ariaLabelCloseBtn,
    onOpenChange,
}: {
    n: NotificationT
    ariaLabelCloseBtn: string
    onOpenChange: (open: boolean) => void
}) {
    return (
        <div className={clsx("notification pointer-events-auto", n.type)}>
            <div className="notification-title">
                <div className="icon">
                    {n.type === "info" ? <Info /> : <Warning />}
                </div>
                {n.title}
            </div>
            <button
                type="button"
                className="notification-close"
                aria-label={ariaLabelCloseBtn}
                onClick={() => onOpenChange(false)}
            >
                <X />
            </button>
            <div className="notification-message">{n.message}</div>

            <div className="notification-buttons">
                {n.buttons?.map((btn) => (
                    <Toast.Action
                        asChild
                        key={btn.ariaLabel}
                        altText={btn.ariaLabel || ""}
                    >
                        <Button
                            {...btn}
                            onClick={(e) => {
                                onOpenChange(false)
                                btn.onClick?.(e)
                            }}
                        />
                    </Toast.Action>
                ))}
            </div>
        </div>
    )
}