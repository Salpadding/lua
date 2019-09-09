-- 返回一个"iterator"。每次被调用，从给定的值的序列中找出一个单词。
function allwords ()
    local line = io.read()    -- 输入给定的值的序列。
    local pos = 1    -- 值的序列中搜索到了哪里。
    return function ()    -- "iterator function"。
        while line do
            -- 从给定的值的序列中找出一个单词，返回单词的起始和终止位置。
            local s, e = string.find(line, "%w+", pos)
            if s then    -- 找到了一个单词。
                pos = e + 1    -- 更新搜索位置。
                return string.sub(line, s, e)    -- 返回找到的单词。
            else    -- 没找到单词。
                line = io.read()    -- 读取下一行。
                pos = 1             -- 值的序列中搜索位置重置。
            end
        end
        return nil
    end
end

-- 将两个索引值组合成为一个索引。
function prefix (w1, w2)
    return w1 .. ' ' .. w2
end

-- 将索引与值的映射存入"statetab"。
function insert (index, value)
    if not statetab[index] then
        -- 一个索引可能对应多个值（["more we"] = {"try", "do"}），所以"table"存储。
        statetab[index] = {}
    end
    table.insert(statetab[index], value)
end

local N = 2    -- 索引值限定为两个。
local MAXGEN = 10000    -- 最大产生的结果数量。
local NOWORD = "\\n"    -- 默认的索引值是"\n"。

-- 创建映射关系并存储。
statetab = {}    -- 存储映射关系的"table"。
local w1, w2 = NOWORD, NOWORD    -- 初始的两个索引值都是"\n"。
for w in allwords() do    -- 产生索引值。
    insert(prefix(w1, w2), w)    -- 存储映射关系。
    w1 = w2; w2 = w;    -- 向后依次更替索引值。
end
insert(prefix(w1, w2), NOWORD)    -- 存储最后一个映射关系（[we do] = {"\n", }）。

-- 打印"statetab"中的内容。
for k, v in pairs(statetab) do
    io.write(string.format("[%s] = {", k))
    for m, n in ipairs(v) do
        io.write(string.format("\"%s\", ", n))
    end
    io.write("}\n")
end

--[[ 以初始索引在"statetab"中随机取值，
     不断的以新的值与旧的索引值组合成为新的索引，再次在"statetab"中随机取值，
     循环往复，打印出找打的值。]]
w1 = NOWORD; w2 = NOWORD    -- 重新初始化索引值为默认的索引值。
for i = 1, MAXGEN do    -- 最大结果数量为"MAXGEN"个。
    local list = statetab[prefix(w1, w2)]
    -- 产生随机数种子。
    math.randomseed(tonumber(tostring(os.time()):reverse():sub(1, 6)))
    --[[ 从"list"中随机选择一个元素，
         比如[more we] = {"try", "do", }对应"try"和"do"两个元素，随机选择一个。]]
    local r = math.random(#list)    -- 生成随机数。
    local nextword = list[r]
    if nextword == NOWORD then    -- 如果到了默认的索引值，就不再找了。
        break
    end
    io.write(nextword, " ")
    w1 = w2; w2 = nextword    -- 不断的以新的值与旧的索引值组合成为新的索引。
end
io.write("\n")
--————————————————
--版权声明：本文为CSDN博主「vermilliontear」的原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接及本声明。
--原文链接：https://blog.csdn.net/VermillionTear/article/details/50555556