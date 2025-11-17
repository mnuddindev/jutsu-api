// wasm-bridge.js - Complete Node.js bridge for WASM execution
// This file is embedded into the Go binary and written to disk on demand.

const https = require('https');
const http = require('http');
const { URL } = require('url');
const util = require('util');
const { webcrypto } = require('crypto');

const userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0";

async function fetchData(targetUrl, headers = {}) {
    return new Promise((resolve, reject) => {
        const parsedUrl = new URL(targetUrl);
        const protocol = parsedUrl.protocol === 'https:' ? https : http;

        const options = {
            hostname: parsedUrl.hostname,
            port: parsedUrl.port,
            path: parsedUrl.pathname + parsedUrl.search,
            method: 'GET',
            headers: {
                'User-Agent': userAgent,
                ...headers,
            },
        };

        const req = protocol.request(options, (res) => {
            const chunks = [];
            res.on('data', (chunk) => chunks.push(chunk));
            res.on('end', () => resolve(Buffer.concat(chunks)));
            res.on('error', reject);
        });

        req.on('error', reject);
        req.end();
    });
}

let wasm;
let arr = new Array(128).fill(void 0);
arr.push(void 0, null, true, false);
let size = 0;
let memoryBuff;
let dataView;
let pointer = arr.length;

function get(index) {
    return arr[index];
}

function isDetached(buffer) {
    if (buffer.byteLength === 0) {
        const formatted = util.format(buffer);
        return formatted.includes("detached");
    }
    return false;
}

function getMemBuff() {
    return (memoryBuff =
        null !== memoryBuff && 0 !== memoryBuff.byteLength
            ? memoryBuff
            : new Uint8Array(wasm.memory.buffer));
}

function getDataView() {
    return (dataView =
        dataView === null ||
        isDetached(dataView.buffer) ||
        dataView.buffer !== wasm.memory.buffer
            ? new DataView(wasm.memory.buffer)
            : dataView);
}

function shift(QP) {
    QP < 132 || ((arr[QP] = pointer), (pointer = QP));
}

function shiftGet(QP) {
    const Qn = get(QP);
    shift(QP);
    return Qn;
}

function addToStack(item) {
    pointer === arr.length && arr.push(arr.length + 1);
    const Qn = pointer;
    pointer = arr[Qn];
    arr[Qn] = item;
    return Qn;
}

const encoder = new TextEncoder();
const decoder = new TextDecoder("utf-8", { fatal: true, ignoreBOM: true });

function parse(text, func, func2) {
    if (func2 === void 0) {
        const encoded = encoder.encode(text);
        const parsedIndex = func(encoded.length, 1) >>> 0;
        getMemBuff()
            .subarray(parsedIndex, parsedIndex + encoded.length)
            .set(encoded);
        size = encoded.length;
        return parsedIndex;
    }

    let len = text.length;
    let parsedLen = func(len, 1) >>> 0;
    const newArr = getMemBuff();
    let i = 0;
    for (; i < len; i++) {
        const char = text.charCodeAt(i);
        if (char > 127) break;
        newArr[parsedLen + i] = char;
    }
    if (i !== len) {
        if (i !== 0) {
            text = text.slice(i);
        }
        parsedLen = func2(parsedLen, len, (len = i + 3 * text.length), 1) >>> 0;
        const encoded = getMemBuff().subarray(parsedLen + i, parsedLen + len);
        i += encoder.encodeInto(text, encoded).written;
        parsedLen = func2(parsedLen, len, i, 1) >>> 0;
    }
    size = i;
    return parsedLen;
}

function decodeSub(index, offset) {
    index >>>= 0;
    return decoder.decode(getMemBuff().subarray(index, index + offset));
}

function isNull(value) {
    return value == null;
}

function applyToWindow(func, args, fakeWindow) {
    try {
        return func.apply(fakeWindow, args);
    } catch (error) {
        wasm.__wbindgen_export_6(addToStack(error));
    }
}

function addCallback(QP, Qn, QT, func) {
    const Qx = { a: QP, b: Qn, cnt: 1, dtor: QT };
    const wrapper = (...Qw) => {
        Qx.cnt++;
        try {
            return func(Qx.a, Qx.b, ...Qw);
        } finally {
            if (--Qx.cnt === 0) {
                wasm.__wbindgen_export_2.get(Qx.dtor)(Qx.a, Qx.b);
                Qx.a = 0;
            }
        }
    };
    wrapper.original = Qx;
    return wrapper;
}

function export3(QP, Qn) {
    return shiftGet(wasm.__wbindgen_export_3(QP, Qn));
}

function export4(Qy, QO, QX) {
    wasm.__wbindgen_export_4(Qy, QO, addToStack(QX));
}

function export5(QP, Qn) {
    wasm.__wbindgen_export_5(QP, Qn);
}

function encodeView(QP, Qn) {
    Qn = Qn(+QP.length, 1) >>> 0;
    getMemBuff().set(QP, Qn);
    size = QP.length;
    return Qn;
}

async function instantiateModule(QP, Qn) {
    let instance;
    if (typeof Response === "function" && QP instanceof Response) {
        const buf = await QP.arrayBuffer();
        instance = await WebAssembly.instantiate(buf, Qn);
        return Object.assign(instance, { bytes: buf });
    }
    instance = await WebAssembly.instantiate(QP, Qn);
    if (instance instanceof WebAssembly.Instance) {
        return { instance, module: QP };
    }
    return instance;
}

function initWasm(fakeWindow, meta, imageData, canvas, nodeList) {
    const wasmObj = {
        wbg: {
            __wbindgen_is_function: (index) => typeof get(index) === "function",
            __wbindgen_is_string: (index) => typeof get(index) === "string",
            __wbindgen_is_object: (index) => {
                const obj = get(index);
                return typeof obj === "object" && obj !== null;
            },
            __wbindgen_number_get: (offset, index) => {
                const value = get(index);
                const view = getDataView();
                if (isNull(value)) {
                    view.setFloat64(offset + 8, 0, true);
                    view.setInt32(offset, 0, true);
                } else {
                    view.setFloat64(offset + 8, value, true);
                    view.setInt32(offset, 1, true);
                }
            },
            __wbindgen_string_get: (offset, index) => {
                const str = get(index);
                const val = parse(str, wasm.__wbindgen_export_0, wasm.__wbindgen_export_1);
                const view = getDataView();
                view.setInt32(offset + 4, size, true);
                view.setInt32(offset, val, true);
            },
            __wbindgen_object_drop_ref: (index) => shiftGet(index),
            __wbindgen_cb_drop: (index) => {
                const original = shiftGet(index).original;
                return original.cnt-- === 1 ? !(original.a = 0) : false;
            },
            __wbindgen_string_new: (index, offset) => addToStack(decodeSub(index, offset)),
            __wbindgen_is_null: (index) => get(index) === null,
            __wbindgen_is_undefined: (index) => get(index) === void 0,
            __wbindgen_boolean_get: (index) => {
                const value = get(index);
                return typeof value === "boolean" ? (value ? 1 : 0) : 2;
            },
            __wbg_instanceof_CanvasRenderingContext2d_4ec30ddd3f29f8f9: () => true,
            __wbg_subarray_adc418253d76e2f1: (index, from, to) =>
                addToStack(get(index).subarray(from >>> 0, to >>> 0)),
            __wbg_randomFillSync_5c9c955aa56b6049: () => {},
            __wbg_getRandomValues_3aa56aa6edec874c: () =>
                applyToWindow(function (index1, index2) {
                    get(index1).getRandomValues(get(index2));
                }, arguments, fakeWindow),
            __wbg_msCrypto_eb05e62b530a1508: (index) => addToStack(get(index).msCrypto),
            __wbg_toString_6eb7c1f755c00453: () => addToStack("[object Storage]"),
            __wbg_toString_139023ab33acec36: (index) => addToStack(get(index).toString()),
            __wbg_require_cca90b1a94a0255b: () =>
                applyToWindow(() => addToStack(module.require), arguments, fakeWindow),
            __wbg_crypto_1d1f22824a6a080c: (index) => addToStack(get(index).crypto),
            __wbg_process_4a72847cc503995b: (index) => addToStack(get(index).process),
            __wbg_versions_f686565e586dd935: (index) => addToStack(get(index).versions),
            __wbg_node_104a2ff8d6ea03a2: (index) => addToStack(get(index).node),
            __wbg_localStorage_3d538af21ea07fcc: () =>
                applyToWindow(() => {
                    const data = fakeWindow.localStorage;
                    return isNull(data) ? 0 : addToStack(data);
                }, arguments, fakeWindow),
            __wbg_setfillStyle_59f426135f52910f: () => {},
            __wbg_setshadowBlur_229c56539d02f401: () => {},
            __wbg_setshadowColor_340d5290cdc4ae9d: () => {},
            __wbg_setfont_16d6e31e06a420a5: () => {},
            __wbg_settextBaseline_c3266d3bd4a6695c: () => {},
            __wbg_drawImage_cb13768a1bdc04bd: () => {},
            __wbg_getImageData_66269d289f37d3c7: () =>
                applyToWindow(() => addToStack(imageData), arguments, fakeWindow),
            __wbg_rect_2fa1df87ef638738: () => {},
            __wbg_fillRect_4dd28e628381d240: () => {},
            __wbg_fillText_07e5da9e41652f20: () => {},
            __wbg_setProperty_5144ddce66bbde41: () => {},
            __wbg_createElement_03cf347ddad1c8c0: () =>
                applyToWindow(() => addToStack(canvas), arguments, fakeWindow),
            __wbg_querySelector_118a0639aa1f51cd: () =>
                applyToWindow(() => addToStack(meta), arguments, fakeWindow),
            __wbg_querySelectorAll_50c79cd4f7573825: () =>
                applyToWindow(() => addToStack(nodeList), arguments, fakeWindow),
            __wbg_getAttribute_706ae88bd37410fa: (offset, index) => {
                const attr = meta.content;
                const val = isNull(attr)
                    ? 0
                    : parse(attr, wasm.__wbindgen_export_0, wasm.__wbindgen_export_1);
                const view = getDataView();
                view.setInt32(offset + 4, size, true);
                view.setInt32(offset, val, true);
            },
            __wbg_target_6795373f170fd786: (index) => {
                const target = get(index).target;
                return isNull(target) ? 0 : addToStack(target);
            },
            __wbg_addEventListener_f984e99465a6a7f4: () => {},
            __wbg_instanceof_HtmlCanvasElement_1e81f71f630e46bc: () => true,
            __wbg_setwidth_233645b297bb3318: (index, value) => {
                get(index).width = value >>> 0;
            },
            __wbg_setheight_fcb491cf54e3527c: (index, value) => {
                get(index).height = value >>> 0;
            },
            __wbg_getContext_dfc91ab0837db1d1: () =>
                applyToWindow((index) => addToStack(get(index).context2d), arguments, fakeWindow),
            __wbg_toDataURL_97b108dd1a4b7454: () =>
                applyToWindow((offset) => {
                    const dataUrl = parse(
                        "data:image/png;base64,",
                        wasm.__wbindgen_export_0,
                        wasm.__wbindgen_export_1,
                    );
                    const view = getDataView();
                    view.setInt32(offset + 4, size, true);
                    view.setInt32(offset, dataUrl, true);
                }, arguments, fakeWindow),
            __wbg_instanceof_HtmlDocument_1100f8a983ca79f9: () => true,
            __wbg_style_ca229e3326b3c3fb: (index) => addToStack(get(index).style),
            __wbg_instanceof_HtmlImageElement_9c82d4e3651a8533: () => true,
            __wbg_src_87a0e38af6229364: (offset, index) => {
                const src = parse(
                    get(index).src,
                    wasm.__wbindgen_export_0,
                    wasm.__wbindgen_export_1,
                );
                const view = getDataView();
                view.setInt32(offset + 4, size, true);
                view.setInt32(offset, src, true);
            },
            __wbg_width_e1a38bdd483e1283: (index) => get(index).width,
            __wbg_height_e4cc2294187313c9: (index) => get(index).height,
            __wbg_complete_1162c2697406af11: (index) => get(index).complete,
            __wbg_data_d34dc554f90b8652: (offset, index) => {
                const data = encodeView(get(index).data, wasm.__wbindgen_export_0);
                const view = getDataView();
                view.setInt32(offset + 4, size, true);
                view.setInt32(offset, data, true);
            },
            __wbg_origin_305402044aa148ce: () =>
                applyToWindow((offset, index) => {
                    const origin = parse(
                        get(index).origin,
                        wasm.__wbindgen_export_0,
                        wasm.__wbindgen_export_1,
                    );
                    const view = getDataView();
                    view.setInt32(offset + 4, size, true);
                    view.setInt32(offset, origin, true);
                }, arguments, fakeWindow),
            __wbg_length_8a9352f7b7360c37: (index) => get(index).length,
            __wbg_get_c30ae0782d86747f: (index) => {
                const image = get(index).image;
                return isNull(image) ? 0 : addToStack(image);
            },
            __wbg_timeOrigin_f462952854d802ec: (index) => get(index).timeOrigin,
            __wbg_instanceof_Window_cee7a886d55e7df5: () => true,
            __wbg_document_eb7fd66bde3ee213: (index) => {
                const doc = get(index).document;
                return isNull(doc) ? 0 : addToStack(doc);
            },
            __wbg_location_b17760ac7977a47a: (index) => addToStack(get(index).location),
            __wbg_performance_4ca1873776fdb3d2: (index) => {
                const perf = get(index).performance;
                return isNull(perf) ? 0 : addToStack(perf);
            },
            __wbg_origin_e1f8acdeb3a39a2b: (offset, index) => {
                const origin = parse(
                    get(index).origin,
                    wasm.__wbindgen_export_0,
                    wasm.__wbindgen_export_1,
                );
                const view = getDataView();
                view.setInt32(offset + 4, size, true);
                view.setInt32(offset, origin, true);
            },
            __wbg_get_8986951b1ee310e0: (index, decode1, decode2) => {
                const data = get(index)[decodeSub(decode1, decode2)];
                return isNull(data) ? 0 : addToStack(data);
            },
            __wbg_setTimeout_6ed7182ebad5d297: () =>
                applyToWindow(() => 7, arguments, fakeWindow),
            __wbg_self_05040bd9523805b9: () =>
                applyToWindow(() => addToStack(fakeWindow), arguments, fakeWindow),
            __wbg_window_adc720039f2cb14f: () =>
                applyToWindow(() => addToStack(fakeWindow), arguments, fakeWindow),
            __wbg_globalThis_622105db80c1457d: () =>
                applyToWindow(() => addToStack(fakeWindow), arguments, fakeWindow),
            __wbg_global_f56b013ed9bcf359: () =>
                applyToWindow(() => addToStack(fakeWindow), arguments, fakeWindow),
            __wbg_newnoargs_cfecb3965268594c: (index, offset) =>
                addToStack(new Function(decodeSub(index, offset))),
            __wbindgen_object_clone_ref: (index) => addToStack(get(index)),
            __wbg_eval_c824e170787ad184: () =>
                applyToWindow((index, offset) => {
                    const fakeStr = "fake_" + decodeSub(index, offset);
                    // eslint-disable-next-line no-eval
                    const ev = eval(fakeStr);
                    return addToStack(ev);
                }, arguments, fakeWindow),
            __wbg_call_3f093dd26d5569f8: () =>
                applyToWindow((index, index2) => addToStack(get(index).call(get(index2))), arguments, fakeWindow),
            __wbg_call_67f2111acd2dfdb6: () =>
                applyToWindow((index, index2, index3) =>
                    addToStack(get(index).call(get(index2), get(index3))), arguments, fakeWindow),
            __wbg_set_961700853a212a39: () =>
                applyToWindow((index, index2, index3) =>
                    Reflect.set(get(index), get(index2), get(index3)), arguments, fakeWindow),
            __wbg_buffer_b914fb8b50ebbc3e: (index) => addToStack(get(index).buffer),
            __wbg_newwithbyteoffsetandlength_0de9ee56e9f6ee6e: (index, val, val2) =>
                addToStack(new Uint8Array(get(index), val >>> 0, val2 >>> 0)),
            __wbg_newwithlength_0d03cef43b68a530: (length) =>
                addToStack(new Uint8Array(length >>> 0)),
            __wbg_new_b1f2d6842d615181: (index) => addToStack(new Uint8Array(get(index))),
            __wbg_buffer_67e624f5a0ab2319: (index) => addToStack(get(index).buffer),
            __wbg_length_21c4b0ae73cba59d: (index) => get(index).length,
            __wbg_set_7d988c98e6ced92d: (index, index2, val) => {
                get(index).set(get(index2), val >>> 0);
            },
            __wbindgen_debug_string: () => {},
            __wbindgen_throw: (index, offset) => {
                throw new Error(decodeSub(index, offset));
            },
            __wbindgen_memory: () => addToStack(wasm.memory),
            __wbindgen_closure_wrapper117: (Qn, QT) => addToStack(addCallback(Qn, QT, 2, export3)),
            __wbindgen_closure_wrapper119: (Qn, QT) => addToStack(addCallback(Qn, QT, 2, export4)),
            __wbindgen_closure_wrapper121: (Qn, QT) => addToStack(addCallback(Qn, QT, 2, export5)),
            __wbindgen_closure_wrapper123: (Qn, QT) => addToStack(addCallback(Qn, QT, 9, export4)),
        },
    };
    return wasmObj;
}

function assignWasm(resp) {
    wasm = resp.exports;
    dataView = null;
    memoryBuff = null;
    return wasm;
}

function initSync(bytes, imports) {
    if (!(bytes instanceof WebAssembly.Module)) {
        bytes = new WebAssembly.Module(bytes);
    }
    return assignWasm(new WebAssembly.Instance(bytes, imports));
}

async function loadWasmModule(wasmBytes, imports) {
    const module = await WebAssembly.compile(wasmBytes);
    const instance = await WebAssembly.instantiate(module, imports);
    wasm = instance.exports;
    dataView = null;
    memoryBuff = null;
    return instance;
}

async function runWasm(baseURL, wasmURL, imageURL, metaContent) {
    const wasmBytes = await fetchData(wasmURL);
    const dateNow = Date.now();
    const fakeWindow = {
        localStorage: {
            setItem: function (item, value) {
                this[item] = value;
            },
        },
        navigator: { webdriver: false, userAgent },
        document: { cookie: "" },
        origin: baseURL,
        location: { href: baseURL, origin: baseURL },
        performance: { timeOrigin: dateNow },
        crypto: webcrypto,
        msCrypto: webcrypto,
        bytes: wasmBytes,
        jwt_plugin: function (bytes) {
            this.bytes = bytes;
        },
        navigate: function () {
            return "";
        },
    };

    const meta = { content: metaContent };
    const imageBytes = await fetchData(imageURL);
    const imageData = {
        height: 50,
        width: 65,
        data: new Uint8ClampedArray(imageBytes),
    };
    const canvas = { baseUrl: baseURL, width: 0, height: 0, context2d: {}, style: {} };
    const nodeList = { image: { src: "", height: 50, width: 65, complete: true }, length: 1 };

    const imports = initWasm(fakeWindow, meta, imageData, canvas, nodeList);
    await loadWasmModule(wasmBytes, imports);
    if (wasm.groot) {
        wasm.groot();
    }

    return {
        pid: fakeWindow.pid || fakeWindow.localStorage.pid || "",
        kversion: fakeWindow.localStorage.kversion || "",
        kid: fakeWindow.localStorage.kid || "",
        bytes: Array.from(new Uint8Array(wasmBytes)),
    };
}

async function main() {
    try {
        const args = process.argv.slice(2);
        if (args.length < 4) {
            throw new Error('Usage: node wasm-bridge.js <baseURL> <wasmURL> <imageURL> <metaContent>');
        }

        const [baseURL, wasmURL, imageURL, metaContent] = args;
        const result = await runWasm(baseURL, wasmURL, imageURL, metaContent);
        console.log(JSON.stringify(result));
    } catch (error) {
        console.error(JSON.stringify({ error: error.message, stack: error.stack }));
        process.exit(1);
    }
}

main();

