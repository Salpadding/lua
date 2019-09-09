math.randomseed(os.time())
local trials = 100
local iter = 1
local targetpercent = 0.99
function recur()
    local tot = 0
    for c = 1, trials do
        local b = false
        local x = 1
        while b == false and x <= iter do
            r = math.random(127)
            if r <= 127 and r > 63 then
                tot = tot + 1
                b = true
            end
            x = x + 1
        end

    end
    if tot/trials <= targetpercent then
        iter = iter + 1
        x = 1
        tot = 0
        recur()
    end
end
local totiter = 0
for c = 1, 100 do
    recur()
    totiter = iter + totiter
    iter = 1
end
print (totiter/trials)
print ("588 tier 7, 268 tier 6, 120 tier 5, 58 tier 4, 31 tier 3, 15 tier 2, 7 tier 1")

function entry1 (o)
    N=N + 1
    local title = o.title or o.org or 'org'
    fwrite('<HR>\n<H3>\n')
    local href = ''

    if o.url then
        href = string.format(' HREF="%s"', o.url)
    end
    fwrite('<A NAME="%d"%s>%s</A>\n', N, href, title)

    if o.title and o.org then
        fwrite('<BR>\n<SMALL><EM>%s</EM></SMALL>', o.org)
    end
    fwrite('\n</H3>\n')

    if o.description then
        fwrite('%s', string.gsub(o.description,
            '\n\n\n*', '<P>\n'))
        fwrite('<P>\n')
    end

    if o.email then
        fwrite('Contact: <A HREF="mailto:%s">%s</A>\n',
            o.email, o.contact or o.email)
    elseif o.contact then
        fwrite('Contact: %s\n', o.contact)
    end
end