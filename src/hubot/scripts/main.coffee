module.exports = (robot) ->

  SERVER_URL= "http://localhost:5000"

  robot.respond /post (.*) (.*)/i, (msg) ->
    name = msg.match[1]
    url = msg.match[2]
    msg.send "Posting #{url} to #{name}"

    data = "url=" + url + "&slave-id=" + name

    robot.http(SERVER_URL + "/form-submit")
    .header("content-length", data.length)
    .header("Content-Type", "application/x-www-form-urlencoded")
    .post(data) (err, res, body) ->
      try
        msg.send "SUCCESS!!!"
      catch err
        msg.send "Error while sending POST request: " + err

  robot.respond /\b(h+e+l+l+o+)|(h+i+)|(h+e+l+p+)|(h+e+y+)\b/i, (msg) ->
    msg.send "Hi! I'm a very friendly robot sending urls to different dashboards around Prezi Office.\n
    You can get the list of available dashboards by typing info.\n
    Post on dashboards by:\n
    @DashboardController post slave_name _url_to_display"


  robot.respond /\b(i+n+f+o+)\b/i, (msg) ->
    robot.http(SERVER_URL + "/slavemap")
    .header("Content-Type", "application/json")
    .get() (err, response, body) ->
      try
        slaveNameArray = JSON.parse body
        if slaveNameArray == null
          msg.send "No slaves in service, Captain!"
        else
          slaveNames = ""
          for slaveName in slaveNameArray
            slaveNames += (slaveName + " ")
          msg.send "List of available slaves: " + slaveNames
      catch err
        msg.send "Error getting slave data: " + err
