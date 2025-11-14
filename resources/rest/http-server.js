import console from 'console'
import time from 'time'
import url from 'net/url'
import http from 'net/http'
import httputil from 'net/http/httputil'

const target = url.parse('https://mirrors.zju.edu.cn/')
const proxy = httputil.newSingleHostReverseProxy(target)

let cnt = 0
const mux = http.newServeMux()
mux.handleFunc('/', (w, r) => {
    cnt++
    w.header().set('Content-Type', 'text/html; charset=utf8')
    const content = `<h1>你访问了<kbd>${r.url.path}</kbd> ，网站总访问量 ${cnt} </h1>`
    w.write(content)

    const log = `[${new Date().toJSON()}] ${r.remoteAddr} 第 ${cnt} 次访问：${r.url.path} ${r.url.rawQuery}`
    console.log(log)
})
mux.handleFunc('/favicon.ico', (w, r) => {
   w.writeHeader(404)
})

const opt = {
    addr: '0.0.0.0:8888',
    handler: mux,
    readTimeout: time.minute,
    readHeaderTimeout: 5 * time.second
}
const handle = http.serve(opt)
handle.wait()
