import { UTCDateMini } from "@date-fns/utc"
import { format, parse, transpose } from "date-fns"

const sqliteDateTimeFormat = "yyyy-MM-dd HH:mm:ss'Z'"

export function dateFromSQLite(value: string): Date {
    return parse(value, sqliteDateTimeFormat, new UTCDateMini())
}

export function dateToSQLite(value: Date | string | undefined): string | null {
    if (!value) {
        return null
    }

    let d = value
    if (typeof d === "string") {
        d = new Date(value)
    }
    return format(transpose(d, UTCDateMini), sqliteDateTimeFormat)
}
