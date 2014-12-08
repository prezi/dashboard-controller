twoArbitraryStringsPattern = /(.*) (.*)/i

parseName = (string) ->
    array = string.match twoArbitraryStringsPattern
    name = array[1]

parseUrl = (string) ->
  array = string.match twoArbitraryStringsPattern
  url = array[2]

module.exports = {parseName, parseUrl}
