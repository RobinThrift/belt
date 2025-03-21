export type Result<T, E extends Error = Error> =
    | {
          ok: true
          value: T
      }
    | {
          ok: false
          err: E
      }

export type ResultValue<T> = {
    ok: true
    value: T
}

export type ResultErr<E extends Error = Error> = {
    ok: false
    err: E
}

export type AsyncResult<T, E extends Error = Error> = Promise<Result<T, E>>

export function fromThrowing<T, E extends Error = Error>(
    fn: () => T,
): Result<T, E> {
    try {
        return { ok: true, value: fn() }
    } catch (e) {
        return { ok: false, err: e as E }
    }
}

export function wrapThrowing<
    T,
    Fn extends (...args: any) => T,
    E extends Error = Error,
>(fn: Fn, ...args: Parameters<Fn>): Result<T, E> {
    try {
        return { ok: true, value: fn(...args) }
    } catch (e) {
        return { ok: false, err: e as E }
    }
}

export async function fromPromise<T, E extends Error = Error>(
    p: PromiseLike<T>,
): AsyncResult<T, E> {
    try {
        return { ok: true, value: await p }
    } catch (e) {
        return { ok: false, err: e as E }
    }
}

export async function fromAsyncFn<
    T,
    Fn extends () => Promise<T>,
    E extends Error = Error,
>(fn: Fn): AsyncResult<T, E> {
    try {
        return { ok: true, value: await fn() }
    } catch (e) {
        return { ok: false, err: e as E }
    }
}

export function Ok<T, E extends Error = Error>(
    value: PromiseLike<T>,
): AsyncResult<T, E>
export function Ok<T, E extends Error = Error>(value: T): Result<T, E>
export function Ok<T, E extends Error = Error>(
    value: T | PromiseLike<T>,
): Result<T, E> | AsyncResult<T, E> {
    if (
        value !== null &&
        typeof value === "object" &&
        "then" in value &&
        typeof value.then === "function"
    ) {
        return Promise.resolve(Ok(value as T))
    }

    return { ok: true, value: value as T }
}

export function Err<T, E extends Error = Error>(err: E): Result<T, E> {
    return { ok: false, err }
}

export function ErrAsync<T, E extends Error = Error>(
    err: E,
): AsyncResult<T, E> {
    return Promise.resolve(Err(err))
}

export function mapResult<A, B, E extends Error = Error>(
    r: AsyncResult<A, E>,
    fn: (v: A) => B,
): AsyncResult<B, E>
export function mapResult<A, B, E extends Error = Error>(
    r: Result<A, E>,
    fn: (v: A) => B,
): Result<B, E>
export function mapResult<A, B, E extends Error = Error>(
    r: Result<A, E> | AsyncResult<A, E>,
    fn: (v: A) => B,
): Result<B, E> | AsyncResult<B, E> {
    if ("then" in r && typeof r.then === "function") {
        return r.then((rr) => {
            if (!rr.ok) {
                return rr
            }
            return Ok(fn(rr.value))
        })
    }

    let rr = r as Result<A, E>

    if (!rr.ok) {
        return rr
    }
    return Ok(fn(rr.value))
}

export function fmtErr<A, E extends Error = Error>(
    fmt: string,
    r: AsyncResult<A, E>,
): AsyncResult<A, E>
export function fmtErr<A, E extends Error = Error>(
    fmt: string,
    r: Result<A, E>,
): Result<A, E>
export function fmtErr<A, E extends Error = Error>(
    fmt: string,
    r: Result<A, E> | AsyncResult<A, E>,
): Result<A, Error> | AsyncResult<A, Error> {
    if ("then" in r && typeof r.then === "function") {
        return r.then((rr) => {
            if (!rr.ok) {
                let err = new Error(fmt.replace("%w", rr.err.message), {
                    cause: rr.err,
                })
                err.stack = rr.err.stack
                return Err(err)
            }
            return rr
        })
    }

    let rr = r as Result<A, E>

    if (!rr.ok) {
        let err = new Error(fmt.replace("%w", rr.err.message), {
            cause: rr.err,
        })
        err.stack = rr.err.stack
        return Err(err)
    }
    return rr
}

export async function toPromise<T, E extends Error = Error>(
    r: AsyncResult<T, E>,
): Promise<T> {
    let res = await r
    if (!res.ok) {
        return Promise.reject(res.err)
    }

    return Promise.resolve(res.value)
}

export async function all<T>(
    results: Iterable<AsyncResult<T>>,
): AsyncResult<T[]> {
    let resolved = await Promise.all(results)
    let values: T[] = []
    for (let r of resolved) {
        if (!r.ok) {
            return Err(r.err)
        }
        values.push(r.value)
    }

    return Ok(values)
}

export function match<
    A,
    B,
    ErrA extends Error = Error,
    ErrB extends Error = ErrA,
    R extends Result<B, ErrB> = Result<B, ErrB>,
>(
    r: Result<A, ErrA>,
    mapVal: (v: A) => R,
    mapErr?: (err: ResultErr<ErrA>) => R,
): R
export function match<
    A,
    B,
    ErrA extends Error = Error,
    ErrB extends Error = ErrA,
    R extends AsyncResult<B, ErrB> = AsyncResult<B, ErrB>,
>(
    r: AsyncResult<A, ErrA>,
    mapVal: (v: A) => R,
    mapErr?: (err: ResultErr<ErrA>) => R,
): R
export function match<
    A,
    B,
    ErrA extends Error,
    ErrB extends Error,
    R extends Result<B, ErrB> | AsyncResult<B, ErrB>,
>(
    r: Result<A, ErrA> | AsyncResult<A, ErrA>,
    mapVal: (v: A) => R,
    mapErr?: (err: ResultErr<ErrA>) => R,
): R {
    if ("then" in r && typeof r.then === "function") {
        return r.then((rr) => {
            if (!rr.ok) {
                return (mapErr ? mapErr(rr) : rr) as R
            }
            return mapVal(rr.value)
        }) as R
    }

    let rr = r as Result<A, ErrA>

    if (!rr.ok) {
        return (mapErr ? mapErr(rr) : rr) as R
    }

    return mapVal(rr.value)
}
