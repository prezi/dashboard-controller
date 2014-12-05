module.exports = (robot) ->

  robot.hear /post (.*) (.*)/i, (msg) ->
    name = msg.match[1]
    url = msg.match[2]
    msg.send "Posting #{url} to #{name}"

    data = "url=" + url + "&slave-id=" + name

    robot.http("http://10.0.1.214:5000/form-submit")
    .header("content-length", data.length)
    .header("Content-Type", "application/x-www-form-urlencoded")
    .post(data) (err, res, body) ->
    try
      msg.send "SUCCESS!!!"
    catch err
      msg.send "Error while sending POST request: " + err
