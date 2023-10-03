function include(fpath)
    local f = assert(io.open(fpath, "rb"))
    local content = f:read("*all")
    f:close()

    f = loadstring(content)
    f()
end