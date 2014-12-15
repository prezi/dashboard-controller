module.exports = (robot) ->

  SERVER_URL= "http://10.0.0.210:5000"

  robot.respond /post (.*) (.*)/i, (msg) ->
    name = msg.match[1]
    url = msg.match[2]
    msg.send "Posting #{url} to #{name}"

    data = "url=" + url + "&slave-id=" + name

    robot.http(SERVER_URL= "http://10.0.0.210:5000" + "/form-submit")
    .header("content-length", data.length)
    .header("Content-Type", "application/x-www-form-urlencoded")
    .post(data) (err, res, body) ->
      try
        msg.send "Success!"
      catch err
        msg.send "Error while sending POST request."
        msg.send err

  robot.respond /\b(h+e+l+l+o+)|(h+i+)|(help\s*(.*)?$)|(h+e+y+)\b/i, (msg) ->
    msg.send "Hi! I'm a very friendly robot sending urls to different dashboards around Prezi Office."
    msg.send "You can get the list of available dashboards by typing info."
    msg.send "Post on dashboards by: "
    msg.send "@DashboardController post slave_name _url_to_display"

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
            slaveNames += (slaveName + "\n")
          msg.send "List of available slaves: \n" + slaveNames
      catch err
        msg.send "Error getting slave data: " + err

  GoT = [
    "http://img2.wikia.nocookie.net/__cb20140518100049/gameofthrones/images/4/47/Ygritte-Profile-HD.png",
    "http://img3.wikia.nocookie.net/__cb20140326175142/gameofthrones/images/e/ed/Charles-Dance-as-Tywin-Lannister_photo-Macall-B.Polay_HBO.jpg",
    "http://img3.wikia.nocookie.net/__cb20110626030942/gameofthrones/images/9/9c/EddardStark.jpg",
    "http://img2.wikia.nocookie.net/__cb20140515200719/gameofthrones/images/f/fb/Catelyn-Stark-Profile-HD.png",
    "http://img3.wikia.nocookie.net/__cb20140515152752/gameofthrones/images/9/90/Oberyn-Martell-Profile-HD.png"
  ]

  robot.hear /(G+o+T me)|(g+o+t m+e)/i, (msg) ->
    msg.send msg.random GoT

  robot.hear /(V+a+l+a+r M+o+r+g+h+u+l+i+s)|(v+a+l+a+r m+o+r+g+h+u+l+i+s)/i, (msg) ->
    msg.send "Valar Dohaeris."

  robot.hear /(W+i+n+t+e+r i+s c+o+m+i+n+g)|(w+i+n+t+e+r i+s c+o+m+i+n+g)/i, (msg) ->
    msg.send "such winter. oh no. much soon. wow. snow."

  dracarys = [
    "http://www.majoh.com/art/2d/Dracarys_large.jpg",
    "http://i.ytimg.com/vi/Nae4VnNC9AE/maxresdefault.jpg",
    "http://sd.keepcalm-o-matic.co.uk/i/-keep-calm-and-dracarys-.png"
  ]

  robot.hear /d+r+a+c+a+r+y+s/i, (msg) ->
    msg.send msg.random dracarys

  robot.respond /s+p+o+i+l+e+r/i, (msg) ->
    msg.send "All men must die."
    msg.send "http://www.slate.com/articles/arts/television/2014/04/game_of_thrones_deaths_mourn_dead_characters_at_their_virtual_graveyard.html"
