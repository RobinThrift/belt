import * as ThumbHash from "./thumbhash"

export async function thumbhashFromFile(file: File): Promise<string> {
    let { promise, resolve, reject } = Promise.withResolvers()

    let img = new Image()
    img.onload = resolve
    img.onerror = (err) => reject(err)
    img.src = URL.createObjectURL(file)

    await promise

    let canvas = document.createElement("canvas")
    let context = canvas.getContext("2d")
    if (!context) {
        throw new Error(
            "error calculating thumbhash: error getting canvas context",
        )
    }

    let scale = 100 / Math.max(img.width, img.height)
    canvas.width = Math.round(img.width * scale)
    canvas.height = Math.round(img.height * scale)
    context.drawImage(img, 0, 0, canvas.width, canvas.height)
    let pixels = context.getImageData(0, 0, canvas.width, canvas.height)
    let binaryThumbHash = ThumbHash.rgbaToThumbHash(
        pixels.width,
        pixels.height,
        pixels.data,
    )

    return btoa(String.fromCharCode(...binaryThumbHash))
}

export function thumbhashToDataURL(str: string): string {
    return ThumbHash.thumbHashToDataURL(
        new Uint8Array(
            atob(str)
                .split("")
                .map((x) => x.charCodeAt(0)),
        ),
    )
}
