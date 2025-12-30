import http from 'net/http'
import console from 'console'

const mux = http.newServeMux()
mux.handleFunc('/ping', (w, r) => {
    console.log(`${r.remoteAddr} 访问了-${r.url.path}`)
    w.write(`PONG`)
})

mux.handleFunc('/hi', (w, r) => {
    w.write(`HELLO: ${r.remoteAddr}。`)
})

http.listenAndServe(':8080', mux)
