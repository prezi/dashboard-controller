DOMAIN = "secret.com"
AUTHORIZATION_HEADER = "secret"

def request(context, flow):
    if flow.request.pretty_host(hostheader=True).endswith(DOMAIN):
        flow.request.headers["Authorization"] = [AUTHORIZATION_HEADER]

