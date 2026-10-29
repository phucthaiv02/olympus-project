# Each case: (test group, case name, Lua chunk). The chunk is loaded with
# chunk name "=case" and its results are rendered by `show` in the prelude.
PRELUDE = r'''
local function show(v)
  if type(v) ~= "table" then return tostring(v) end
  local keys = {}
  for k in pairs(v) do keys[#keys + 1] = k end
  table.sort(keys)
  local out = {}
  for _, k in ipairs(keys) do
    local x = v[k]
    if k == "func" then
      x = type(x)
    elseif k == "activelines" then
      local lines = {}
      for line in pairs(x) do lines[#lines + 1] = line end
      table.sort(lines)
      x = "{" .. table.concat(lines, ",") .. "}"
    elseif type(x) == "table" then
      x = "table"
    end
    out[#out + 1] = k .. "=" .. tostring(x)
  end
  return "{" .. table.concat(out, " ") .. "}"
end

function run(f)
  local function pack(...) return select("#", ...), {...} end
  local n, results = pack(pcall(f))
  local out = {}
  for i = 1, n do out[i] = show(results[i]) end
  return table.concat(out, " | ")
end
'''

CASES = []
def case(group, name, code):
    CASES.append((group, name, code.strip("\n")))

# ---------------------------------------------------------------- getinfo fields
G = "GetInfo"
case(G, "source of a Lua function", '''
local function f()
  return 1
end
return debug.getinfo(f, "S")
''')
case(G, "source of the main chunk", '''
return debug.getinfo(1, "S")
''')
case(G, "source of a Go function", '''
return debug.getinfo(print, "S")
''')
case(G, "current line", '''
local a = 1

local t = debug.getinfo(1, "l")
return t.currentline
''')
case(G, "current line of the caller", '''
local function f()
  return debug.getinfo(2, "l").currentline
end
local x = 0
local y = f()
return y
''')
case(G, "current line of a function value", '''
return debug.getinfo(function() end, "l"), debug.getinfo(print, "l")
''')
case(G, "upvalues and parameters", '''
local a, b = 1, 2
local function f(x, y, z) return a + b end
local function g(...) return f end
return debug.getinfo(f, "u"), debug.getinfo(g, "u"), debug.getinfo(print, "u"), debug.getinfo(load("return 1"), "u")
''')
case(G, "tail call flag", '''
local function a() return debug.getinfo(1, "t") end
local function b() return a() end
local x = a()
local y = b()
return x, y, debug.getinfo(a, "t")
''')
case(G, "func field", '''
local function f() return debug.getinfo(1, "f").func end
return f() == f, debug.getinfo(0, "f").func == debug.getinfo, debug.getinfo(print, "f").func == print, debug.getinfo(f, "f").func == f
''')
case(G, "active lines", '''
local f = function(a)
  if a then
    return 1
  end

  -- comment
  return 2
end
return debug.getinfo(f, "L"), debug.getinfo(print, "L")
''')
case(G, "active lines of a chunk", '''
return debug.getinfo(load("local a = 1\\n\\nlocal b = 2\\nreturn a"), "L")
''')
case(G, "default options", '''
local function f(x, ...)
  return debug.getinfo(1)
end
local t = f()
return t
''')
case(G, "selected options only", '''
return debug.getinfo(1, ""), debug.getinfo(1, "lt"), debug.getinfo(print, "")
''')
case(G, "level out of range", '''
return debug.getinfo(100), debug.getinfo(-1), debug.getinfo(100, "S")
''')
case(G, "numeric string level", '''
local t = debug.getinfo("1", "l")
return t
''')
case(G, "level of a Go function", '''
local r
pcall(function() r = debug.getinfo(2, "Sl") end)
return r
''')
case(G, "levels through Go and tail calls", '''
local function a() return debug.getinfo(2, "Sl"), debug.getinfo(3, "Sl") end
local function b() return a() end
local function c()
  local x, y = b()
  return x, y
end
return c()
''')
case(G, "info of a hook frame", '''
local seen
local function f()
  local x = 1
  return x
end
debug.sethook(function(event, line) if not seen then seen = debug.getinfo(2, "Sl") end end, "l")
f()
debug.sethook()
return seen
''')

# ---------------------------------------------------------------- names
N = "GetInfoNames"
case(N, "global", '''
function gfn() local t = debug.getinfo(1, "n") return t end
local r = gfn()
return r
''')
case(N, "local", '''
local function lfn() return debug.getinfo(1, "n") end
local r = lfn()
return r
''')
case(N, "method", '''
local o = {}
function o:m() return debug.getinfo(1, "n") end
local r = o:m()
return r
''')
case(N, "field", '''
local o = {f = function() return debug.getinfo(1, "n") end}
T = {}
function T.g() return debug.getinfo(1, "n") end
local r = o.f()
local s = T.g()
return r, s
''')
case(N, "upvalue", '''
local up = function() return debug.getinfo(1, "n") end
local function g() local r = up() return r end
return g()
''')
case(N, "metamethods", '''
local out = {}
local function name() local t = debug.getinfo(2, "n") return tostring(t.name) .. "/" .. t.namewhat end
local mt = {}
mt.__index = function() out[#out + 1] = name() end
mt.__newindex = function() out[#out + 1] = name() end
mt.__add = function() out[#out + 1] = name() return 0 end
mt.__concat = function() out[#out + 1] = name() return "" end
mt.__lt = function() out[#out + 1] = name() return true end
mt.__le = function() out[#out + 1] = name() return true end
mt.__eq = function() out[#out + 1] = name() return true end
mt.__len = function() out[#out + 1] = name() return 0 end
mt.__unm = function() out[#out + 1] = name() return 0 end
mt.__call = function() out[#out + 1] = name() end
local a, b = setmetatable({}, mt), setmetatable({}, mt)
local _ = a.x
a.y = 1
_ = a + 1
_ = a .. "s"
_ = a < b
_ = a <= b
_ = a == b
_ = #a
_ = -a
a()
return table.concat(out, ",")
''')
case(N, "for iterator", '''
local r
for _ in function(s, c) if c then return nil end r = debug.getinfo(1, "n") return 1 end do end
return r
''')
case(N, "called from Go", '''
local r = select(2, pcall(function() return debug.getinfo(1, "n") end))
local s
table.sort({3, 1, 2}, function(x, y) s = debug.getinfo(1, "n") return x < y end)
return r, s
''')
case(N, "tail call", '''
local function t() return debug.getinfo(1, "n") end
local function u() return t() end
local r = u()
return r
''')
case(N, "Go functions", '''
local r
pcall(function() r = debug.getinfo(2, "n") end)
local s
table.sort({2, 1}, function(x, y) s = debug.getinfo(2, "n") return x < y end)
return debug.getinfo(0, "n"), r, s
''')
case(N, "function value", '''
return debug.getinfo(print, "n"), debug.getinfo(function() end, "n")
''')

# ---------------------------------------------------------------- short_src
S = "ShortSource"
case(S, "literal and file names", '''
local function src(code, name) return load(code, name)() end
local q = "return debug.getinfo(1, 'S').short_src"
return src(q, "=my chunk"), src(q, "@dir/file.lua"), src(q, "=" .. string.rep("n", 80)), src(q, "@" .. string.rep("d/", 40) .. "file.lua")
''')
case(S, "string sources", '''
local q = "return debug.getinfo(1, 'S').short_src"
local exact = q .. " --" .. string.rep("z", 44 - #q - 3)
local over = q .. " --" .. string.rep("z", 45 - #q - 3)
return load(q)(), load("local a = 1\\n" .. q)(), load(q .. " --" .. string.rep("x", 80))(), #exact, load(exact)(), #over, load(over)()
''')
case(S, "error positions", '''
local ok1, m1 = pcall(load("local x = nil\\nerror('boom')\\n-- " .. string.rep("y", 60)))
local ok2, m2 = pcall(load("error('boom')", "@" .. string.rep("p", 70) .. ".lua"))
local ok3, m3 = pcall(load("error('boom')", "=" .. string.rep("q", 70)))
return m1, m2, m3
''')

# ---------------------------------------------------------------- argument errors
E = "Errors"
case(E, "getinfo", '''
local function check(...)
  local ok, m = pcall(debug.getinfo, ...)
  return tostring(ok) .. ":" .. tostring(string.find(m, "bad argument #", 1, true) ~= nil and m:sub(m:find("#", 1, true), m:find("#", 1, true) + 1)) .. ":" .. tostring(string.find(m, "invalid option", 1, true) ~= nil) .. ":" .. tostring(string.find(m, "function or level expected", 1, true) ~= nil)
end
return check(1, "X"), check(print, "S>"), check("abc"), check({}), check()
''')
case(E, "getlocal and setlocal", '''
local function check(f, ...)
  local ok, m = pcall(f, ...)
  return tostring(ok) .. ":" .. m:sub(m:find("#", 1, true), m:find("#", 1, true) + 1) .. ":" .. tostring(string.find(m, "level out of range", 1, true) ~= nil)
end
return check(debug.getlocal, 100, 1), check(debug.getlocal), check(debug.getlocal, 1, {}), check(debug.setlocal, 100, 1, 1), check(debug.setlocal, 1, 1)
''')

# ---------------------------------------------------------------- getlocal
L = "GetLocal"
case(L, "locals in order", '''
local a, b = 10, "s"
local out = {}
local i = 1
while true do
  local name, v = debug.getlocal(1, i)
  if not name then break end
  out[#out + 1] = name .. "=" .. (type(v) == "table" and "table" or tostring(v))
  i = i + 1
end
return table.concat(out, ",")
''')
case(L, "parameters", '''
local function f(x, y)
  local z = x + y
  return debug.getlocal(1, 1), debug.getlocal(1, 2), debug.getlocal(1, 3)
end
return f(7, 8)
''')
case(L, "scopes", '''
do local inner = 1 end
local after = 2
local function f()
  local early = debug.getlocal(1, 1)
  local later = 1
  return early, debug.getlocal(1, 1), debug.getlocal(1, 2)
end
return debug.getlocal(1, 1), f()
''')
case(L, "loop control variables", '''
local function names(level)
  local out = {}
  local j = 1
  while true do
    local name, v = debug.getlocal(level + 1, j)
    if not name then break end
    out[#out + 1] = name .. "=" .. (type(v) == "function" and "fn" or type(v) == "table" and "table" or tostring(v))
    j = j + 1
  end
  return table.concat(out, ",")
end
local r1, r2
for i = 5, 5 do r1 = names(1) end
for k, v in pairs({a = 1}) do r2 = names(1) end
return r1, r2
''')
case(L, "temporaries and out of range", '''
local a = 1
return debug.getlocal(1, 2), debug.getlocal(1, 50), debug.getlocal(1, 0)
''')
case(L, "varargs", '''
local function f(a, ...)
  return debug.getlocal(1, -1), debug.getlocal(1, -2), debug.getlocal(1, -3)
end
local function g(...) return debug.getlocal(1, -1) end
local function h(a) return debug.getlocal(1, -1) end
return f(1, "v1", "v2"), g(), h(1)
''')
case(L, "varargs of the main chunk", '''
return load("return debug.getlocal(1, -1)")("va")
''')
case(L, "function parameters", '''
local function f(a, b, ...) local c = 1 end
return debug.getlocal(f, 1), debug.getlocal(f, 2), debug.getlocal(f, 3), debug.getlocal(f, -1), debug.getlocal(print, 1)
''')
case(L, "caller levels", '''
local outer = "o"
local function f() return debug.getlocal(2, 1) end
local r = {}
pcall(function() r = {debug.getlocal(2, 1)} end)
return f(), r[1], type(r[2]), debug.getlocal(0, 1), debug.getlocal(0, 2), debug.getlocal("1", "1")
''')

# ---------------------------------------------------------------- setlocal
T = "SetLocal"
case(T, "current function", '''
local a = 1
local name = debug.setlocal(1, 1, 42)
return name, a
''')
case(T, "caller", '''
local x = "old"
local function f() return debug.setlocal(2, 1, "new") end
local n = f()
return n, x
''')
case(T, "out of range", '''
local a = 1
return debug.setlocal(1, 10, 5), a, debug.setlocal(1, 0, 5), a
''')
case(T, "varargs and parameters", '''
local function f(...)
  debug.setlocal(1, -2, "changed")
  return ...
end
local function g(p)
  debug.setlocal(1, 1, p * 2)
  return p
end
return f("a", "b", "c"), g(21)
''')
case(T, "captured locals", '''
local x = 1
local function get() return x end
debug.setlocal(1, 1, 99)
local y = 1
debug.setlocal(1, 3, "s")
local function get2() return y end
return x, get(), y, get2()
''')
case(T, "temporaries", '''
local a = 1
return debug.setlocal(1, 2, 5)
''')

# ---------------------------------------------------------------- hooks
H = "Hooks"
case(H, "line events", '''
local lines = {}
local function f(n)
  local a = 1
  for i = 1, n do
    a = a + i
  end
  if a > 3 then
    a = 0
  end
  return a
end
debug.sethook(function(event, line) lines[#lines + 1] = line end, "l")
f(2)
debug.sethook()
return table.concat(lines, ",")
''')
case(H, "line events of one-line loops", '''
local lines = {}
local x = 0
debug.sethook(function(event, line) lines[#lines + 1] = line end, "l")
for i = 1, 3 do x = x + i end
local function g() return 1 end local y = g() + g()
debug.sethook()
return table.concat(lines, ","), x
''')
case(H, "line events in functions called from Go", '''
local lines = {}
local function cmp(a, b)
  return a < b
end
debug.sethook(function(event, line) lines[#lines + 1] = line end, "l")
pcall(cmp, 1, 2)
debug.sethook()
return table.concat(lines, ",")
''')
case(H, "current line in hooks", '''
local out = {}
local function f()
  local x = 1
  x = x + 1
  return x
end
debug.sethook(function(event, line)
  local t = debug.getinfo(2, "lf")
  if t.func == f then out[#out + 1] = event .. ":" .. tostring(line) .. ":" .. t.currentline end
end, "crl")
f()
debug.sethook()
return table.concat(out, ",")
''')
case(H, "count events", '''
local function f() local x = 0 for i = 1, 10 do x = x + i end return x end
local n1, n3 = 0, 0
debug.sethook(function(event) n1 = n1 + 1 end, "", 1)
f()
debug.sethook()
debug.sethook(function(event) n3 = n3 + 1 end, "", 3)
f()
debug.sethook()
return n1, n3
''')
case(H, "current line in count hooks", '''
local lines = {}
local function f()
  local x = 1
  x = x + 1
  return x
end
debug.sethook(function(event)
  local t = debug.getinfo(2, "lf")
  if t.func == f then lines[#lines + 1] = event .. t.currentline end
end, "", 1)
f()
debug.sethook()
return table.concat(lines, ",")
''')
case(H, "call and return events", '''
local evs = {}
local function a() return 1 end
local function b() return a() end
debug.sethook(function(event)
  local t = debug.getinfo(2, "nSt")
  evs[#evs + 1] = event .. ":" .. tostring(t.name) .. ":" .. t.what .. ":" .. tostring(t.istailcall)
end, "cr")
b()
local s = tostring(1)
debug.sethook()
return table.concat(evs, ",")
''')
case(H, "locals in line hooks", '''
local seen = {}
local function f()
  local watch = 1
  watch = 2
  watch = 3
  return watch
end
debug.sethook(function(event, line)
  local t = debug.getinfo(2, "f")
  local name, v = debug.getlocal(2, 1)
  if t.func == f and name == "watch" then
    seen[#seen + 1] = line .. "=" .. tostring(v)
  end
end, "l")
f()
debug.sethook()
return table.concat(seen, ",")
''')
case(H, "setlocal from a line hook", '''
local function f()
  local v = 1
  v = v + 1
  return v
end
debug.sethook(function(event, line)
  if line == 4 and debug.getinfo(2, "f").func == f then debug.setlocal(2, 1, 100) end
end, "l")
local r = f()
debug.sethook()
return r
''')
case(H, "gethook", '''
local h = function() end
debug.sethook(h, "crl", 5)
local f, mask, count = debug.gethook()
debug.sethook()
local g, mask2, count2 = debug.gethook()
debug.sethook(h, "l")
local _, mask3, count3 = debug.gethook()
debug.sethook()
return f == h, mask, count, g, mask2, count2, mask3, count3
''')
