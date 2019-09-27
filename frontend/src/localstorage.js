let db = {
    save(key, value) {
        await window.backend.ls.setItem(key, JSON.stringify(value))
    },
    get(key, defaultValue = {}) {
        return await JSON.parse(window.backend.ls.getItem(key)) || defaultValue
    },
    remove(key) {
        await window.backend.ls.removeItem(key)
    },
    clear() {
        await window.backend.ls.clear()
    }
}

export default db