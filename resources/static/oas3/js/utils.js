class Utils {
    static oas3() {
        const key = 'oas3'
        const quires = new URLSearchParams(window.location.search)
        let name = quires.get(key)
        if (!name) {
            name = sessionStorage.getItem(key);
        }

        if (name) {
            return name
        }

        const warning = document.createElement('a')
        warning.href = './'
        warning.innerText = '未找到接口文档，点击跳转至接口文档主页'
        warning.className = 'sl-text-6xl'
        document.body.replaceChildren(warning)

        return ''
    }

}
