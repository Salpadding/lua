--
-- Created by IntelliJ IDEA.
-- User: Sal
-- Date: 2019/9/22
-- Time: 17:39
-- To change this template use File | Settings | File Templates.
--

local t = {"a", "b", "c" }
t[2] = "B"
t["foo"] = "Bar"
local s = t[3] .. t[2] .. t[1] .. t["foo"] .. #t

