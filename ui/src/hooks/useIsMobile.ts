import * as React from "react"

const MOBILE_BREAKPOINT = 1024

export function useIsMobile() {
    const [isMobile, setIsMobile] = React.useState<boolean | undefined>(
        undefined,
    )

    React.useEffect(() => {
        let mql = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`)
        let onChange = () => {
            setIsMobile(window.innerWidth < MOBILE_BREAKPOINT)
        }
        mql.addEventListener("change", onChange)
        setIsMobile(window.innerWidth < MOBILE_BREAKPOINT)
        return () => mql.removeEventListener("change", onChange)
    }, [])

    return !!isMobile
}
