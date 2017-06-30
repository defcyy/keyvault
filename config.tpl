APP_NAME = Gogs
RUN_USER = git
RUN_MODE = prod

[database]
DB_TYPE  = mysql
HOST     = {{ .XDOSSQLDBHOST01 }}
NAME     = gogs
USER     = {{ .XDOSSQLDBUSER01 }}
PASSWD   = {{ .XDOSSQLDBPAWD01 }}
SSL_MODE = disable
PATH     = data/gogs.db

[repository]
ROOT = /data/git/gogs-repositories

[server]
DOMAIN       = gogs.cn133.azure.net
HTTP_PORT    = 3000
ROOT_URL     = https://gogs.cn133.azure.net/
DISABLE_SSH  = true
OFFLINE_MODE = false

[mailer]
ENABLED = false

[service]
REGISTER_EMAIL_CONFIRM = false
ENABLE_NOTIFY_MAIL     = false
DISABLE_REGISTRATION   = false
ENABLE_CAPTCHA         = true
REQUIRE_SIGNIN_VIEW    = false

[picture]
DISABLE_GRAVATAR        = false
ENABLE_FEDERATED_AVATAR = true

[session]
PROVIDER = file

[log]
MODE      = console, file
LEVEL     = Info
ROOT_PATH = /app/gogs/log

[security]
INSTALL_LOCK = true
SECRET_KEY   = xSTmVudqm1eg6QV
