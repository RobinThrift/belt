import type { ListMemosQuery, Memo, MemoID, MemoList } from "@/domain/Memo"
import type { Pagination } from "@/domain/Pagination"
import { format, formatRFC3339, parse } from "date-fns"
import { APIError, UnauthorizedError } from "./APIError"

export type Filter = ListMemosQuery

export async function get({
    id,
    baseURL = "",
    signal,
}: {
    id: string
    baseURL?: string
    signal?: AbortSignal
}): Promise<Memo> {
    let url = new URL(`${baseURL}/api/v1/memos/${id}`, globalThis.location.href)

    let res = await fetch(url, { signal })

    if (res.status === 401) {
        throw new UnauthorizedError("error fetching memo list")
    }

    if (res.status !== 200) {
        let err = await APIError.fromHTTPResponse(res)
        throw err.withPrefix("error fetching memo")
    }

    if (!res.ok) {
        throw new Error(`error fetching memo: ${res.status} ${res.statusText}`)
    }

    let memo = (await res.json()) as Memo

    return {
        ...memo,
        createdAt: new Date(memo.createdAt),
        updatedAt: new Date(memo.updatedAt),
    }
}

export async function list({
    pagination,
    baseURL = "",
    filter,
    signal,
}: {
    pagination: Pagination<Date>
    filter: Filter
    baseURL?: string
    signal?: AbortSignal
}): Promise<MemoList> {
    let url = new URL(`${baseURL}/api/v1/memos`, globalThis.location.href)
    addFilterToSearchParams(url.searchParams, filter, pagination)

    let res = await fetch(url, { signal })

    if (res.status === 401) {
        throw new UnauthorizedError("error fetching memo list")
    }

    if (res.status !== 200) {
        let err = await APIError.fromHTTPResponse(res)
        throw err.withPrefix("error fetching memo list")
    }

    if (!res.ok) {
        throw new Error(
            `unknown error fetching memo list: ${res.status} ${res.statusText}`,
        )
    }

    let list = (await res.json()) as MemoList

    list.items = list.items.map((memo) => ({
        ...memo,
        createdAt: new Date(memo.createdAt),
        updatedAt: new Date(memo.updatedAt),
    }))

    list.next = list.next && new Date(list.next)

    return list
}

export interface CreateMemoRequest {
    content: string
}

export async function create({
    memo,
    signal,
    baseURL = "",
}: {
    memo: CreateMemoRequest
    baseURL?: string
    signal?: AbortSignal
}): Promise<Memo> {
    let url = new URL(`${baseURL}/api/v1/memos`, globalThis.location.href)

    let res = await fetch(url, {
        signal,
        method: "POST",
        body: JSON.stringify(memo),
    })

    if (res.status === 401) {
        throw new UnauthorizedError("error creating memo")
    }

    if (res.status !== 201) {
        let err = await APIError.fromHTTPResponse(res)
        throw err.withPrefix("error creating memo")
    }

    if (!res.ok) {
        throw new Error(
            `unknown error creating memo: ${res.status} ${res.statusText}`,
        )
    }

    let created = await res.json()

    return {
        ...created,
        createdAt: new Date(created.createdAt),
        updatedAt: new Date(created.updatedAt),
    }
}

export interface UpdateMemoRequest {
    id: MemoID
    content?: string
    isArchived?: boolean
    isDeleted?: boolean
}

export async function update({
    memo,
    signal,
    baseURL = "",
}: {
    memo: UpdateMemoRequest
    baseURL?: string
    signal?: AbortSignal
}): Promise<void> {
    let url = new URL(
        `${baseURL}/api/v1/memos/${memo.id}`,
        globalThis.location.href,
    )

    let res = await fetch(url, {
        signal,
        method: "PATCH",
        body: JSON.stringify({
            content: memo.content,
            isArchived: memo.isArchived,
            isDeleted: memo.isDeleted,
        }),
    })

    if (res.status === 401) {
        throw new UnauthorizedError("error updating memo")
    }

    if (res.status !== 204) {
        let err = await APIError.fromHTTPResponse(res)
        throw err.withPrefix("error updating memo")
    }

    if (!res.ok) {
        throw new Error(
            `unknown error updating memo: ${res.status} ${res.statusText}`,
        )
    }
}

export interface FilterQueryParams {
    "filter[tag]"?: string
    "filter[content]"?: string
    "filter[created_at]"?: string
    "op[created_at]"?: string
    "filter[is_deleted]"?: string
    "filter[is_archived]"?: string
}

export function filterFromQuery(query: FilterQueryParams): Filter {
    let filter: Filter = {}

    if (query["filter[tag]"]) {
        filter.tag = query["filter[tag]"]
    }

    if (query["filter[content]"]) {
        filter.query = query["filter[content]"]
    }

    if (query["filter[is_deleted]"] === "true") {
        filter.isDeleted = true
    }

    if (query["filter[is_archived]"] === "true") {
        filter.isArchived = true
    }

    let createdAt = query["filter[created_at]"]
    if (!createdAt) {
        return filter
    }

    let opCreatedAt = query["op[created_at]"]
    if (opCreatedAt && opCreatedAt === "<=") {
        filter.startDate = parse(createdAt, "yyyy-MM-dd", new Date())
    } else {
        filter.exactDate = parse(createdAt, "yyyy-MM-dd", new Date())
    }

    return filter
}

export function filterToSearchParams(
    filter: Filter,
    pagination?: Pagination<Date>,
) {
    let searchParams = new URLSearchParams()
    addFilterToSearchParams(searchParams, filter, pagination)
    return searchParams
}

function addFilterToSearchParams(
    searchParams: URLSearchParams,
    filter: Filter,
    pagination?: Pagination<Date>,
) {
    if (pagination) {
        searchParams.set("page[size]", `${pagination.pageSize}`)

        if (pagination.after) {
            searchParams.set(
                "page[after]",
                `${formatRFC3339(pagination.after)}`,
            )
        }
    }

    if (filter.tag) {
        searchParams.set("filter[tag]", filter.tag)
    }

    if (filter.query) {
        searchParams.set("filter[content]", filter.query)
    }

    if (filter.exactDate) {
        searchParams.set(
            "filter[created_at]",
            format(filter.exactDate, "yyyy-MM-dd"),
        )
    }

    if (filter.startDate) {
        searchParams.set(
            "filter[created_at]",
            format(filter.startDate, "yyyy-MM-dd"),
        )
        searchParams.set("op[created_at]", "<=")
    }

    if (filter.isArchived) {
        searchParams.set("filter[is_archived]", "true")
    }

    if (filter.isDeleted) {
        searchParams.delete("filter[is_archived]")
        searchParams.set("filter[is_deleted]", "true")
    }

    return searchParams
}
