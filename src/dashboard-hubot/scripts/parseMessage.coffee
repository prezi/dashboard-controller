module.exports = (robot) ->

  data = "url=www.9gag.com&slave-id=Jucc"

  robot.hear /post/i, (msg) ->
    msg.send "Posting url to slave"

    robot.http("http://localhost:5000/form-submit")
    .header("content-length",data.length)
    .header("Content-Type","application/x-www-form-urlencoded")
    .post(data) (err, res, body) ->
      try
        msg.send "SUCCESS!!!"
      catch err
        msg.send "Error while sending POST request: " + err
