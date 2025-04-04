import type { UnlockController } from "@/control/UnlockController"
import { BaseContext } from "@/lib/context"
import type { StartListening } from "@/ui/state/rootStore"

import * as settings from "../settings"
import * as setup from "../setup"
import * as sync from "../sync"
import { slice } from "./slice"

export const registerEffects = (
    startListening: StartListening,
    {
        unlockCtrl,
    }: {
        unlockCtrl: UnlockController
    },
) => {
    startListening({
        actionCreator: slice.actions.unlock,
        effect: async (
            { payload },
            { cancelActiveListeners, getState, dispatch, signal },
        ) => {
            if (slice.selectors.isUnlocked(getState())) {
                return
            }

            cancelActiveListeners()

            let unlocked = await unlockCtrl.unlock(
                BaseContext.withSignal(signal),
                {
                    plaintextKeyData: payload.plaintextKeyData,
                    storeKey: payload.storeKey,
                    db: payload.db,
                },
            )

            if (signal.aborted) {
                return
            }

            if (!unlocked.ok) {
                dispatch(
                    slice.actions.setIsUnlocked({
                        error: unlocked.err,
                        isUnlocked: false,
                    }),
                )
                return
            }

            dispatch(
                slice.actions.setIsUnlocked({
                    isUnlocked: true,
                }),
            )

            dispatch(settings.actions.loadStart())
            dispatch(sync.actions.loadSyncInfo())
            dispatch(sync.actions.syncStart())
        },
    })

    startListening({
        actionCreator: setup.actions.setupCandidatePrivateCryptoKey,
        effect: async (
            { payload },
            { cancelActiveListeners, dispatch, signal },
        ) => {
            cancelActiveListeners()

            let ctx = BaseContext.withSignal(signal)

            await unlockCtrl.reset(ctx)

            let unlocked = await unlockCtrl.unlock(ctx, {
                plaintextKeyData: payload.plaintextKeyData,
            })

            if (signal.aborted) {
                return
            }

            if (!unlocked.ok) {
                dispatch(
                    setup.actions.setStep({
                        step: "configure-encryption",
                        error: unlocked.err,
                    }),
                )
                return
            }

            dispatch(
                slice.actions.setIsUnlocked({
                    isUnlocked: true,
                }),
            )

            dispatch(settings.actions.loadStart())
            dispatch(sync.actions.loadSyncInfo())
        },
    })
}
