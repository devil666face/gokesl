package installer

const (
	dpkg     = "dpkg"
	rpm      = "rpm"
	klnagent = "klnagent64"
	kesl     = "kesl"
)

const IpReplacer = "{{ ip }}"

var AgentConfig string = `
EULA_ACCEPTED=yes
KLNAGENT_SERVER={{ ip }}
KLNAGENT_PORT=14000
KLNAGENT_SSLPORT=13000
KLNAGENT_USESSL=Y
KLNAGENT_GW_MODE=1
`

var KeslConfig string = `EULA_AGREED=yes
PRIVACY_POLICY_AGREED=yes
USE_KSN=no
UPDATE_SOURCE=KSCServer
ADMIN_USER=admin
GROUP_CLEAN=no
USE_GUI=yes

INSTALL_LICENSE={{ key }}
UPDATE_EXECUTE=yes
`
